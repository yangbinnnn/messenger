package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
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
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}
}
