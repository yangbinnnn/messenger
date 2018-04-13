package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func (h *Handler) SendWeChat(c echo.Context) error {
	if !h.cfg.Wechat.Enable {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	token := h.ParamString(c, "token")
	if token != h.cfg.Http.Token {
		return echo.ErrForbidden
	}

	tosStr := h.ParamString(c, "tos")
	content := h.ParamString(c, "content")
	if tosStr == "" || content == "" {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.wechat.SendMsg(tosStr, content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "success")
}
