package natsMapper

import (
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"
)

type Msg interface {
	RequestSubject() string
	ResponseSubject() string

	Get() Message
}
type Message struct {
	Hash               string `json:"hash"`
	Data               []byte `json:"data"`
	EndpointIdentifier string `json:"endpointIdentifier"`
}

func (m *Message) Get() Message {
	return *m
}

func (m *Message) RequestSubject() string {
	s := asSubject("request", m.EndpointIdentifier, "*")
	return s
}

func (m *Message) ResponseSubject() string {
	s := asSubject("response", m.EndpointIdentifier, m.Hash)
	return s
}

func asSubject(pre, str, hash string) string {
	return "endpoint." + pre + "." + str + "." + hash
}

func NewMsg(endpointIdentifier string, data []byte) Msg {
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
	return &Message{
		EndpointIdentifier: endpointIdentifier,
		Data:               data,
		Hash:               hashIt(endpointIdentifier),
	}

}
