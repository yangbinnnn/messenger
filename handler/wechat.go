package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
	"github.com/toolkits/web/param"
	"github.com/yangbinnnn/messenger/sender"
)

var (
	wc  sender.Wechat
	cfg *g.GlobalConfig
)

func init() {
	cfg = g.Config()
	wc = sender.NewWechat(cfg.Wechat.CorpID, cfg.Wechat.AgentId,
		cfg.Wechat.Secret, cfg.Wechat.EncodingAESKEY)
	go wc.GetAccessTokenFromWeixin()
}

func (h *Handler) SendWeChat(c echo.Context) error {
	if !cfg.Wechat.Enable {
		http.Error(w, "wechat not enable", http.StatusBadRequest)
		return
	}

	token := param.String(r, "token", "")
	if token != cfg.Http.Token {
		http.Error(w, "no privilege", http.StatusForbidden)
		return
	}

	tosStr := param.String(r, "tos", "")
	content := param.String(r, "content", "")
	if tosStr == "" || content == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := wc.SendMsg(tosStr, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write([]byte("success"))
	}
	return
}
