package httpRequest

import (
	"net/http"

	"github.com/Ishan27g/jetstreaminX/pkg/natsMapper"
	"github.com/Ishan27g/jetstreaminX/pkg/provider"
)

type request struct {
	JS natsMapper.JsSender
}

func (r *request) Publish(req *http.Request) http.Response {
	return r.JS.PublishHTTPRequest(req, 0)
}

func New(js natsMapper.JsSender) provider.Sender[natsMapper.JsSender] {
	return &request{js}
}
