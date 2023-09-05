package kafka

import (
	"context"
	"errors"
	"flag"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/tronglv92/accounts/component/asyncjob"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

type KafkaConsumerConfig struct {
	GroupId string
	Topic   string
}
type appKafka struct {
	prefix string
	// connection    *amqp.Connection
	logger            logger.Logger
	address           string
	numPartitions     int
	replicationFactor int
	maxRetry          int
	sleepTime         time.Duration
	multiplier        int
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
	flag.IntVar(&kafka.numPartitions, prefix+"number-partitions", 1, "numbers partitions")
	flag.IntVar(&kafka.replicationFactor, prefix+"replication", 1, "replication")
	flag.IntVar(&kafka.maxRetry, prefix+"max-retry", 3, "max retry")
	flag.DurationVar(&kafka.sleepTime, prefix+"sleep-time", 100*time.Millisecond, "sleep time in millisecon")
	flag.IntVar(&kafka.multiplier, prefix+"multiplier", 2, "multiplier")

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
	w := &kafka.Writer{
		Addr:                   kafka.TCP(address...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Compression:            kafka.Snappy,
	}

	messagesKafka := []kafka.Message{}

	for _, msg := range messages {
		msgkafka := kafka.Message{
			Key:   []byte(msg.Id),
			Value: msg.EncodeToBytes(),
		}
		messagesKafka = append(messagesKafka, msgkafka)
	}

	var err error
	retries := k.maxRetry
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = w.WriteMessages(ctx, messagesKafka...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			var sleepTime time.Duration
			if i == 0 {
				sleepTime = k.sleepTime

			} else {

				sleepTime = k.sleepTime * time.Duration(i*k.multiplier)
			}
			time.Sleep(sleepTime)
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
func (k *appKafka) Subscribe(ctx context.Context, config KafkaConsumerConfig) (ch <-chan *Message, close func(), errChan <-chan error) {
	msgCh := make(chan *Message)
	errCh := make(chan error)
	address := strings.Split(k.address, ",")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  address,
		GroupID:  config.GroupId,
		Topic:    config.Topic,
		MaxBytes: 10e6, // 10MB
	})

	go func() {
		for {

			m, err := r.ReadMessage(ctx)
			if err != nil {
				k.logger.Errorf("subscribe messsage err %v \n", err)
				errCh <- err
				break
			}
			msg := DecodeToMessage(m.Value)
			k.logger.Infof("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), msg)
			msgCh <- msg
		}
	}()

	return msgCh, func() {
		r.Close()
	}, errCh
}
func (k *appKafka) CreateTopics(ctx context.Context, topics ...string) error {
	job := asyncjob.NewJob(func(ctx context.Context) error {
		return k.doCreateTopics(ctx, topics...)
	})
	group := asyncjob.NewGroup(false, job)

	if err := group.Run(ctx); err != nil {
		return err
	}
	return nil
}
func (k *appKafka) doCreateTopics(ctx context.Context, topics ...string) error {
	address := strings.Split(k.address, ",")
	if len(address) == 0 {
		return errors.New("must provide address")
	}
	conn, err := kafka.DialContext(ctx, "tcp", address[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{}
	for _, tp := range topics {
		topicConfig := kafka.TopicConfig{

			Topic:             tp,
			NumPartitions:     k.numPartitions,
			ReplicationFactor: k.replicationFactor,
		}
		topicConfigs = append(topicConfigs, topicConfig)
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}
func (k *appKafka) CheckTopicsCreated(ctx context.Context, topicsName ...string) error {
	topics, err := k.getTopics(ctx)
	if err != nil {
		return err
	}

	retryCount := -1
	for _, tpName := range topicsName {

		for !k.isTopicCreated(topics, tpName) {
			retryCount++
			if retryCount >= k.maxRetry {
				k.logger.Infoln("isCheckTopicCreate: {}", false)
				return errors.New("reached max number of retry for reading kafka topic(s)")
			}

			var sleepTime time.Duration
			if retryCount == 0 {
				sleepTime = k.sleepTime
			} else {
				sleepTime = k.sleepTime * time.Duration(retryCount*k.multiplier)
			}
			time.Sleep(sleepTime)
		}
	}
	return nil

}

func (k *appKafka) isTopicCreated(topics map[string]struct{}, topicname string) bool {
	if topics == nil {
		return false
	}
	_, ok := topics[topicname]
	// If the key exists
	return ok
}
func (k *appKafka) getTopics(ctx context.Context) (map[string]struct{}, error) {
	address := strings.Split(k.address, ",")
	if len(address) == 0 {
		return nil, errors.New("must provide address")
	}
	conn, err := kafka.Dial("tcp", address[0])
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, err
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}

	return m, nil
}
