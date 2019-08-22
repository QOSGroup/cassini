package prometheus

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ExportMetric defined an interface about prometheus exporter metric
type ExportMetric interface {
	GetValue() float64
	SetValue(float64)
	GetValueType() prometheus.ValueType
	SetValueType(prometheus.ValueType)
	GetLabelValues() []string
	SetLabelValues([]string)
}

// ImmutableGaugeMetric stores an immutable value
type ImmutableGaugeMetric struct {
	value       float64
	labelValues []string
}

// GetValue returns value
func (m *ImmutableGaugeMetric) GetValue() float64 {
	return m.value
}

// SetValue sets value
func (m *ImmutableGaugeMetric) SetValue(v float64) {
	m.value = v
}

// GetValueType returns value type
func (m *ImmutableGaugeMetric) GetValueType() prometheus.ValueType {
	return prometheus.GaugeValue
}

// SetValueType sets nothing
func (m *ImmutableGaugeMetric) SetValueType(prometheus.ValueType) {}

// GetLabelValues returns label values
func (m *ImmutableGaugeMetric) GetLabelValues() []string {
	return m.labelValues
}

// SetLabelValues sets label values
func (m *ImmutableGaugeMetric) SetLabelValues(values []string) {
	m.labelValues = values
}

// CounterMetric stores a counter value
type CounterMetric struct {
	value       float64
	labelValues []string
	mux         sync.RWMutex
}

// GetValue returns value
func (m *CounterMetric) GetValue() float64 {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.value
}

// SetValue sets value
func (m *CounterMetric) SetValue(v float64) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.value += v
}

// GetValueType returns value type
func (m *CounterMetric) GetValueType() prometheus.ValueType {
	return prometheus.CounterValue
}

// SetValueType sets nothing
func (m *CounterMetric) SetValueType(prometheus.ValueType) {}

// GetLabelValues returns label values
func (m *CounterMetric) GetLabelValues() []string {
	return m.labelValues
}

// SetLabelValues sets label values
func (m *CounterMetric) SetLabelValues(values []string) {
	m.labelValues = values
}

// GaugeMetric wraps prometheus export data
type GaugeMetric struct {
	value       float64
	labelValues []string
	mux         sync.RWMutex
}

// GetValue returns value
func (m *GaugeMetric) GetValue() float64 {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.value
}

// SetValue sets value
func (m *GaugeMetric) SetValue(v float64) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.value += v
}

// GetValueType returns value type
func (m *GaugeMetric) GetValueType() prometheus.ValueType {
	return prometheus.GaugeValue
}

// SetValueType sets nothing
func (m *GaugeMetric) SetValueType(prometheus.ValueType) {}

// GetLabelValues returns label values
func (m *GaugeMetric) GetLabelValues() []string {
	return m.labelValues
}

// SetLabelValues sets label values
func (m *GaugeMetric) SetLabelValues(values []string) {
	m.labelValues = values
}

// TickerGaugeMetric wraps prometheus export data with ticker job
type TickerGaugeMetric struct {
	calculating float64
	value       float64
	labelValues []string
	mux         sync.RWMutex
}

// Init starts ticker goroutine
func (m *TickerGaugeMetric) Init() {
	t := time.NewTicker(time.Duration(1) * time.Second)

	go func() {
		for {
			select {
			case <-t.C:
				{
					m.SetValue(0)
				}
			}
		}
	}()
}

// GetValue returns value
func (m *TickerGaugeMetric) GetValue() float64 {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.value
}

// SetValue sets value
func (m *TickerGaugeMetric) SetValue(v float64) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if v == 0 {
		m.value = m.calculating
		m.calculating = 0
		return
	}
	m.calculating += v
}

// GetValueType returns value type
func (m *TickerGaugeMetric) GetValueType() prometheus.ValueType {
	return prometheus.GaugeValue
}

// SetValueType sets nothing
func (m *TickerGaugeMetric) SetValueType(prometheus.ValueType) {}

// GetLabelValues returns label values
func (m *TickerGaugeMetric) GetLabelValues() []string {
	return m.labelValues
}

// SetLabelValues sets label values
func (m *TickerGaugeMetric) SetLabelValues(values []string) {
	m.labelValues = values
}
