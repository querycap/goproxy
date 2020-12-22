package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/goproxy/goproxy"
	"github.com/goproxy/goproxy/cacher"
	"golang.org/x/mod/module"
	"k8s.io/klog/v2/klogr"
)

var goprivate = os.Getenv("GOPRIVATE")
var proxiedsumdbs = os.Getenv("PROXIEDSUMDBS")

func main() {
	g := goproxy.New()

	g.Cacher = &cacher.Disk{Root: "/data"}

	g.ProxiedSUMDBs = func(rules []string) (finalRules []string) {
		for i := range rules {
			v := strings.TrimSpace(rules[i])
			if v != "" {
				finalRules = append(finalRules, v)
			}
		}
		return
	}(strings.Split(proxiedsumdbs, ";"))

	l := klogr.New()

	s := &http.Server{}
	s.Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if len(req.URL.Path) > 0 {
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
