package server

import (
	"context"
	"fmt"
	"net/http"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/auth"
	"github.com/karthik2304/xm-go-service/internal/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// SignUp implements Openapi.ServerInterface.
func (s *Server) SignUp(ctx echo.Context) error {
	body := Openapi.LoginJSONRequestBody{}
	err := ctx.Bind(&body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	// verify userName
	data, err := s.db.GetAuth(context.Background(), body.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}
	if data != nil && data.Username == body.Username {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: "UserName can't be Same"})
	}

	body.Password, err = auth.HashPassword(body.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}
	err = s.db.CreateAuth(context.Background(), Openapi.SignUpJSONRequestBody(body))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}
	err = s.WriteMessage(Openapi.Event{
		EventType: "SIGNUP-EVENT",
		Id:        utils.GetUniqueId(),
		Timestamp: utils.GetCurrentTime(),
		UserName:  body.Username,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	return ctx.JSON(http.StatusOK, Openapi.Success{Message: "User Created Successfully"})
}
