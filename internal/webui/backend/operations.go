package backend

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/joyrex2001/nightshift/internal/agent"
	"github.com/joyrex2001/nightshift/internal/scanner"
)

// GetObjects will return the list of currently scanned objects.
func (f *handler) GetObjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(agent.New().GetObjects()); err != nil {
		f.Error(w, r, http.StatusInternalServerError, err)
	}
	return
}

// GetScanners will return the list of active scanners.
func (f *handler) GetScanners(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := []scanner.Config{}
	for _, scnr := range agent.New().GetScanners() {
		res = append(res, scnr.GetConfig())
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		f.Error(w, r, http.StatusInternalServerError, err)
	}
	return
}
