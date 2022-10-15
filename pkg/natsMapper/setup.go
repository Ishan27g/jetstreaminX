package natsMapper

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func init() {

	// Connect Options.
	log.SetFlags(log.LstdFlags)

}

type Jetstream struct {
	JS nats.JetStreamContext
}

func NewJetstream() (*Jetstream, func()) {
	log.Println(os.Getenv("NATS"))
	nc, err := nats.Connect("")
	if err != nil {
		log.Println("error connecting to nats " + err.Error())
		return nil, nil
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("error getting js " + err.Error())
		return nil, nil
	}
	return &Jetstream{JS: js}, nc.Close
}
