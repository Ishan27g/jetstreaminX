package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Ishan27g/jetstreaminX/internal/natsMapper"
	"github.com/Ishan27g/jetstreaminX/pkg/provider/httpHandler"
	"github.com/Ishan27g/jetstreaminX/pkg/provider/httpRequest"
)

const IdentifierHandleSample = "nice"

func init() {
	os.Setenv("NATS", "NATS: nats://localhost:4222")
}

func httpClient(js *natsMapper.Jetstream, endpoint string) {

	r, _ := http.NewRequest("GET", "/unusedUrl/since/url/is/mapped/to/id", nil)

	jsSender := natsMapper.RegisterJSSender(js, endpoint, 10*time.Second)

	rsp := httpRequest.New().Publish(jsSender, r)

	fmt.Println("client : rsp.StatusCode", rsp.StatusCode)

}
func httpHandle(js *natsMapper.Jetstream, endpoint string) {

	jsListener := natsMapper.RegisterJSListener(js, endpoint)

	httpHandler.New().Start(jsListener, func(req *http.Request, rsp *http.Response) {
		rsp.StatusCode = 201
		fmt.Println("http handler called : ", req.URL.String())
	})

}

func main() {
	js, c := natsMapper.NewJetstream()
	defer c()
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "pub":
			httpClient(js, IdentifierHandleSample)
		case "sub":
			httpHandle(js, IdentifierHandleSample)
			<-make(chan bool)
		}
	}
}
