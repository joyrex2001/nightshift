package backend

// Based on: github.com/gleicon/go-httplogger

import (
	"net/http"
	"time"

	"github.com/golang/glog"
)

type stResponseWriter struct {
	http.ResponseWriter
	HTTPStatus   int
	ResponseSize int
}

func (w *stResponseWriter) WriteHeader(status int) {
	w.HTTPStatus = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *stResponseWriter) Flush() {
	z := w.ResponseWriter
	if f, ok := z.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *stResponseWriter) CloseNotify() <-chan bool {
	z := w.ResponseWriter
	return z.(http.CloseNotifier).CloseNotify()
}

func (w *stResponseWriter) Write(b []byte) (int, error) {
	if w.HTTPStatus == 0 {
		w.HTTPStatus = 200
	}
	w.ResponseSize = len(b)
	return w.ResponseWriter.Write(b)
}

// HTTPLogger middleware will log incoming http requests, unless the request
// path is in the given ignore map.
func HTTPLogger(handler http.Handler, ignore []string) http.Handler {
	skip := map[string]bool{}
	for _, s := range ignore {
		skip[s] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		interceptWriter := stResponseWriter{w, 0, 0}
		handler.ServeHTTP(&interceptWriter, r)
		if nolog, _ := skip[r.URL.Path]; nolog {
			return
		}
		glog.Infof("HTTP - %s - - - \"%s %s %s\" %d %d %s %dus\n",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.Proto,
			interceptWriter.HTTPStatus,
			interceptWriter.ResponseSize,
			r.UserAgent(),
			time.Since(t),
		)
	})
}
