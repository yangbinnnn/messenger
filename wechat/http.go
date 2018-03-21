package wechat

import (
	"log"
	"net/http"

	"github.com/toolkits/web/param"
	"github.com/yangbinnnn/messenger/config"
)

// ConfigRoute sender/wechat
func ConfigRoute() {
	cfg := config.Config()
	wc := NewWechat(cfg.Wechat.CorpID, cfg.Wechat.AgentId,
		cfg.Wechat.Secret, cfg.Wechat.EncodingAESKEY)
	go wc.GetAccessTokenFromWeixin()
	log.Println("config wechat route")
	http.HandleFunc("/sender/wechat", func(w http.ResponseWriter, r *http.Request) {
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
	})
}
