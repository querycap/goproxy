package main

import (
	"context"
	"github.com/goproxy/goproxy"
	"github.com/goproxy/goproxy/cacher"
	"golang.org/x/mod/module"
	"k8s.io/klog/v2/klogr"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var goprivate = os.Getenv("GOPRIVATE")

func main() {
	g := goproxy.New()
	g.Cacher = &cacher.Disk{Root: "/data"}

	s := &http.Server{}
	s.Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if goprivate != "" && len(req.URL.Path) > 0 && module.MatchPrefixPatterns(goprivate, req.URL.Path[1:]) {
			rw.WriteHeader(http.StatusNotFound)
			_, _ = rw.Write(nil)
			return
		}

		g.ServeHTTP(rw, req)
	})
	s.Addr = "0.0.0.0:80"

	l := klogr.New()

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
