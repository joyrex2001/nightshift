package webui

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/viper"

	"github.com/joyrex2001/nightswitch/internal/webui/backend"
)

// healthz will start listening for /healthz http requests on given address.
func healthz(addr string) {
	glog.Infof("Starting /healthz on %s...", addr)
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "{ status: 'OK', timestamp: %d }", time.Now().Unix())
		})
		glog.Fatal(http.ListenAndServe(addr, nil))
	}()
}

// Main is the main entry point for starting this service, based the settings
// initiated by cmd.
func Main() {
	hlth := viper.GetString("health.listen-addr")
	if hlth != "" {
		healthz(hlth)
	}

	addr := viper.GetString("web.listen-addr")
	hndlr := backend.NewHandler()
	srv := http.Server{
		Addr:         addr,
		Handler:      backend.HTTPLogger(hndlr),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	cert := viper.GetString("web.cert-file")
	key := viper.GetString("web.key-file")
	tls := viper.GetBool("web.enable-tls")

	glog.Infof("Starting webui on %s...", addr)
	if tls {
		glog.Fatal(srv.ListenAndServeTLS(cert, key))
	} else {
		glog.Fatal(srv.ListenAndServe())
	}
}
