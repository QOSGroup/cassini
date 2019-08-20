package prometheus

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/QOSGroup/cassini/log"
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
		nil, nil)
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
		nil, nil)
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
			log.Error("Collect error: can not convert key into a string")
			return false
		}
		var value float64
		value, ok = v.(float64)
		if !ok {
			log.Error("Collect error: can not convert value into a float64")
			return false
		}
		log.Debugf("collect: %s, %d", key, value)
		var desc *prometheus.Desc
		desc, ok = c.descs[key]
		if !ok {
			log.Errorf("Collect error: can not find desc - %s", key)
			return false
		}
		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			value)
		return true
	}
	c.mapper.Range(exports)
}

func (c *cassiniCollector) Set(key string, value interface{},
	labelValue ...string) {
	c.mapper.Store(key, value)
}

// Set key and value to the collector mapper
func Set(key string, value interface{}, labelValue ...string) {
	collector.Set(key, value, labelValue...)
}

// StartMetrics prometheus exporter("/metrics") service
func StartMetrics() {

	prometheus.MustRegister(Collector())

	http.Handle("/metrics", promhttp.Handler())
	log.Errorf("Prometheus metrics running error: %v",
		http.ListenAndServe(":39099", nil))
}
