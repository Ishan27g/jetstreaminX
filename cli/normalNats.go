package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Ishan27g/jetstreaminX/internal/natsMapper"
)

func init() {
	os.Setenv("NATS", "NATS: nats://localhost:4222")
}

func httpClient(endpoint string) {
	s := natsMapper.RegisterSender(endpoint, time.Second*10)
	r, _ := http.NewRequest("GET", "/unusedUrl/since/url/is/mapped/to/id", nil)
	rsp := s.PublishHTTPRequest(r)
	fmt.Println("client : rsp.StatusCode", rsp.StatusCode)
}
func httpHandler(endpoint string) {
	l := natsMapper.RegisterListener(endpoint)
	l.Start(func(req *http.Request, rsp *http.Response) {
		rsp.StatusCode = 200
		fmt.Println("http handler : ", req.URL.String())
	})

}

//
//func main() {
//	//
//	go httpHandler(IdentifierHandleSample)
//
//	<-time.After(1 * time.Second)
//
//	httpClient(IdentifierHandleSample)
//	<-time.After(3000 * time.Second)
//
//	<-make(chan bool)
//}
