package server

import (
	"context"
	"fmt"
	"net/http"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/utils"
	"github.com/labstack/echo/v4"
)

// ListCompanies implements Openapi.ServerInterface.
func (s *Server) ListCompanies(ctx echo.Context) error {
	userName, ok := ctx.Get("username").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, Openapi.Error{Message: "No UserName in context"})
	}
	data, err := s.db.GetAll(context.Background())
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: err.Error()})
	}
	err = s.WriteMessage(Openapi.Event{
		EventType: "GET-COMPANIES-ALL-EVENT",
		Id:        utils.GetUniqueId(),
		Timestamp: utils.GetCurrentTime(),
		UserName:  userName,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}
	return ctx.JSON(http.StatusOK, data)
}
