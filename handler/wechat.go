package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/yangbinnnn/messenger/g"
)

type Chat struct {
	Token   string `json:"token" form:"token" query:"token"`
	TOS     string `json:"tos" form:"tos" query:"tos"`
	Content string `json:"content" form:"content" query:"content"`
}

func (c Chat) Validate() error {
	if c.Token != g.Config().Http.Token {
		return echo.ErrForbidden
	}
	if c.TOS == "" || c.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "tos or content requried")
	}
	return nil
}

func (h *Handler) SendWeChat(c echo.Context) error {
	if !g.Config().Wechat.Enable {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	chat := new(Chat)
	if err := c.Bind(chat); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := chat.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.wechat.SendMsg(chat.TOS, chat.Content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "success")
}
