package pubsub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Message struct {
	Id string

	Data      map[string]interface{}
	CreatedAt time.Time
	buf       bytes.Buffer
}

func NewMessage(data map[string]interface{}) *Message {
	now := time.Now().UTC()
	buf := bytes.Buffer{}
	return &Message{
		Id:        fmt.Sprintf("%v", now.UnixNano()),
		Data:      data,
		CreatedAt: now,
		buf:       buf,
	}
}

func (evt *Message) String() string {
	return fmt.Sprintf("Message %v", evt.Id)
}

// func (evt *Message) Channel() string {
// 	return evt.channel
// }
// func (evt *Message) SetChannel(channel string) {
// 	evt.channel = channel
// }
// func (evt *Message) Data() map[string]interface{} {
// 	return evt.data
// }
func (evt *Message) EncodeToBytes() []byte {
	result, err := json.Marshal(evt)
	if err != nil {
		log.Fatal("decode:", err)
	}

	return result
}
func DecodeToMessage(data []byte) *Message {

	var msg Message
	json.Unmarshal(data, &msg)
	return &msg

}
