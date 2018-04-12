package handler

import (
	"github.com/labstack/echo"
)

type Handler struct {
}

func (h *Handler) ParamString(c echo.Context, key string) string {
	value := c.QueryParam(key)
	if value != "" {
		return value
	}
	value = c.FormValue(key)
	return value
}
