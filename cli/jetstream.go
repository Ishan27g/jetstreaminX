package main

import (
	"fmt"
	"log"

	natsMapper2 "github.com/Ishan27g/jetstreaminX/internal/natsMapper"
)

const IdentifierHandleSample = "Okay"

func jsPub(data string) {
	js, c := natsMapper2.NewJetstream()
	defer c()
	st := js.GetOrCreate(IdentifierHandleSample)
	if st == nil {
		log.Println("could not get/create stream ")
		return
	}
	js.Publish(IdentifierHandleSample, natsMapper2.NewMsg(IdentifierHandleSample, []byte(data)).Get())
	fmt.Println("pub: send : " + data)
}
func jsSub() {
	js, c := natsMapper2.NewJetstream()
	defer c()
	st := js.GetOrCreate(IdentifierHandleSample)
	if st == nil {
		log.Println("could not get/create stream ")
		return
	}
	ok := make(chan bool)
	js.Subscribe(IdentifierHandleSample, func(msg *natsMapper2.Message) {
		fmt.Println("sub: received : " + string(msg.Data))
		ok <- true
	})
	<-ok
}
