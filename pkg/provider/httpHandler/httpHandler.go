package httpHandler

import (
	"net/http"

	"github.com/Ishan27g/jetstreaminX/internal/natsMapper"
	"github.com/Ishan27g/jetstreaminX/pkg/provider"
)

type request struct{}

func (r *request) Start(l natsMapper.JsListener, handlerFunc func(req *http.Request, rsp *http.Response)) {
	go l.Start(handlerFunc)
}

func New() provider.Listener[natsMapper.JsListener] {
	return &request{}
}
