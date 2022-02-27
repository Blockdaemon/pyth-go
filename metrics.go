package pyth

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "pyth"

	subsystemClient = "client"
)

var (
	metricsWsActiveConns = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystemClient,
		Name:      "ws_active_conns",
		Help:      "Number of active WebSockets between Pyth client and RPC nodes",
	})
	metricsWsEventsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemClient,
		Name:      "ws_events_total",
		Help:      "Number of WebSocket events delivered from RPC nodes to Pyth client",
	})
)
