package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/yangbinnnn/messenger/g"
)

type Chat struct {
	Token   string `json:"token" form:"token" query:"token"`
	TOS     string `json:"tos" form:"tos" query:"tos"`
	Content string `json:"content" form:"content" query:"content"`
	ChatID  string `json:"chatid" form:"chatid" query:"chatid"`
}

type Group struct {
	Name   string `json:"name" form:"name" query:"name"`
	Users  string `json:"users" form:"users" query:"users"`
	ChatID string `json:"chatid" form:"chatid" query:"chatid"`
}

func (c Chat) Validate() error {
	if c.Token != g.Config().Http.Token {
		return echo.ErrForbidden
	}

	if c.TOS == "" && c.ChatID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "tos or chatid requried")
	}

	if c.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "content requried")
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

	var err error
	if chat.ChatID != "" {
		err = h.wechat.SendMsgToGroup(chat.ChatID, chat.Content)
	} else {
		err = h.wechat.SendMsg(chat.TOS, chat.Content)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

func (h *Handler) NewChatGroup(c echo.Context) error {
	if !g.Config().Wechat.Enable {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	group := new(Group)
	if err := c.Bind(group); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if group.Name == "" || group.Users == "" || group.ChatID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "chatid or name or users required")
	}

	userList := strings.Split(group.Users, ",")
	if len(userList) < 2 {
		return echo.NewHTTPError(http.StatusBadRequest, "at least two users required")
	}

	err := h.wechat.CreateChatGroup(group.Name, userList[0], group.ChatID, userList)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "success")
}
