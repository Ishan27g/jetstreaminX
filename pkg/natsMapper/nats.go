package natsMapper

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/nats-io/nats.go"
)

var urls = os.Getenv("NATS")

func init() {

	// Connect Options.
	log.SetFlags(log.LstdFlags)

}

func publishJson(subj string, msg Message, reply *string) bool {

	nc, err := nats.Connect(urls)
	if err != nil {
		log.Println(err)
		return false
	}

	defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Println(err)
		return false
	}
	defer ec.Close()
	if reply != nil && *reply != "" {
		err = ec.PublishRequest(subj, *reply, &msg)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}
	if err := ec.Publish(subj, &msg); err != nil {
		log.Println(err)
		return false
	}

	err = ec.Flush()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if err := ec.LastError(); err != nil {
		log.Fatal(err)
	} else {
		// log.Printf("Published  on [%s] \n", subj)
	}
	return true
}

func subscribeJson(subj string, cb func(msg *Message)) {
	var opts []nats.Option

	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
		if s != nil {
			log.Printf("Async error in %q/%q: %v", s.Subject, s.Queue, err)
		} else {
			log.Printf("Async error outside subscription: %v", err)
		}
	}))
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Println(err)
		return
	}
	defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)

	if _, err := ec.Subscribe(subj, func(s *Message) {
		cb(s)
		wg.Done()
	}); err != nil {
		log.Fatal(err.Error())
	}

	//log.Printf("Listening on [%s]", subj)
	wg.Wait()

	if err := ec.LastError(); err != nil {
		log.Fatal(err)
	}
}
