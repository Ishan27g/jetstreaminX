package natsMapper

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Jetstream struct {
	JS nats.JetStreamContext
}

func NewJetstream() (*Jetstream, func()) {
	nc, err := nats.Connect(urls)
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

func (j *Jetstream) CheckStream(endpointIdentifier string) bool {
	stInfo := j.get(&endpointIdentifier, nil)
	fmt.Println(stInfo)
	return stInfo != nil
}

func (j *Jetstream) GetOrCreate(endpointIdentifier string) nats.JetStreamContext {
	stInfo := j.get(&endpointIdentifier, nil)

	if stInfo == nil && j.Add(endpointIdentifier) != nil {
		return nil
	}

	return j.JS
}

func (j *Jetstream) get(endpointIdentifier *string, subject *string) *nats.StreamInfo {
	timeOut := make(chan bool)
	go func() {
		<-time.After(3 * time.Second)
		timeOut <- true
		close(timeOut)
	}()
	for {
		stInfo, ok := <-j.JS.Streams()
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
		}
		if _, ok := <-timeOut; !ok {
			fmt.Println("timed out")
			return nil
		}

	}
	return nil
}
func (j *Jetstream) AddSubjects(endpointIdentifier string, subjects ...string) error {

	//stInfo := j.get(&endpointIdentifier, nil)
	//if stInfo != nil {
	//	subjects = append(subjects, stInfo.Config.Subjects...)
	//}
	_, err := j.JS.UpdateStream(j.defaultStreamConfig(endpointIdentifier, subjects))
	if err != nil {
		log.Println("could not update stream " + endpointIdentifier + " : " + err.Error())
		if err.Error() == "nats: duplicate subjects detected" {
			err = nil
		}
	}
	return err
}
func (j *Jetstream) Add(endpointIdentifier string, subjects ...string) error {
	var err error
	_, err = j.JS.AddStream(j.defaultStreamConfig(endpointIdentifier, subjects))
	if err != nil {
		log.Println("could not create stream " + endpointIdentifier + " : " + err.Error())
	}
	return err
}

func (j *Jetstream) defaultStreamConfig(endpointIdentifier string, subjects []string) *nats.StreamConfig {
	return &nats.StreamConfig{
		Name:      endpointIdentifier,
		Subjects:  subjects,
		Retention: nats.WorkQueuePolicy,
		//MaxAge:    250 * time.Millisecond,
	}
}

func (j *Jetstream) Publish(subject string, msg Message) {
	//if j.get(nil, &subject) == nil {
	//	log.Println("pub-stream not found : ", subject)
	//	return
	//}

	b, _ := json.Marshal(msg)

	_, err := j.JS.PublishAsync(subject, b)
	if err != nil {
		return
	}
	select {
	case <-j.JS.PublishAsyncComplete():
	case <-time.After(500 * time.Millisecond):
		fmt.Println("publish did not resolve in time")
	}
	//ack, err := j.JS.PublishMsg(&nats.Msg{
	//	Subject: subject,
	//	Data:    b,
	//})
	if err != nil {
		log.Println("error publishing to jetstream " + err.Error())
		return
	}

	log.Println("good publish to " + subject)
	return

}

func (j *Jetstream) Subscribe(subject string, cb func(msg *Message)) bool {
	//if j.get(nil, &subject) == nil {
	//	log.Println("sub-stream not found :", subject)
	//	return false
	//}
	//var ss *nats.Subscription
	var err error
	var ss *nats.Subscription
	log.Println("subbed to " + subject)
	if ss, err = j.JS.Subscribe(subject, func(natsMsg *nats.Msg) {
		err = natsMsg.Ack()
		if err != nil {
			log.Println("unable to ack natsMsg " + err.Error())
			return
		}
		var msg = &Message{}
		err := json.Unmarshal(natsMsg.Data, msg)
		if err != nil {
			log.Println("unable to unmarshal natsMsg " + err.Error())
			return
		}
		cb(msg)
		if err = ss.Unsubscribe(); err != nil {
			log.Println("unable to unsubscribe " + err.Error())
			return
		}
		//if err = ss.Drain(); err != nil {
		//	log.Println("unable to drain " + err.Error())
		//	return
		//}
	}, nats.AckExplicit()); err != nil {
		log.Println("unable to subscribe " + err.Error())
	}
	//err = ss.AutoUnsubscribe(1)
	//if err != nil {
	//	log.Println("unable to unsubscribe " + err.Error())
	//}
	return true
}
