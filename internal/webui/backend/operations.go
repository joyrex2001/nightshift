package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/joyrex2001/nightshift/internal/agent"
	"github.com/joyrex2001/nightshift/internal/config"
	"github.com/joyrex2001/nightshift/internal/metrics"
	"github.com/joyrex2001/nightshift/internal/scanner"
)

// GetVersion will version details of nightshift.
func (f *handler) GetVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := struct {
		Version string `json:"version"`
		Build   string `json:"build"`
		Date    string `json:"date"`
	}{
		Version: config.Version,
		Build:   config.Build,
		Date:    config.Date,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		f.Error(w, r, http.StatusInternalServerError, err)
	}
	return
}

// GetObjects will return the list of currently scanned objects.
func (f *handler) GetObjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := []*scanner.Object{}
	for _, obj := range agent.New().GetObjects() {
		if len(obj.Schedule) > 0 {
			res = append(res, obj)
		}
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

// PostObjectsRestore will restore the provided pods to the previous known
// state of the given objects.
func (f *handler) PostObjectsRestore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	in := []*scanner.Object{}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		f.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if err := restoreObjects(in); err != nil {
		f.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

// scaleObjects will scale the array of objects to given amount of replicas.
func scaleObjects(objects []*scanner.Object, replicas int) error {
	errs := []string{}
	metrics.Increase("manual_scale")
	for _, obj := range objects {
		if _err := obj.Scale(replicas); _err != nil {
			errs = append(errs, _err.Error())
		}
	}
	if len(errs) > 0 {
		metrics.Increase("manual_scale_error")
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}
	return nil
}

// restoreObjects will scale the array of objects to the previous known state.
func restoreObjects(objects []*scanner.Object) error {
	errs := []string{}
	metrics.Increase("manual_restore")
	for _, obj := range objects {
		if obj.State == nil {
			_err := fmt.Sprintf("no state available on %s/%s", obj.Namespace, obj.Name)
			errs = append(errs, _err)
			continue
		}
		if obj.State != nil {
			if _err := obj.Scale(obj.State.Replicas); _err != nil {
				errs = append(errs, _err.Error())
			}
		}
	}
	if len(errs) > 0 {
		metrics.Increase("manual_restore_error")
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}
	return nil
}
