package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-zeromq/zmq4"
)

type Listener struct {
	endpoints []string
	ctx       context.Context
}

func NewListener(endpoints []string) *Listener {
	return &Listener{
		endpoints: endpoints,
		ctx:       context.Background(),
	}
}

func (l *Listener) Start() {
	for _, endpoint := range l.endpoints {
		go l.listenToEndpoint(endpoint)
	}
}

func (l *Listener) listenToEndpoint(endpoint string) {
	sub := zmq4.NewSub(l.ctx)
	defer sub.Close()

	err := sub.Dial(endpoint)
	if err != nil {
		log.Fatalf("Failed to dial %s: %v", endpoint, err)
	}

	err = sub.SetOption(zmq4.OptionSubscribe, "")
	if err != nil {
		log.Fatalf("Failed to subscribe to %s: %v", endpoint, err)
	}

	fmt.Printf("Listening for messages on %s...\n", endpoint)

	for {
		msg, err := sub.Recv()
		if err != nil {
			log.Printf("Error receiving from %s: %v", endpoint, err)
			continue
		}

		if len(msg.Frames) == 0 {
			log.Printf("Received empty message from %s", endpoint)
			continue
		}

		envelope := string(msg.Frames[0])

		incomingMessagesTotal.WithLabelValues(endpoint, envelope).Inc()

		log.Printf("Received message with envelope: %s, endpoint: %s", envelope, endpoint)
	}
}
