package server

import (
	"context"
	"fmt"
	"net/http"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"github.com/karthik2304/xm-go-service/internal/utils"
	"github.com/labstack/echo/v4"
)

// UpdateCompanyDetails implements Openapi.ServerInterface.
func (s *Server) UpdateCompanyDetails(ctx echo.Context, companyUuid string) error {

	userName, ok := ctx.Get("username").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, Openapi.Error{Message: "No UserName in context"})
	}
	body := Openapi.UpdateCompanyDetailsJSONRequestBody{}
	err := ctx.Bind(&body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}
	err = s.db.Update(context.Background(), map[string]interface{}{
		"companyuuid": companyUuid,
	}, body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}
	err = s.WriteMessage(Openapi.Event{
		EventType: "UPDATE-EVENT",
		Id:        utils.GetUniqueId(),
		Timestamp: utils.GetCurrentTime(),
		UserName:  userName,
		EventDetails: &Openapi.EventDetails{
			CompanyName:    body.CompanyName,
			CompanyUUID:    companyUuid,
			Type:           body.Type,
			Registered:     body.Registered,
			TotalEmployees: body.TotalEmployees,
		},
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Openapi.Error{Message: fmt.Sprintf("%v", err)})
	}

	return ctx.JSON(http.StatusOK, Openapi.Success{Message: "Company updated successfully"})
}
