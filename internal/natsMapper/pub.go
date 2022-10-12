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

type Listener interface {
	Start(handlerFunc func(req *http.Request, rsp *http.Response))
}
type Sender interface {
	PublishHTTPRequest(req *http.Request) http.Response
}

func RegisterSender(endpointIdentifier string, timeout time.Duration) Sender {
	return &natsHandler{endpointIdentifier: endpointIdentifier, httpTimeout: timeout}
}
func RegisterListener(endpointIdentifier string) Listener {
	return &natsHandler{endpointIdentifier: endpointIdentifier}
}

type natsHandler struct {
	endpointIdentifier string
	httpTimeout        time.Duration
}

func (s *natsHandler) Start(handlerFunc func(req *http.Request, rsp *http.Response)) {
	x := NewMsg(s.endpointIdentifier, nil)
	subscribeJson(x.RequestSubject(), func(msg *Message) {
		go s.Start(handlerFunc)

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
		publishJson(msg.ResponseSubject(), *msg, nil)
	})
}

type HttpResponse struct {
	*http.Response `json:"http_response"`
}

func (s *natsHandler) PublishHTTPRequest(req *http.Request) http.Response {

	var b = &bytes.Buffer{}
	var rspChan = make(chan http.Response, 1)
	err := req.Write(b)
	if err != nil {
		log.Println("sender:" + err.Error())
		return http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader("sorry"))}
	}

	m := NewMsg(s.endpointIdentifier, b.Bytes())

	go subscribeJson(m.ResponseSubject(), func(natsMsg *Message) {
		var rsp = HttpResponse{}
		err = json.Unmarshal(natsMsg.Data, &rsp)
		if err != nil {
			log.Println("sender:" + err.Error())
			return
		}
		r := *rsp.Response
		rspChan <- r
	})
	publishJson(m.RequestSubject(), m.Get(), nil)

	select {
	case <-time.After(s.httpTimeout):
		log.Println("timed out")
		return http.Response{StatusCode: http.StatusRequestTimeout, Body: io.NopCloser(strings.NewReader("timed out"))}
	case rsp := <-rspChan:
		return rsp
	}
}
