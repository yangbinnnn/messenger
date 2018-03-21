package mail

import (
	"log"
	"net/http"
	"strings"

	"github.com/toolkits/web/param"
	"github.com/yangbinnnn/messenger/config"
)

// ConfigRoute sender/mail
func ConfigRoute() {
	log.Println("config mail route")
	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		if !cfg.Smtp.Enable {
			http.Error(w, "mail not enable", http.StatusBadRequest)
			return
		}

		client, err := NewClient(
			cfg.Smtp.Addr,
			cfg.Smtp.Username,
			cfg.Smtp.Password,
			cfg.Smtp.From,
			cfg.Smtp.Timeout,
			cfg.Smtp.TLS,
			false,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token := param.String(r, "token", "")
		if token != cfg.Http.Token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tosStr := param.String(r, "tos", "")
		subject := param.String(r, "subject", "")
		content := param.String(r, "content", "")
		if tosStr == "" || subject == "" || content == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		tos := strings.Split(tosStr, ",")
		err = client.Send(tos, subject, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write([]byte("success"))
		}
		return
	})
}
