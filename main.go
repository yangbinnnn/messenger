package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/yangbinnnn/messenger/config"
	"github.com/yangbinnnn/messenger/mail"
	"github.com/yangbinnnn/messenger/wechat"
)

const (
	VERSION = "0.0.1"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	prepare()

	version := flag.Bool("v", false, "show version")
	cfg := flag.String("c", "cfg.json", "config path")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	config.Parse(*cfg)

	mail.ConfigRoute()
	wechat.ConfigRoute()

	addr := config.Config().Http.Listen
	if addr == "" {
		return
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(VERSION))
	})

	log.Println("listen on", addr)
	http.ListenAndServe(addr, nil)
}
