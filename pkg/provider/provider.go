package provider

import (
	"net/http"

	"github.com/Ishan27g/jetstreaminX/internal/natsMapper"
)

type Sender[s natsMapper.JsSender] interface { // HL
	Publish(s, *http.Request) http.Response
}

type Listener[l natsMapper.JsListener] interface { // HL
	Start(l, func(req *http.Request, rsp *http.Response))
}
