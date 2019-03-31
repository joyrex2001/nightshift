package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/joyrex2001/nightshift/internal/agent"
	"github.com/joyrex2001/nightshift/internal/scanner"
)

// GetObjects will return the list of currently scanned objects.
func (f *handler) GetObjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := []*scanner.Object{}
	for _, obj := range agent.New().GetObjects() {
		res = append(res, obj)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
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

// PostObjectsScale will scale the provided pods to the number of specified
// replicas.
func (f *handler) PostObjectsScale(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	replicas, err := strconv.Atoi(ps.ByName("replicas"))
	if err != nil {
		f.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if replicas < 0 {
		f.Error(w, r, http.StatusInternalServerError, fmt.Errorf("invalid number of replicas: %d", replicas))
		return
	}
	in := []*scanner.Object{}
	if err = json.NewDecoder(r.Body).Decode(&in); err != nil {
		f.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if err := scaleObjects(in, replicas); err != nil {
		f.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

// scaleObjects will scale the array of objects to given amount of replicas.
func scaleObjects(objects []*scanner.Object, replicas int) error {
	for _, obj := range objects {
		if err := obj.Scale(replicas); err != nil {
			return err
		}
	}
	agent.New().UpdateSchedule()
	return nil
}
