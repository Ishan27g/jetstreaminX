package natsMapper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type JsListener interface {
	Start(handlerFunc func(req *http.Request, rsp *http.Response))
}
type JsSender interface {
	PublishHTTPRequest(req *http.Request) http.Response
}

func RegisterJSSender(js *Jetstream, endpointIdentifier string, timeout time.Duration) JsSender {
	return &natsJsHandler{endpointIdentifier: endpointIdentifier, httpTimeout: timeout, Jetstream: js}
}
func RegisterJSListener(js *Jetstream, endpointIdentifier string) JsListener {
	return &natsJsHandler{endpointIdentifier: endpointIdentifier, Jetstream: js}
}

type natsJsHandler struct {
	endpointIdentifier string
	httpTimeout        time.Duration
	*Jetstream
}

func (n *natsJsHandler) Start(handlerFunc func(req *http.Request, rsp *http.Response)) {
	x := NewMsg(n.endpointIdentifier, nil)

	n.GetOrCreate(x.Get().EndpointIdentifier)
	// n.AddSubjects(x.Get().EndpointIdentifier, x.Get().EndpointIdentifier)

	log.Println("listener: subbed to " + x.Get().EndpointIdentifier)
	n.Subscribe(x.Get().EndpointIdentifier, func(msg *Message) {

		go n.Start(handlerFunc)

		var rsp = http.Response{}
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(msg.Data)))
		if err != nil {
			log.Println("listener :" + err.Error())
			return
		}

		handlerFunc(req, &rsp)

		msg.Data, _ = json.Marshal(&HttpResponse{Response: &rsp})
		if err != nil {
			log.Println("listener:" + err.Error())
			return
		}
		// n.AddSubjects(m.EndpointIdentifier, m.Hash)
		s := msg.Hash
		log.Println("listener: pub to " + s)
		n.Publish(s, *msg)
	})
}

func (n *natsJsHandler) PublishHTTPRequest(req *http.Request) http.Response {
	var b = &bytes.Buffer{}
	var rspChan = make(chan http.Response, 1)
	err := req.Write(b)
	if err != nil {
		log.Println("sender:" + err.Error())
		return http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader("sorry"))}
	}

	x := NewMsg(n.endpointIdentifier, b.Bytes())

	if !n.CheckStream(x.Get().EndpointIdentifier) {
		log.Println("no subscribers")
	}

	n.GetOrCreate(x.Get().Hash)
	//n.AddSubjects(x.Get().EndpointIdentifier, x.Get().EndpointIdentifier+"."+x.Get().Hash)

	log.Println("sender: subbed to : ", x.Get().Hash)
	go n.Subscribe(x.Get().Hash, func(natsMsg *Message) {
		var rsp = HttpResponse{}
		err = json.Unmarshal(natsMsg.Data, &rsp)
		if err != nil {
			log.Println("sender:" + err.Error())
			return
		}
		r := *rsp.Response
		rspChan <- r
	})
	s := x.Get().EndpointIdentifier //+ ".*" //+ x.Get().Hash
	log.Println("sender: publish to " + s)
	n.Publish(s, x.Get())

	select {
	case <-time.After(n.httpTimeout):
		log.Println("timed out")
		return http.Response{StatusCode: http.StatusRequestTimeout, Body: io.NopCloser(strings.NewReader("timed out"))}
	case rsp := <-rspChan:
		return rsp
	}
}
