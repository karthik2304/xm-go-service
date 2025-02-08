package pkg

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type RequestValidator struct {
	validator *validator.Validate
}

func (rv *RequestValidator) Validate(i interface{}) error {
	if err := rv.validator.Struct(i); err != nil {
		return fmt.Errorf("failed to validated struct: %w", err)
	}

	return nil
}

func New() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Validator = &RequestValidator{validator: validator.New()}
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		ErrorHandler(err, c, nil)
	}
	e.Use(middleware.CORS())
	return e
}

func ErrorHandler(err error, c echo.Context, logger *zap.Logger) {
	if c.Response().Committed {
		return
	}
	jsonErr := c.JSON(http.StatusBadRequest, &echo.HTTPError{
		Message: err.Error(),
	})

	if jsonErr != nil {
		logger.Warn("failed to return json error", zap.Error(err))
	}
}
