package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/yangbinnnn/messenger/g"
	"github.com/yangbinnnn/messenger/sender"
)

func (h *Handler) SendMail(c echo.Context) error {
	cfg := config.Config()
	if !cfg.Smtp.Enable {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	client, err := sender.NewMailClient(
		cfg.Smtp.Addr,
		cfg.Smtp.Username,
		cfg.Smtp.Password,
		cfg.Smtp.From,
		cfg.Smtp.Timeout,
		cfg.Smtp.TLS,
		false,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	token := h.ParamString(c, "token")
	if token != cfg.Http.Token {
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
