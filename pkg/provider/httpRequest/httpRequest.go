package httpRequest

import (
	"net/http"

	"github.com/Ishan27g/jetstreaminX/internal/natsMapper"
	"github.com/Ishan27g/jetstreaminX/pkg/provider"
)

type request struct{}

func (r *request) Publish(rr natsMapper.JsSender, req *http.Request) http.Response {
	return rr.PublishHTTPRequest(req)
}

func New() provider.Sender[natsMapper.JsSender] {
	return &request{}
}
