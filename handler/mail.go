package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/yangbinnnn/messenger/sender"
)

func (h *Handler) SendMail(c echo.Context) error {
	if !h.cfg.Smtp.Enable {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	client, err := sender.NewMailClient(
		h.cfg.Smtp.Addr,
		h.cfg.Smtp.Username,
		h.cfg.Smtp.Password,
		h.cfg.Smtp.From,
		h.cfg.Smtp.Timeout,
		h.cfg.Smtp.TLS,
		false,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	token := h.ParamString(c, "token")
	if token != h.cfg.Http.Token {
		return echo.ErrForbidden
	}

	tosStr := h.ParamString(c, "tos")
	subject := h.ParamString(c, "subject")
	content := h.ParamString(c, "content")
	if tosStr == "" || subject == "" || content == "" {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	tos := strings.Split(tosStr, ",")
	err = client.Send(tos, subject, content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "success")
}
