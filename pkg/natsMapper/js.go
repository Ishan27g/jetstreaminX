package natsMapper

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func (j *Jetstream) defaultStreamConfig(endpointIdentifier string, subjects []string) *nats.StreamConfig {
	return &nats.StreamConfig{
		Name:     endpointIdentifier,
		Subjects: subjects,
		//		Retention: nats.WorkQueuePolicy,
		Retention:   nats.WorkQueuePolicy,
		MaxAge:      10 * time.Second,
		Description: "stream for " + endpointIdentifier + " http requests and responses",
		Duplicates:  0,
		// MaxConsumers: 1,
	}
}
func (j *Jetstream) checkStream(endpointIdentifier string) bool {
	stInfo := j.get(&endpointIdentifier, nil)
	return stInfo != nil
}

func (j *Jetstream) getOrCreate(endpointIdentifier string) nats.JetStreamContext {
	stInfo := j.get(&endpointIdentifier, nil)

	if stInfo == nil && j.addStream(endpointIdentifier) != nil {
		return nil
	}

	return j.JS
}

func (j *Jetstream) get(endpointIdentifier *string, subject *string) *nats.StreamInfo {
	ctx, can := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer can()
	for {
		stInfo, ok := <-j.JS.Streams(nats.Context(ctx))
		if !ok {
			break
		}
		if endpointIdentifier == nil && subject != nil {
			for _, sub := range stInfo.Config.Subjects {
				if sub == *subject {
					return stInfo
				}
			}
			break
		}
		if subject == nil && endpointIdentifier != nil {
			if stInfo.Config.Name == *endpointIdentifier {
				return stInfo
			}
			break
		}
	}
	return nil
}

// func (j *Jetstream) AddSubjects(endpointIdentifier string, subjects ...string) error {
//
//		_, err := j.JS.UpdateStream(j.defaultStreamConfig(endpointIdentifier, subjects))
//		if err != nil {
//			log.Println("could not update stream " + endpointIdentifier + " : " + err.Error())
//			if err.Error() == "nats: duplicate subjects detected" {
//				err = nil
//			}
//		}
//		return err
//	}
func (j *Jetstream) addStream(endpointIdentifier string, subjects ...string) error {
	var err error
	_, err = j.JS.AddStream(j.defaultStreamConfig(endpointIdentifier, subjects))
	if err != nil {
		log.Println("could not create stream " + endpointIdentifier + " : " + err.Error())
	}
	return err
}
func (j *Jetstream) deleteStream(endpointIdentifier string) error {
	var err error
	err = j.JS.DeleteStream(endpointIdentifier)
	if err != nil {
		log.Println("could not create stream " + endpointIdentifier + " : " + err.Error())
	}
	return err
}

func (j *Jetstream) publish(subject string, msg message, timeout time.Duration) bool {
	b, _ := json.Marshal(msg)

	_, err := j.JS.PublishAsync(subject, b)
	if err != nil {
		log.Println("error in publishing " + err.Error())
		return false
	}
	select {
	case <-j.JS.PublishAsyncComplete():
		log.Println("good publish to " + subject)
		return true
	case <-time.After(timeout):
		log.Println("publish did not resolve in time")
		return false
	}

}

func (j *Jetstream) subscribe(subject string, cb func(msg *message)) bool {
	var err error
	log.Println("subbed to " + subject)
	consumerName := subject + "-consumer"
	j.JS.AddConsumer(subject, &nats.ConsumerConfig{
		Durable:   consumerName,
		AckPolicy: nats.AckExplicitPolicy,
	})
	sub1, err := j.JS.PullSubscribe(subject, consumerName, nats.DeliverAll())
	if err != nil {
		log.Println("unable to pull-sub " + err.Error())
		return false
	}
	natsMsg, err := sub1.Fetch(1, nats.MaxWait(time.Hour))
	if err != nil {
		log.Println("unable to fetch natsMsg " + err.Error())
		return false
	}
	go natsMsg[0].Ack()
	var msg = &message{}
	err = json.Unmarshal(natsMsg[0].Data, msg)
	if err != nil {
		log.Println("unable to unmarshal natsMsg " + err.Error())
		return false
	}
	cb(msg)
	//
	//if _, err = j.JS.Subscribe(subject, func(natsMsg *nats.Msg) {
	//	err = natsMsg.Ack()
	//	if err != nil {
	//		log.Println("unable to ack natsMsg " + err.Error())
	//		return
	//	}
	//	var msg = &message{}
	//	err := json.Unmarshal(natsMsg.Data, msg)
	//	if err != nil {
	//		log.Println("unable to unmarshal natsMsg " + err.Error())
	//		return
	//	}
	//	cb(msg)
	//}, nats.AckExplicit(), nats.DeliverNew()); err != nil {
	//	log.Println("unable to subscribe " + err.Error())
	//}
	return true
}
