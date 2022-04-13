package context

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Error   string      `json:"error,omitempty"`
	ErrHint string      `json:"hint,omitempty"`
	Data    interface{} `json:"data"`
}

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   data,
	})
}

func Error(c echo.Context, status int, hint string, err error) error {
	return c.JSON(status, Response{
		Status:  status,
		Error:   err.Error(),
		ErrHint: hint,
	})
}
