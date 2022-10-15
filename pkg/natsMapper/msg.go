package natsMapper

import (
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"
)

type message struct {
	Hash               string `json:"hash"`
	Data               []byte `json:"data"`
	EndpointIdentifier string `json:"endpointIdentifier"`
}

func newMessage(endpointIdentifier string, data []byte) *message {
	var hashIt = func(endpointIdentifier string) (sha string) {
		var randomBytes = func(size int) []byte {
			var b []byte
			for i := 0; i < size; i++ {
				b = append(b, []byte(strconv.Itoa(rand.Intn((i+1)*100)))...)
			}
			return b
		}
		rand.Seed(time.Now().UnixNano())
		hr := sha1.New()
		hr.Write(append([]byte(endpointIdentifier), randomBytes(len(endpointIdentifier))...))
		return base64.URLEncoding.EncodeToString(hr.Sum(nil))
	}
	return &message{
		EndpointIdentifier: endpointIdentifier,
		Data:               data,
		Hash:               hashIt(endpointIdentifier),
	}

}
