package prometheus

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// nolint
const (
	KeyPrefix       = "cassini_"
	KeyQueueSize    = "queue_size"
	KeyQueue        = "queue"
	KeyAdaptors     = "adaptors"
	KeyTxsWait      = "txs_wait"
	KeyTxCost       = "tx_cost"
	KeyTxsPerSecond = "txs_per_second"
	KeyErrors       = "errors"
)

var collector *cassiniCollector

func init() {

	initConfig()

	initCollector()

	// Set(KeyQueueSize, 0,"local")
	Set(KeyQueue, 0)
	// Set(KeyAdaptors, 0, "?")
	Set(KeyTxsWait, 0)
	Set(KeyTxCost, 0)
	Set(KeyTxsPerSecond, 0)
	Set(KeyErrors, 0)

	// go func() {
	// 	t := time.NewTicker(time.Duration(100) * time.Millisecond)
	// 	for {
	// 		select {
	// 		case <-t.C:
	// 			{
	// 				Set(KeyTxsPerSecond, 123)
	// 			}
	// 		}
	// 	}
	// }()
}

func initConfig() {
	initMcs := []*MetricConfig{
		&MetricConfig{
			Key:  KeyQueueSize,
			Type: "ImmutableGaugeMetric"},
		&MetricConfig{
			Key:  KeyErrors,
			Type: "CounterMetric"},
		&MetricConfig{
			Key:  KeyTxsPerSecond,
			Type: "TickerGaugeMetric"}}
	viper.Set(KeyMetricType, initMcs)
}

func initCollector() {
	var mcs []*MetricConfig
	viper.UnmarshalKey(KeyMetricType, &mcs)

	collector = &cassiniCollector{
		metricConfigs: mcs,
		descs:         make(map[string]*prometheus.Desc)}

	collector.descs[KeyQueueSize] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyQueueSize),
		"Size of queue",
		[]string{"type"}, nil)
	collector.descs[KeyQueue] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyQueue),
		"Current size of tx in queue",
		nil, nil)
	collector.descs[KeyTxsPerSecond] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyTxsPerSecond),
		"Number of relayed tx per second",
		nil, nil)
	collector.descs[KeyTxsWait] = prometheus.NewDesc(
		fmt.Sprint(KeyPrefix, KeyTxsWait),
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

// Collector returns a collector
// which exports metrics about status code of network service response
func Collector(errChannel chan<- error) prometheus.Collector {
	collector.SetErrorChannel(errChannel)
	return collector
}

type cassiniCollector struct {
	metricConfigs []*MetricConfig
	descs         map[string]*prometheus.Desc
	mapper        sync.Map
	errChannel    chan<- error
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
			c.errChannel <- fmt.Errorf("%s%s%s",
				"Collect error: can not convert key(",
				key, ") into a string")
			return true
		}
		var metric ExportMetric
		metric, ok = v.(ExportMetric)
		if !ok {
			// var metrics []*CassiniMetric
			// metrics, ok = v.([]*CassiniMetric)
			// if !ok {
			// 	c.ch <- fmt.Errorf("%s%s%s",
			// 		"Collect error: can not convert value(", key,
			// 		") into a *cassiniMetric or a []*cassiniMetric")
			// 	return true
			// }
			// for _, metric = range metrics {
			// 	c.export(ch, key, metric)
			// }
		} else {
			c.export(ch, key, metric)
		}
		return true
	}
	c.mapper.Range(exports)
}

// SetErrorChannel set a channel for error
func (c *cassiniCollector) SetErrorChannel(errChannel chan<- error) {
	c.errChannel = errChannel
}

func (c *cassiniCollector) Set(key string, value float64,
	labelValues ...string) {
	// c.mapper.Store(key, value)
	v, loaded := c.mapper.Load(key)
	if v == nil || !loaded {
		metric := c.createMetric(key)
		if metric == nil {
			return
		}
		metric.SetValue(value)
		metric.SetLabelValues(labelValues)
		if v, loaded = c.mapper.LoadOrStore(key, metric); !loaded {
			return
		}
	}
	metric, ok := v.(ExportMetric)
	if !ok {
		c.errChannel <- fmt.Errorf("%s%s%s",
			"Count error: can not convert value(",
			key, ") into a *cassiniMetric")
		return
	}
	metric.SetValue(value)
}

// Set the value to the collector mapper
func Set(key string, value float64,
	labelValues ...string) {
	collector.Set(key, value, labelValues...)
}

func (c *cassiniCollector) export(ch chan<- prometheus.Metric,
	key string, metric ExportMetric) {
	desc, ok := c.descs[key]
	if !ok {
		c.errChannel <- fmt.Errorf("Collect error: can not find desc(%s)", key)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		desc,
		metric.GetValueType(),
		metric.GetValue(), metric.GetLabelValues()...)
}

func (c *cassiniCollector) createMetric(key string) ExportMetric {
	mc := c.getMetricConfig(key)
	if mc != nil {
		if strings.EqualFold(mc.Type, "ImmutableGaugeMetric") {
			return &ImmutableGaugeMetric{}
		} else if strings.EqualFold(mc.Type, "CounterMetric") {
			return &CounterMetric{}
		} else if strings.EqualFold(mc.Type, "TickerGaugeMetric") {
			m := &TickerGaugeMetric{}
			m.Init()
			return m
		}

	}
	return &GaugeMetric{}
}

func (c *cassiniCollector) getMetricConfig(key string) *MetricConfig {
	for _, mc := range c.metricConfigs {
		if strings.EqualFold(key, mc.Key) {
			return mc
		}
	}
	return nil
}

// StartMetrics prometheus exporter("/metrics") service
func StartMetrics(errChannel chan<- error) {

	prometheus.MustRegister(Collector(errChannel))

	http.Handle("/metrics", promhttp.Handler())
	errChannel <- http.ListenAndServe(":39099", nil)
}
