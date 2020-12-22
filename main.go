package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goproxy/goproxy"
	"github.com/goproxy/goproxy/cacher"
	"golang.org/x/mod/module"
	"k8s.io/klog/v2/klogr"
)

var goprivate = os.Getenv("GOPRIVATE")

func main() {
	g := goproxy.New()
	g.Cacher = &cacher.Disk{Root: "/tmp/data"}
	l := klogr.New()

	s := &http.Server{}
	s.Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if len(req.URL.Path) > 0 {
			l.Info(req.RequestURI[1:])

			if goprivate != "" && module.MatchPrefixPatterns(goprivate, req.RequestURI[1:]) {
				rw.WriteHeader(http.StatusNotFound)
				_, _ = rw.Write(nil)
				return
			}
		}

		g.ServeHTTP(rw, req)
	})
	s.Addr = "0.0.0.0:80"

	go func() {
		l.Info("serve on", "addr", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			l.Error(err, "")
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh

	timeout := 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_ = s.Shutdown(ctx)
}
