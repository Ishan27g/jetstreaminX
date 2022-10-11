package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Ishan27g/jetstreaminX/pkg/natsMapper"
)

const IdentifierHandleSample = "nice"

func init() {
	os.Setenv("NATS", "NATS: nats://localhost:4222")
}

func httpClient(js *natsMapper.Jetstream, endpoint string) {

	jsSender := natsMapper.RegisterJSSender(js, endpoint, 10*time.Second)
	r, _ := http.NewRequest("GET", "/unusedUrl/since/url/is/mapped/to/id", nil)
	rsp := jsSender.PublishHTTPRequest(r)
	fmt.Println("client : rsp.StatusCode", rsp.StatusCode)
}
func httpHandler(js *natsMapper.Jetstream, endpoint string) {

	jsListener := natsMapper.RegisterJSListener(js, endpoint)
	jsListener.Start(func(req *http.Request, rsp *http.Response) {
		rsp.StatusCode = 200
		fmt.Println("http handler called ")
		fmt.Println("http handler : ", req.URL.String())

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
			httpHandler(js, IdentifierHandleSample)
		}
		<-make(chan bool)
	}
}
