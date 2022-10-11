package main

import (
	"fmt"
	"log"

	"github.com/Ishan27g/jetstreaminX/pkg/natsMapper"
)

const IdentifierHandleSample = "Okay"

func jsPub(data string) {
	js, c := natsMapper.NewJetstream()
	defer c()
	st := js.GetOrCreate(IdentifierHandleSample)
	if st == nil {
		log.Println("could not get/create stream ")
		return
	}
	js.Publish(IdentifierHandleSample, natsMapper.NewMsg(IdentifierHandleSample, []byte(data)).Get())
	fmt.Println("pub: send : " + data)
}
func jsSub() {
	js, c := natsMapper.NewJetstream()
	defer c()
	st := js.GetOrCreate(IdentifierHandleSample)
	if st == nil {
		log.Println("could not get/create stream ")
		return
	}
	ok := make(chan bool)
	js.Subscribe(IdentifierHandleSample, func(msg *natsMapper.Message) {
		fmt.Println("sub: received : " + string(msg.Data))
		ok <- true
	})
	<-ok
}
