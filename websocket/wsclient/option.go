package wsclient

import (
	"time"
)

// configurableFields for WebsocketClient.
type configurableFields struct {
	ReconnectSleep time.Duration
	PingPeriod     time.Duration
	PongTimeout    time.Duration
}

func defaultFields() *configurableFields {
	return &configurableFields{
		ReconnectSleep: 5 * time.Second,

		// cloud ping pong duration is 40s
		PingPeriod:  10 * time.Second,
		PongTimeout: 15 * time.Second,
	}
}

// Option is the function signature that applies configurable option for WebsocketClient.
type Option func(*configurableFields)

// ReconnectSleep how long it'll sleep before a re-reconnect.
func ReconnectSleep(d time.Duration) Option {
	return func(cf *configurableFields) {
		cf.ReconnectSleep = d
	}
}

// PingPeriod is how often client will handle a ping.
func PingPeriod(d time.Duration) Option {
	return func(cf *configurableFields) {
		cf.PingPeriod = d
	}
}
