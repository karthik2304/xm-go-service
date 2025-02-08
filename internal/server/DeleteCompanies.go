package server

import (
	"context"
	"fmt"
	"net/http"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/utils"
	"github.com/labstack/echo/v4"
)

// DeleteCompanyDetails implements Openapi.ServerInterface.
func (s *Server) DeleteCompanyDetails(ctx echo.Context, companyUuid string) error {

	userName, ok := ctx.Get("username").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, Openapi.Error{Message: "No UserName in context"})
	}
	err := s.db.Delete(context.Background(), map[string]interface{}{
		"companyuuid": companyUuid,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	err = s.WriteMessage(Openapi.Event{
		EventDetails: &Openapi.EventDetails{
			CompanyUUID: companyUuid,
		},
		EventType: "DELETE-COMPANY-EVENT",
		Id:        utils.GetUniqueId(),
		Timestamp: utils.GetCurrentTime(),
		UserName:  userName,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}
	return ctx.JSON(http.StatusOK, Openapi.Success{Message: "Company deleted successfully"})
}
