package webui

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
)

// Main is the main entry point for starting this service, based the settings
// initiated by cmd.
func Main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{ status: 'OK', timestamp: %d }", time.Now().Unix())
	})
	glog.Fatal(http.ListenAndServe(":8088", nil))
}
