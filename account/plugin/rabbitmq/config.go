package apprabbitmq

import (
	"encoding/json"
	"log"
)

type PublishConfig struct {
	ExchangeType string
	ExchangeName string
	QueueName    string
	RoutingKey   string

	ExchangeNameDLX string
	RoutingKeyDLX   string
	TTL             int

	Data *Message
}

func (config *PublishConfig) String() string {
	result, err := json.Marshal(config)
	if err != nil {
		log.Fatal("decode:", err)
	}
	return string(result)
}

type QueueDelayExpireConfig struct {
	Message     *Message
	QueueName   string
	ExchangeDLX string
	RouteDLX    string
	TTL         int
}

func (config *QueueDelayExpireConfig) String() string {
	result, err := json.Marshal(config)
	if err != nil {
		log.Fatal("decode:", err)
	}
	return string(result)
}
