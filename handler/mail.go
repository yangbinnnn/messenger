package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/yangbinnnn/messenger/g"
)

type Mail struct {
	Token   string `json:"token" form:"token" query:"token"`
	TOS     string `json:"tos" form:"tos" query:"tos"`
	Subject string `json:"subject" form:"subject" query:"subject"`
	Content string `json:"content" form:"content" query:"content"`
}

func (m Mail) Validate() error {
	if m.Token != g.Config().Http.Token {
		return echo.ErrForbidden
	}
	if m.TOS == "" || m.Subject == "" || m.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "tos or subject or content requried")
	}
	return nil
}

func (h *Handler) SendMail(c echo.Context) error {
	if !g.Config().Smtp.Enable {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	m := new(Mail)
	if err := c.Bind(m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := m.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tos := strings.Split(m.TOS, ",")
	err := h.mailcli.Send(tos, m.Subject, m.Content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "success")
}
