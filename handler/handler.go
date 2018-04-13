package handler

import (
	"github.com/labstack/echo"
	"github.com/yangbinnnn/messenger/g"
	"github.com/yangbinnnn/messenger/sender"
)

type Handler struct {
	cfg    *g.GlobalConfig
	wechat *sender.Wechat
}

func (h *Handler) Prepre() {
	h.cfg = g.Config()
	if h.cfg.Wechat.Enable {
		h.wechat = sender.NewWechat(
			h.cfg.Wechat.CorpID,
			h.cfg.Wechat.AgentId,
			h.cfg.Wechat.Secret,
			h.cfg.Wechat.EncodingAESKEY,
		)
		go h.wechat.GetAccessTokenFromWeixin()
	}
}

func (h *Handler) ParamString(c echo.Context, key string) string {
	value := c.QueryParam(key)
	if value != "" {
		return value
	}
	value = c.FormValue(key)
	return value
}
