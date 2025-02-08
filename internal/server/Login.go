package server

import (
	"context"
	"fmt"
	"net/http"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/auth"
	"github.com/karthik2304/xm-go-service/internal/utils"
	"github.com/labstack/echo/v4"
)

// Login implements Openapi.ServerInterface.

func (s *Server) Login(ctx echo.Context) error {

	body := Openapi.LoginJSONRequestBody{}
	err := ctx.Bind(&body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	userDBPassword, err := s.db.GetAuth(context.Background(), body.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}
	// Verify password
	if !auth.VerifyPassword(userDBPassword.Password, body.Password) {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}
	token, err := auth.GenerateToken(body.Username)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating token"})
	}

	err = s.WriteMessage(Openapi.Event{
		EventType: "LOGIN-EVENT",
		Id:        utils.GetUniqueId(),
		Timestamp: utils.GetCurrentTime(),
		UserName:  body.Username,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	return ctx.JSON(http.StatusOK, Openapi.JwtResponse{JwtToken: token})

}
