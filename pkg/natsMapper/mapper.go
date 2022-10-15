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

type Response struct {
	*http.Response `json:"http_response"`
}
type JsListener interface {
	Start(handlerFunc func(req *http.Request, rsp *http.Response))
}
type JsSender interface {
	PublishHTTPRequest(req *http.Request, retry int) http.Response
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
	x := newMessage(n.endpointIdentifier, nil)

	n.getOrCreate(x.EndpointIdentifier)

	log.Println("listener: subbed to " + x.EndpointIdentifier)
	n.subscribe(x.EndpointIdentifier, func(msg *message) {

		var rsp = http.Response{}
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(msg.Data)))
		if err != nil {
			log.Println("listener :" + err.Error())
			return
		}

		handlerFunc(req, &rsp)

		msg.Data, _ = json.Marshal(&Response{Response: &rsp})
		if err != nil {
			log.Println("listener:" + err.Error())
			return
		}
		// n.AddSubjects(m.EndpointIdentifier, m.Hash)
		s := msg.Hash
		log.Println("listener: pub to " + s)
		n.publish(s, *msg)
		n.deleteStream(s)
	})
}

func (n *natsJsHandler) PublishHTTPRequest(req *http.Request, retry int) http.Response {
	var b = &bytes.Buffer{}
	var rspChan = make(chan http.Response, 1)
	err := req.Write(b)
	if err != nil {
		log.Println("sender:" + err.Error())
		return http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader("sorry"))}
	}

	x := newMessage(n.endpointIdentifier, b.Bytes())

	if !n.checkStream(x.EndpointIdentifier) {
		log.Println("no subscribers")
		if retry >= 2 {
			return http.Response{StatusCode: http.StatusBadGateway, Body: io.NopCloser(strings.NewReader("No handlers subscribed for - " + n.endpointIdentifier))}
		}
		<-time.After(100 * time.Millisecond)
		n.PublishHTTPRequest(req, retry+1)
	}

	n.getOrCreate(x.Hash)

	log.Println("sender: subbing to : ", x.Hash)
	go n.subscribe(x.Hash, func(natsMsg *message) {

		var rsp = Response{}
		err = json.Unmarshal(natsMsg.Data, &rsp)
		if err != nil {
			log.Println("sender:" + err.Error())
			return
		}
		r := *rsp.Response
		rspChan <- r
	})
	s := x.EndpointIdentifier //+ ".*" //+ x.Hash
	log.Println("sender: publishing to " + s)
	<-time.After(300 * time.Millisecond)
	n.publish(s, *x)

	select {
	case <-time.After(n.httpTimeout):
		log.Println("timed out")
		return http.Response{StatusCode: http.StatusRequestTimeout, Body: io.NopCloser(strings.NewReader("timed out"))}
	case rsp := <-rspChan:
		return rsp
	}
}
