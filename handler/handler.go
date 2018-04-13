package handler

import (
	"github.com/yangbinnnn/messenger/g"
	"github.com/yangbinnnn/messenger/sender"
)

type Handler struct {
	cfg     *g.GlobalConfig
	wechat  *sender.Wechat
	mailcli *sender.MailClient
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
	if h.cfg.Smtp.Enable {
		h.mailcli = sender.NewMailClient(
			h.cfg.Smtp.Addr,
			h.cfg.Smtp.Username,
			h.cfg.Smtp.Password,
			h.cfg.Smtp.From,
			h.cfg.Smtp.Timeout,
			h.cfg.Smtp.TLS,
			false,
		)
	}
}
