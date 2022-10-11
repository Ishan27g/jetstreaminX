package natsMapper

import (
	"encoding/json"
	"testing"
)

func Test_msg_UnmarshalJSON(t *testing.T) {
	type mJSON struct {
		SubjectRequest  string `json:"subjectRequest"`
		SubjectResponse string `json:"subjectResponse"`

		Data []byte `json:"data"`
	}
	m := mJSON{
		SubjectRequest:  "ok",
		SubjectResponse: "ok",
		Data:            []byte("ok"),
	}

	b, err := json.Marshal(m)
	if err != nil {
		t.Error("Message: marshal error")
		return
	}
	var ms = Message{}
	err = json.Unmarshal(b, &ms)
	if err != nil {
		t.Error("Message: unmarshal error")
		return
	}

	msg := NewMsg("ok", nil)
	msg.RequestSubject()
	msg.ResponseSubject()
}
