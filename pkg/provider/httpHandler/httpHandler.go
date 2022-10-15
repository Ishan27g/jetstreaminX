package httpHandler

import (
	"net/http"

	"github.com/Ishan27g/jetstreaminX/pkg/natsMapper"
	"github.com/Ishan27g/jetstreaminX/pkg/provider"
)

type request struct {
	JS natsMapper.JsListener
}

func (r *request) Start(handlerFunc func(req *http.Request, rsp *http.Response)) {
	go r.JS.Start(handlerFunc)
}

func New(js natsMapper.JsListener) provider.Listener[natsMapper.JsListener] {
	return &request{js}
}
