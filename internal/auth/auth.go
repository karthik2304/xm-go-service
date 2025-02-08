package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt"
	"github.com/karthik2304/xm-go-service/configs"
	echomiddleware "github.com/oapi-codegen/echo-middleware"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(configs.Settings.APP_SECURE_SECRETKEY)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 1 day expiry
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GetTestToken(userName string) string {
	rr, _ := GenerateToken(userName)
	return rr
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid JWT token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	return claims, nil
}

func OpenAPIAuthFunc(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	eCtx := echomiddleware.GetEchoContext(ctx)
	if eCtx == nil {
		return errors.New("could not retrieve Echo context")
	}
	if eCtx.Path() == "/v1/auth/signup" || eCtx.Path() == "/v1/auth/login" {
		return nil
	}
	authHeader := eCtx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return errors.New("invalid Authorization format")
	}

	tokenString := parts[1]

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	eCtx.Set("userClaims", claims)

	if userName, exists := claims["username"].(string); exists {
		eCtx.Set("username", userName)
	} else {
		return errors.New("userName not found in token")
	}

	return nil
}
