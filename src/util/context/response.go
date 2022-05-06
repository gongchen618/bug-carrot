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

// Success 是 http 函数成功运行时的消息返回结构
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   data,
	})
}

// Error 是 http 函数失败运行时的消息返回结构
func Error(c echo.Context, status int, hint string, err error) error {
	return c.JSON(status, Response{
		Status:  status,
		Error:   err.Error(),
		ErrHint: hint,
	})
}
