package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-zeromq/zmq4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	endpoint              = "tcp://pubsub.besteffort.ndovloket.nl:7664"
	incomingMessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "zmq_incoming_messages_total",
			Help: "Total number of incoming ZeroMQ messages by envelope",
		},
		[]string{"uri", "envelope"},
	)
)

func init() {
	prometheus.MustRegister(incomingMessagesTotal)
}

func listen() {
	go func() {
		ctx := context.Background()

		sub := zmq4.NewSub(ctx)
		defer sub.Close()

		err := sub.Dial(endpoint)
		if err != nil {
			log.Fatal(err)
		}

		err = sub.SetOption(zmq4.OptionSubscribe, "")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Listening for messages...")

		for {
			msg, err := sub.Recv()
			if err != nil {
				log.Fatal(err)
			}

			var envelope = string(msg.Frames[0])

			incomingMessagesTotal.WithLabelValues(endpoint, envelope).Inc()

			log.Printf("Received message with envelope: %s, endpoint: %s", envelope, endpoint)
		}
	}()
}

func main() {
	listen()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
