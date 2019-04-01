package backend

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"

	"github.com/joyrex2001/nightshift/internal/webui/backend/internalfs"
)

func NewHandler() *handler {
	return &handler{}
}

type handler struct {
	once sync.Once
	mux  *httprouter.Router
}

func (f *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.once.Do(f.init)
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	f.mux.ServeHTTP(w, r)
}

func (f *handler) init() {
	// web routing
	f.mux = httprouter.New()
	f.mux.GET("/public/*filepath", f.Authenticate(f.ServeFiles("")))
	f.mux.GET("/api/objects", f.Authenticate(f.GetObjects))
	f.mux.POST("/api/objects/scale/:replicas", f.Authenticate(f.PostObjectsScale))
	f.mux.POST("/api/objects/restore", f.Authenticate(f.PostObjectsRestore))
	f.mux.GET("/api/scanners", f.Authenticate(f.GetScanners))
	f.mux.GET("/healthz", f.Healthz)
	f.mux.GET("/", f.Redirect(307, "/public"))
}

// ServeFiles will host http files based on a filestore as determined by the
// FileStore method for given folder.
func (f *handler) ServeFiles(folder string) httprouter.Handle {
	root := f.FileStore(folder)
	fileServer := http.FileServer(root)
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		req.URL.Path = ps.ByName("filepath")
		fileServer.ServeHTTP(w, req)
	}
}

// FileStore will return a http filesystem object for given folder.
func (f *handler) FileStore(folder string) http.FileSystem {
	return &assetfs.AssetFS{
		Asset:     internalfs.Asset,
		AssetDir:  internalfs.AssetDir,
		AssetInfo: internalfs.AssetInfo,
		Prefix:    folder,
	}
}

// Redirect will redirect the user to given location.
func (f *handler) Redirect(status int, location string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.Redirect(w, r, location, status)
		return
	}
}

// Healthz will return a liveness response.
func (f *handler) Healthz(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{ status: 'OK', timestamp: %d }", time.Now().Unix())
	return
}

// Error will return an error response in json.
func (f *handler) Error(w http.ResponseWriter, r *http.Request, code int, cerr error) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{ status: %d, error: %s }", code, cerr)
	glog.Errorf("HTTP %d: %s", code, cerr)
	return
}
