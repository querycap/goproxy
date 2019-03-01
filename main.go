package main

import (
	"net/http"

	"github.com/goproxy/goproxy"
	"github.com/goproxy/goproxy/cacher"
)

func main() {
	g := goproxy.New()

	g.Cacher = &cacher.Disk{Root: "/data"}

	if err := http.ListenAndServe("0.0.0.0:80", g); err != nil {
		panic(err)
	}
}
