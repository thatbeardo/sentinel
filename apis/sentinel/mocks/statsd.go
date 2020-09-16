package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Statsd is a mocked struct that implements the statsd interface
type Statsd struct {
	ExpectedGaugeMetric  string
	ExpectedTimingMetric string
	T                    *testing.T
}

// Increment mocks the original increment function
func (s Statsd) Increment(metric string) {
}

// Gauge mocks the original gauge function
func (s Statsd) Gauge(metric string, value interface{}) {
	assert.Equal(s.T, s.ExpectedGaugeMetric, metric)
}

// Timing mocks the original gauge function
func (s Statsd) Timing(metric string, value interface{}) {
	assert.Equal(s.T, s.ExpectedTimingMetric, metric)
}
