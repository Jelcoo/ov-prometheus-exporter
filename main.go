package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	incomingMessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ov_zmq_incoming_messages_total",
			Help: "Total number of incoming ZeroMQ messages by envelope",
		},
		[]string{"uri", "envelope"},
	)
)

func init() {
	prometheus.MustRegister(incomingMessagesTotal)
}

func getEndpoints() []string {
	return []string{
		"tcp://pubsub.besteffort.ndovloket.nl:7658", // BISON KV6, KV15, KV17
		"tcp://pubsub.besteffort.ndovloket.nl:7817", // KV78Turbo
		"tcp://pubsub.besteffort.ndovloket.nl:7664", // NS InfoPlus
		"tcp://pubsub.besteffort.ndovloket.nl:7666", // SIRI
	}
}

func main() {
	endpoints := getEndpoints()

	log.Printf("Starting listener for %d endpoint(s)", len(endpoints))
	for _, ep := range endpoints {
		log.Printf("  - %s", ep)
	}

	listener := NewListener(endpoints)
	listener.Start()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Prometheus metrics available at :2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
