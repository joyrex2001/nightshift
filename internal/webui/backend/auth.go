package backend

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (f *handler) Authenticate(okhandler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !f.validUser(r) {
			f.unauthorized(w)
			return
		}
		okhandler(w, r, ps)
		return
	}
}

func (f *handler) unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func (f *handler) validUser(r *http.Request) bool {
	return true
}
