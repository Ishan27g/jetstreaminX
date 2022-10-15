package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Ishan27g/jetstreaminX/pkg/natsMapper"
	"github.com/Ishan27g/jetstreaminX/pkg/provider/httpHandler"
	"github.com/Ishan27g/jetstreaminX/pkg/provider/httpRequest"
	"golang.org/x/sync/errgroup"
)

const IdentifierHandleSample = "nice"

func httpClient(js *natsMapper.Jetstream, endpoint string) {

	r, _ := http.NewRequest("GET", "/unusedUrl/since/url/is/mapped/to/id", nil)

	jsSender := natsMapper.RegisterJSSender(js, endpoint, 8*time.Second)

	s := httpRequest.New(jsSender)
	rand.Seed(time.Now().UnixNano())

	rspCodes := make(chan map[int]int, 1)
	rspCodes <- map[int]int{}
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(5)
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			<-time.After(time.Duration(rand.Intn(1000)) * time.Millisecond)
			rsp := s.Publish(r)
			log.Println("client : rsp.StatusCode", rsp.StatusCode)
			r := <-rspCodes
			r[rsp.StatusCode]++
			rspCodes <- r
			return nil
		})
	}
	g.Wait()
	close(rspCodes)
	for code := range rspCodes {
		log.Println(code)
	}

}
func httpHandle(js *natsMapper.Jetstream, endpoint string) {

	jsListener := natsMapper.RegisterJSListener(js, endpoint)

	httpHandler.New(jsListener).Start(func(req *http.Request, rsp *http.Response) {
		rsp.StatusCode = 201
		log.Println("http handler called : ", req.URL.String())
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
