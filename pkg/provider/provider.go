package provider

import (
	"net/http"

	"github.com/Ishan27g/jetstreaminX/pkg/natsMapper"
)

type Sender[s natsMapper.JsSender] interface {
	Publish(*http.Request) http.Response
}

type Listener[l natsMapper.JsListener] interface {
	Start(func(req *http.Request, rsp *http.Response))
}
