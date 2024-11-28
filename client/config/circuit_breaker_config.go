package config

import (
	"time"

	"github.com/sony/gobreaker/v2"
)

func CircuitBreakerConfig() *gobreaker.CircuitBreaker[[]byte] {
	st := gobreaker.Settings{
		Name:    "server-circuit-breaker",
		Timeout: 5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.TotalFailures >= 3
		},
	}
	return gobreaker.NewCircuitBreaker[[]byte](st)
}
