package statsd

import (
	"fmt"

	statsD "gopkg.in/alexcesaro/statsd.v2"
)

// Client is used to connect to a statsd daemon
type Client interface {
	Increment(metric string)
	Gauge(metric string, value interface{})
	Timing(metric string, value interface{})
}

type statsd struct {
	client *statsD.Client
}

func (s statsd) Increment(metric string) {
	s.client.Increment(metric)
}

func (s statsd) Gauge(metric string, value interface{}) {
	s.client.Gauge(metric, value)
}

func (s statsd) Timing(metric string, value interface{}) {
	s.client.Timing(metric, value)
}

// New is a factory method to generate statsd instances
func New(host, port string) (client Client, err error) {
	client, err = statsD.New(statsD.Address(fmt.Sprintf("%s:%s", host, port)))
	if err != nil {
		return nil, err
	}
	return
}
