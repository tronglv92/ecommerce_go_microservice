package kafka

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

type appKafka struct {
	prefix string
	// connection    *amqp.Connection
	logger  logger.Logger
	address string
}

func NewKafka(prefix string) *appKafka {
	return &appKafka{
		prefix: prefix,
	}
}
func (kafka *appKafka) GetPrefix() string {
	return kafka.prefix
}
func (kafka *appKafka) Get() interface{} {
	return kafka
}
func (kafka *appKafka) Name() string {
	return kafka.prefix
}
func (kafka *appKafka) InitFlags() {

	prefix := kafka.prefix
	if kafka.prefix != "" {
		prefix += "-"
	}
	flag.StringVar(&kafka.address, prefix+"address", "", "address of brokers")

}
func (kafka *appKafka) Configure() error {
	kafka.logger = logger.GetCurrent().GetLogger(kafka.prefix)
	return nil
}

func (kafka *appKafka) Run() error {

	return kafka.Configure()
}

func (kafka *appKafka) Stop() <-chan bool {
	c := make(chan bool)

	go func() {

		c <- true
	}()
	return c

}
func (k *appKafka) Publish(ctx context.Context, topic string, messages ...*Message) error {
	address := strings.Split(k.address, ",")
	w := &kafkago.Writer{
		Addr:                   kafkago.TCP(address...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Compression:            kafkago.Snappy,
	}

	messagesKafka := []kafkago.Message{}

	for _, msg := range messages {
		msgkafka := kafkago.Message{
			Key:   []byte(msg.Id),
			Value: msg.EncodeToBytes(),
		}
		messagesKafka = append(messagesKafka, msgkafka)
	}

	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = w.WriteMessages(ctx, messagesKafka...)
		if errors.Is(err, kafkago.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			return err
		}
		break
	}

	if err := w.Close(); err != nil {
		return err
	}
	return nil
}
func (k *appKafka) Subscribe(ctx context.Context, groupId string, topic string) (ch <-chan *Message, close func()) {
	msgCh := make(chan *Message)
	address := strings.Split(k.address, ",")
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  address,
		GroupID:  groupId,
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
	})

	go func() {
		for {

			m, err := r.ReadMessage(ctx)
			if err != nil {
				k.logger.Errorf("subscribe messsage err %v \n", err)
				break
			}
			msg := DecodeToMessage(m.Value)
			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), msg)
			msgCh <- msg
		}
	}()

	return msgCh, func() {
		r.Close()
	}
}
