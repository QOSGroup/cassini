package prometheus

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// nolint
const (
	KeyPrefix    = "cassini_"
	KeyQueueSize = "queue_size"
	KeyQueue     = "queue"
	KeyTxs       = "txs"
	KeyTxWait    = "tx_wait"
	KeyTxCost    = "tx_cost"
	KeyErrors    = "errors"
	KeyAdaptors  = "adaptors"
)

var collector *cassiniCollector

func init() {
	collector = &cassiniCollector{
		descs: make(map[string]*prometheus.Desc)}

	collector.descs[KeyQueueSize] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyQueueSize),
		"Size of queue",
		[]string{"type"}, nil)
	collector.descs[KeyQueue] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyQueue),
		"Current size of tx in queue",
		nil, nil)
	collector.descs[KeyTxs] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyTxs),
		"Number of relayed tx last minute",
		nil, nil)
	collector.descs[KeyTxWait] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyTxWait),
		"Number of tx waiting to be relayed",
		nil, nil)
	collector.descs[KeyTxCost] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyTxCost),
		"Time(milliseconds) cost of lastest tx relay",
		nil, nil)
	collector.descs[KeyAdaptors] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyAdaptors),
		"Number of available adaptors",
		[]string{"node"}, nil)
	// []string{"from", "to"}, nil)
	collector.descs[KeyErrors] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyErrors),
		"Count of running errors",
		nil, nil)
}

type cassiniCollector struct {
	descs map[string]*prometheus.Desc

	mapper sync.Map
}

// CassiniMetric wraps prometheus export data
type CassiniMetric struct {
	value       float64
	LabelValues []string
	mux         sync.RWMutex
}

// Value returns the metric's value
func (m *CassiniMetric) Value() float64 {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.value
}

// Set the metric's value
func (m *CassiniMetric) Set(v float64) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.value = v
}

// Count the value to the collector mapper
func (m *CassiniMetric) Count(increase float64) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.value += increase
}

// Collector returns a collector
// which exports metrics about status code of network service response
func Collector() prometheus.Collector {
	return collector
}

// Describe returns all descriptions of the collector.
func (c *cassiniCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.descs {
		ch <- desc
	}
}

// Collect returns the current state of all metrics of the collector.
func (c *cassiniCollector) Collect(ch chan<- prometheus.Metric) {
	exports := func(k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			panic(errors.New("Collect error: can not convert key into a string"))
			// return false
		}
		var metric *CassiniMetric
		metric, ok = v.(*CassiniMetric)
		if !ok {
			var metrics []*CassiniMetric
			metrics, ok = v.([]*CassiniMetric)
			if !ok {
				panic(fmt.Errorf("%s %s %s", "Collect error:",
					"can not convert value into a *cassiniMetric",
					"or a []*cassiniMetric"))
				// return false
			}
			for _, metric = range metrics {
				c.export(ch, key, metric)
			}
		} else {
			c.export(ch, key, metric)
		}
		return true
	}
	c.mapper.Range(exports)
}

func (c *cassiniCollector) export(ch chan<- prometheus.Metric,
	key string, metric *CassiniMetric) {
	desc, ok := c.descs[key]
	if !ok {
		panic(fmt.Errorf("Collect error: can not find desc - %s", key))
	}
	ch <- prometheus.MustNewConstMetric(
		desc,
		prometheus.GaugeValue,
		metric.Value(), metric.LabelValues...)
}

func (c *cassiniCollector) Set(key string, value interface{}) {
	c.mapper.Store(key, value)
}

func (c *cassiniCollector) Count(key string, increase float64) {
	v, loaded := c.mapper.Load(key)
	if v == nil || !loaded {
		metric := &CassiniMetric{
			value: float64(increase)}
		if v, loaded = c.mapper.LoadOrStore(key, metric); !loaded {
			return
		}
	}
	metric, ok := v.(*CassiniMetric)
	if !ok {
		panic(fmt.Errorf("Count error: can not convert value into a *cassiniMetric"))
		// return
	}
	metric.Count(increase)
}

// Set the value to the collector mapper
func Set(key string, value interface{}) {
	collector.Set(key, value)
}

// Count the value to the collector mapper
func Count(key string, increase float64) {
	collector.Count(key, increase)
}

// StartMetrics prometheus exporter("/metrics") service
func StartMetrics() {

	prometheus.MustRegister(Collector())

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":39099", nil)
	panic(fmt.Errorf("Prometheus metrics running error: %v", err))
}
