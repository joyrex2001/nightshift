package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type mycount struct {
	count int
	add   float64
}

func (m *mycount) Collect(chan<- prometheus.Metric) {}
func (m *mycount) Describe(chan<- *prometheus.Desc) {}
func (m *mycount) Desc() *prometheus.Desc           { return nil }
func (m *mycount) Write(*dto.Metric) error          { return nil }
func (m *mycount) Inc()                             { m.count++ }
func (m *mycount) Add(a float64)                    { m.add += a }

func TestIncrease(t *testing.T) {
	mycounts := map[string]*mycount{}
	for id, cntr := range counters {
		if cntr.prom == nil {
			t.Errorf("failed - prometheus counter object %s is nil", id)
		}
		mycounts[id] = &mycount{}
		cntr.prom = mycounts[id]
	}
	for id := range counters {
		for i := 1; i < 10; i++ {
			Increase(id)
			if mycounts[id].count != i {
				t.Errorf("failed - prometheus counter object %s increase failed, expected %d got %d", id, i, mycounts[id].count)
			}
		}
	}
}
