package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/yangbinnnn/messenger/g"
	"github.com/yangbinnnn/messenger/handler"
)

// use echo
const (
	VERSION = "0.0.2"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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

	g.Parse(*cfg)

	app := echo.New()
	app.Logger.SetLevel(log.ERROR)
	app.Debug = g.Config().Debug
	app.Use(middleware.Logger())

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	h := &handler.Handler{}
	h.Prepre()

	app.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, VERSION)
	})
	app.GET("/healthy", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	app.Match([]string{"GET", "POST"}, "/sender/mail", h.SendMail)
	app.Match([]string{"GET", "POST"}, "/sender/wechat", h.SendWeChat)
	app.Start(addr)
}
