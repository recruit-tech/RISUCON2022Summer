package main

import (
	"flag"

	"github.com/recruit-tech/RISUCON2022Summer/snapshots/generator"
)

func main() {
	var (
		targetUrl string
		minify    bool
	)

	flag.StringVar(&targetUrl, "target", "http://localhost:3000", "target url")
	flag.BoolVar(&minify, "minify", true, "minify output json")
	flag.Parse()

	generator.Run(targetUrl, minify)
}
