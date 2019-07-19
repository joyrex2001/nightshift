package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const metricsPrefix = "nightshift_"

var (
	counters = map[string]*struct {
		Name string
		Help string
		prom prometheus.Counter
	}{
		"scale": {
			Help: "The total number of processed scale events",
		},
		"scale_error": {
			Help: "The total number errors while scaling",
		},
		"manual_scale": {
			Help: "The total number of processed manual scale events",
		},
		"manual_scale_error": {
			Help: "The total number of errors while manual scaling",
		},
		"manual_restore": {
			Help: "The total number of processed manual restore events",
		},
		"manual_restore_error": {
			Help: "The total number of errors while manual restoring",
		},
		"resync_error": {
			Help: "The total number errors while resyncing objects",
		},
		"watch_retries": {
			Help: "The total number of watcher connection retries",
		},
		"watch_event_error": {
			Help: "The total number of error events received from watcher connection",
		},
	}
	// custom metric for exporting current number of replicas
	replicas = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricsPrefix + "replicas",
			Help: "Current expected number of nightshift scaled replicas",
		},
		[]string{"target", "scanner"},
	)
)

func init() {
	for id, m := range counters {
		m.prom = prometheus.NewCounter(prometheus.CounterOpts{
			Name: metricsPrefix + id,
			Help: m.Help,
		})
		prometheus.MustRegister(m.prom)
	}
	prometheus.MustRegister(replicas)
}

// Increase will increase given metric with 1
func Increase(metr string) {
	prom, ok := counters[metr]
	if ok && prom.prom != nil {
		prom.prom.Inc()
	}
}

// SetReplicas will set the replicas metric to given value for given namespace
// and scanner id.
func SetReplicas(ns, scanid string, repl int) {
	if scanid == "" {
		// don't store metrics without scanner id's configured.
		return
	}
	replicas.With(prometheus.Labels{
		"target":  ns,
		"scanner": scanid}).Set(float64(repl))
}
