package prometheus

const (
	// KeyMetricType config key
	KeyMetricType = "metricType"
)

// MetricConfig defined config about metric type
type MetricConfig struct {
	Key  string
	Type string
}
