package apprabbitmq

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"

	"github.com/tronglv92/accounts/common"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

type PublishConfig struct {
	ExchangeType string
	ExchangeName string
	QueueName    string
	RoutingKey   string

	ExchangeNameDLX string
	QueueNameDLX    string
	RoutingKeyDLX   string

	Data *Message
}
type appRabbitMQ struct {
	prefix string
	// connection    *amqp.Connection
	logger logger.Logger
	url    string

	reliable      bool
	continuous    bool
	autoAck       bool
	verbose       bool
	deliveryCount int
	done          chan bool
	errChan       chan error
}

func NewRabbitMQ(prefix string) *appRabbitMQ {
	return &appRabbitMQ{
		prefix: prefix,
	}
}
func (rabbit *appRabbitMQ) GetPrefix() string {
	return rabbit.prefix
}
func (rabbit *appRabbitMQ) Get() interface{} {
	return rabbit
}
func (rabbit *appRabbitMQ) Name() string {
	return rabbit.prefix
}
func (rabbit *appRabbitMQ) InitFlags() {
	flag.StringVar(&rabbit.url, rabbit.prefix+"-url", "amqp://test:dogcute@localhost:5672/", "URL of RabbitMQ service")

	// flag.StringVar(&rabbit.exchangeType, rabbit.prefix+"-exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")

	flag.BoolVar(&rabbit.reliable, rabbit.prefix+"-reliable", true, "Wait for the publisher confirmation before exiting")
	flag.BoolVar(&rabbit.continuous, rabbit.prefix+"-continuous", false, "Keep publishing messages at a 1msg/sec rate")
	flag.BoolVar(&rabbit.autoAck, "auto_ack", false, "enable message auto-ack")
	flag.BoolVar(&rabbit.verbose, "verbose", true, "enable verbose output of message data")

}
func (rabbit *appRabbitMQ) Configure() error {
	rabbit.logger = logger.GetCurrent().GetLogger(rabbit.prefix)
	return nil
}

func (rabbit *appRabbitMQ) Run() error {
	return rabbit.Configure()
}

func (rabbit *appRabbitMQ) Stop() <-chan bool {
	c := make(chan bool)
	// rabbit.done = make(chan bool)
	go func() {
		// if rabbit.connection != nil {
		// 	rabbit.connection.Close()
		// }
		// rabbit.done <- true
		c <- true
	}()
	return c

}

// func (rabbit *appRabbitMQ) Publish(nameConnection string, done chan bool,
//
//	amqpURI string, exchange string, exchangeType string, routingKey string, body string, reliable bool)
func (rabbit *appRabbitMQ) Publish(ctx context.Context, publishConfig PublishConfig) error {
	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", publishConfig.QueueName))
	defer messageSpan.End()

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-producer")
	rabbit.logger.Infof("dialing %q", rabbit.url)
	connection, err := amqp.DialConfig(rabbit.url, config)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}
	defer connection.Close()

	rabbit.logger.Infof("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}
	defer channel.Close()
	rabbit.logger.Infof("got Channel, declaring %q Exchange (%q)", publishConfig.ExchangeType, publishConfig.ExchangeName)

	if err := channel.ExchangeDeclare(
		publishConfig.ExchangeName, // name
		publishConfig.ExchangeType, // type
		true,                       // durable
		false,                      // auto-deleted
		false,                      // internal
		false,                      // noWait
		nil,                        // arguments
	); err != nil {
		return fmt.Errorf("exchange Declare: %s", err)
	}

	if len(publishConfig.QueueName) > 0 {
		rabbit.logger.Infof("producer: declaring queue '%s'", publishConfig.QueueName)

		queue, err := channel.QueueDeclare(
			publishConfig.QueueName, // name of the queue
			true,                    // durable
			false,                   // delete when unused
			false,                   // exclusive
			false,                   // noWait
			nil,                     // arguments
		)
		if err == nil {
			rabbit.logger.Infof("producer: declared queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
				queue.Name, queue.Messages, queue.Consumers, publishConfig.RoutingKey)
		} else {
			return fmt.Errorf("producer: Queue Declare: %s", err)
		}

		log.Printf("producer: declaring binding")
		if err := channel.QueueBind(queue.Name, publishConfig.RoutingKey, publishConfig.ExchangeName, false, nil); err != nil {
			return fmt.Errorf("producer: Queue Bind: %s", err)
		}
	}

	var publishes chan uint64 = nil
	var confirms chan amqp.Confirmation = nil

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if rabbit.reliable {
		rabbit.logger.Infof("enabling publisher confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %s", err)
		}
		// We'll allow for a few outstanding publisher confirms
		publishes = make(chan uint64, 8)
		confirms = channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		go rabbit.confirmHandler(rabbit.done, publishes, confirms)
	}

	rabbit.logger.Infof("declared Exchange, publishing messages")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for {
		seqNo := channel.GetNextPublishSeqNo()
		rabbit.logger.Infof("publishing B body (%q)", publishConfig.Data)

		if err := channel.PublishWithContext(ctx,
			publishConfig.ExchangeName, // publish to an exchange
			publishConfig.RoutingKey,   // routing to 0 or more queues
			false,                      // mandatory
			false,                      // immediate
			amqp.Publishing{
				Headers:         common.InjectAMQPHeaders(amqpContext),
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            publishConfig.Data.EncodeToBytes(),
				DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
				Priority:        0,              // 0-9
				// a bunch of application/implementation-specific fields
			},
		); err != nil {
			return fmt.Errorf("exchange Publish: %s", err)
		}

		rabbit.logger.Infof("published %dB OK", publishConfig.Data)
		if rabbit.reliable {
			rabbit.logger.Infof("vao trong nay 1")
			publishes <- seqNo
		}

		if rabbit.continuous {
			rabbit.logger.Infof("vao trong nay 2")
			select {
			case <-rabbit.done:
				rabbit.logger.Infof("producer is stopping")
				return nil
			case <-time.After(time.Second):
				continue
			}
		} else {
			break
		}
	}

	return nil
}

func (rabbit *appRabbitMQ) confirmHandler(done chan bool, publishes chan uint64, confirms chan amqp.Confirmation) {
	m := make(map[uint64]bool)
	for {
		select {
		case <-done:
			rabbit.logger.Infof("confirmHandler is stopping")
			return
		case publishSeqNo := <-publishes:
			rabbit.logger.Infof("waiting for confirmation of %d", publishSeqNo)
			m[publishSeqNo] = false
		case confirmed := <-confirms:
			if confirmed.DeliveryTag > 0 {
				if confirmed.Ack {
					rabbit.logger.Infof("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
				} else {
					rabbit.logger.Errorf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
				}
				delete(m, confirmed.DeliveryTag)
			}
		}
		if len(m) > 1 {
			rabbit.logger.Infof("outstanding confirmations: %d", len(m))
		}
	}
}
func (rabbit *appRabbitMQ) Subscribe(exchangeType string, exchangeName string, queueName string, routingKey string) (ch <-chan *Message, close func()) {
	msgChan := make(chan *Message)

	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-consumer")
	rabbit.logger.Infof("dialing %q", rabbit.url)
	connection, err := amqp.DialConfig(rabbit.url, config)
	if err != nil {
		rabbit.logger.Errorf("Dial: %s", err)
	}

	go func() {
		rabbit.logger.Infof("closing: %s", <-connection.NotifyClose(make(chan *amqp.Error)))
	}()

	rabbit.logger.Infof("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		rabbit.logger.Errorf("Channel: %s", err)

	}

	rabbit.logger.Infof("got Channel, declaring Exchange (%q)", exchangeName)
	if err = channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		rabbit.logger.Errorf("Exchange Declare: %s", err)

	}

	rabbit.logger.Infof("declared Exchange, declaring Queue %q", queueName)
	queue, err := channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		rabbit.logger.Errorf("Queue Declare: %s", err)

	}

	rabbit.logger.Infof("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, routingKey)

	if err = channel.QueueBind(
		queue.Name,   // name of the queue
		routingKey,   // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		rabbit.logger.Errorf("Queue Bind: %s", err)

	}

	rabbit.logger.Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", "tag")
	deliveries, err := channel.Consume(
		queueName,      // name
		"tag",          // consumerTag,
		rabbit.autoAck, // autoAck
		false,          // exclusive
		false,          // noLocal
		false,          // noWait
		nil,            // arguments
	)
	if err != nil {
		rabbit.logger.Errorf("Queue Consume: %s", err)

	}

	go rabbit.handle(msgChan, deliveries, rabbit.errChan)
	return msgChan, func() {
		connection.Close()
	}
}
func (rabbit *appRabbitMQ) Consumer(exchangeType string, exchangeName string, queueName string, routingKey string, exchangeNameDLX string, routingKeyDLX string) (ch <-chan *Message, close func()) {
	msgChan := make(chan *Message)

	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-consumer")
	rabbit.logger.Infof("dialing %q", rabbit.url)
	connection, err := amqp.DialConfig(rabbit.url, config)
	if err != nil {
		rabbit.logger.Errorf("Dial: %s", err)
	}

	go func() {
		rabbit.logger.Infof("closing: %s", <-connection.NotifyClose(make(chan *amqp.Error)))
	}()

	rabbit.logger.Infof("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		rabbit.logger.Errorf("Channel: %s", err)

	}

	rabbit.logger.Infof("got Channel, declaring Exchange (%q)", exchangeName)
	if err = channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		rabbit.logger.Errorf("Exchange Declare: %s", err)

	}

	rabbit.logger.Infof("declared Exchange, declaring Queue %q", queueName)
	queue, err := channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		amqp.Table{
			"x-dead-letter-exchange":    exchangeNameDLX,
			"x-dead-letter-routing-key": routingKeyDLX,
		}, // arguments (DLX settings)
	)
	if err != nil {
		rabbit.logger.Errorf("Queue Declare: %s", err)

	}

	rabbit.logger.Infof("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, routingKey)

	if err = channel.QueueBind(
		queue.Name,   // name of the queue
		routingKey,   // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		rabbit.logger.Errorf("Queue Bind: %s", err)

	}

	rabbit.logger.Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", "tag")
	deliveries, err := channel.Consume(
		queueName,      // name
		"tag",          // consumerTag,
		rabbit.autoAck, // autoAck
		false,          // exclusive
		false,          // noLocal
		false,          // noWait
		nil,            // arguments
	)
	if err != nil {
		rabbit.logger.Errorf("Queue Consume: %s", err)

	}

	go rabbit.handle(msgChan, deliveries, rabbit.errChan)
	return msgChan, func() {
		connection.Close()
	}
}
func (rabbit *appRabbitMQ) handle(channel chan *Message, deliveries <-chan amqp.Delivery, done chan error) {
	cleanup := func() {
		rabbit.logger.Infof("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()

	for d := range deliveries {
		rabbit.logger.Debugf("Headers ctx %v", d.Headers)
		ctx := common.ExtractAMQPHeaders(context.Background(), d.Headers)
		rabbit.logger.Debugf("ExtractAMQPHeaders ctx %v", ctx)
		// Create a new span
		tr := otel.Tracer("amqp")
		_, messageSpan := tr.Start(ctx, "AMQP - consume - checkout.processed")
		time.Sleep(1 * time.Millisecond)
		messageSpan.End()
		rabbit.deliveryCount++
		if rabbit.verbose {
			rabbit.logger.Infof(
				"got %dB delivery: [%v] %q",
				len(d.Body),
				d.DeliveryTag,
				d.Body,
			)
		} else {
			if rabbit.deliveryCount%65536 == 0 {
				rabbit.logger.Infof("delivery count %d", rabbit.deliveryCount)
			}
		}
		message := DecodeToMessage(d.Body)
		// var message *pubsub.Message

		// message.DecodeToMessage(d.Body)

		rabbit.logger.Infof("delivery message %d", message)
		channel <- message

		retryCount := d.Headers["x-death"]
		// d.MessageCount
		rabbit.logger.Infof("retryCount %v", retryCount)
		if !rabbit.autoAck {
			d.Ack(false)
		}
		// d.Nack()
	}
}

func (rabbit *appRabbitMQ) PublishRetry(ctx context.Context, data *Message) error {

	// key_task := "binding.task.*"
	queue_task := "queue.task"
	// type_task := "direct"

	ex_route_delay := "ex.route.delay"
	key_route_delay := "binding.#.1"
	type_route_delay := "topic"
	queue_route_delay := "queue.route.delay"

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-producer")
	rabbit.logger.Infof("dialing %q", rabbit.url)

	connection, err := amqp.DialConfig(rabbit.url, config)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}
	defer connection.Close()

	rabbit.logger.Infof("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}
	defer channel.Close()

	// if err := channel.ExchangeDeclare(
	// 	publishConfig.ExchangeName, // name
	// 	"direct",                   // type
	// 	true,                       // durable
	// 	false,                      // auto-deleted
	// 	false,                      // internal
	// 	false,                      // noWait
	// 	nil,                        // arguments
	// ); err != nil {
	// 	return fmt.Errorf("exchange Declare: %s", err)
	// }

	queueTask, err := channel.QueueDeclare(
		queue_task, // name of the queue
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // noWait
		amqp.Table{
			"x-dead-letter-exchange":    ex_route_delay,
			"x-dead-letter-routing-key": key_route_delay,
			"x-message-ttl":             10000,
		},
	)
	if err != nil {
		rabbit.logger.Errorf("producer: Queue Declare error: %s", err)
		return fmt.Errorf("producer: Queue Declare error: %s", err)
	}
	if err := channel.PublishWithContext(
		context.Background(),
		"",             // publish to an exchange
		queueTask.Name, // routing to 0 or more queues
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{

			ContentType:     "text/plain",
			ContentEncoding: "",
			// Expiration:      "10000", // TTL in milliseconds (5 seconds),
			Body:         data.EncodeToBytes(),
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:     0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %s", err)
	}
	rabbit.logger.Infof(" publishing messages success")
	// route delay

	if err := channel.ExchangeDeclare(
		ex_route_delay,   // name
		type_route_delay, // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return fmt.Errorf("failed exchange Declare DLX: %s", err)
	}
	rabbit.logger.Infof("declared Exchange DLX success %v", ex_route_delay)

	// Declare a Dead Letter Queue
	queueRouteDelay, err := channel.QueueDeclare(
		queue_route_delay, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a Dead Letter Queue: %s", err)
	}
	log.Printf("producer: declaring binding")
	if err := channel.QueueBind(queueRouteDelay.Name, key_route_delay, ex_route_delay, false, nil); err != nil {
		return fmt.Errorf("producer: Queue Bind: %s", err)
	}

	dlvRouteKey, err := channel.Consume(
		queueRouteDelay.Name, // name
		"",                   // consumerTag,
		false,                // autoAck
		false,                // exclusive
		false,                // noLocal
		false,                // noWait
		nil,                  // arguments
	)
	if err != nil {
		rabbit.logger.Errorf("Queue Consume: %s", err)
		return fmt.Errorf("Queue Consume: %s", err)
	}
	for d := range dlvRouteKey {
		log.Printf("d %v", d)
		message := DecodeToMessage(d.Body)
		message.RetryCount++
		d.Ack(false)
		err := rabbit.sendMessageToRabbitMQDelay(ctx, channel, message)
		if err != nil {
			rabbit.logger.Errorf("Send message to rabbit MQ Delay error: %s", err)
			return fmt.Errorf("send message to rabbit MQ Delay error: %s", err)
		}

	}

	return nil

	// return rabbit.sendMessageToRabbitMQDelay(ctx, data)
}

func (rabbit *appRabbitMQ) PublishRetryWithDLX(ctx context.Context) error {

	// QUEUE_microservice ---exp time-----> EXCHANGE DLX  ----- consumer update retry count ----->
	// Queue delay  ---- exp time -------> queue_microserivce

	//rabbit.logger.Infof("publishConfig %v", publishConfig)

	// queue_microservice:="queue_microservice"
	// exchange


	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-producer")
	rabbit.logger.Infof("dialing %q", rabbit.url)

	connection, err := amqp.DialConfig(rabbit.url, config)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}
	defer connection.Close()

	rabbit.logger.Infof("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}
	defer channel.Close()
	//rabbit.logger.Infof("got Channel, declaring %q Exchange (%q)", publishConfig.ExchangeType, publishConfig.ExchangeName)

	// if err := channel.ExchangeDeclare(
	// 	publishConfig.ExchangeName, // name
	// 	"direct",                   // type
	// 	true,                       // durable
	// 	false,                      // auto-deleted
	// 	false,                      // internal
	// 	false,                      // noWait
	// 	nil,                        // arguments
	// ); err != nil {
	// 	return fmt.Errorf("exchange Declare: %s", err)
	// }

	if len(publishConfig.QueueName) > 0 {
		rabbit.logger.Infof("producer: declaring queue '%s' '%s' '%s'",
			publishConfig.QueueName,
			publishConfig.ExchangeNameDLX,
			publishConfig.RoutingKeyDLX)
		queue, err := channel.QueueDeclare(
			publishConfig.QueueName, // name of the queue
			true,                    // durable
			false,                   // delete when unused
			false,                   // exclusive
			false,                   // noWait
			// nil,
			amqp.Table{
				"x-dead-letter-exchange":    publishConfig.ExchangeNameDLX,
				"x-dead-letter-routing-key": publishConfig.RoutingKeyDLX,
				"x-message-ttl":             10000,
			}, // arguments (DLX settings)
		)
		if err != nil {
			rabbit.logger.Errorf("producer: Queue Declare error: %s", err)
			return fmt.Errorf("producer: Queue Declare error: %s", err)
		} else {
			rabbit.logger.Infof("producer: declared queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
				queue.Name, queue.Messages, queue.Consumers, publishConfig.RoutingKey)
		}

		log.Printf("producer: declaring binding")
		if err := channel.QueueBind(queue.Name, publishConfig.RoutingKey, publishConfig.ExchangeName, false, nil); err != nil {
			return fmt.Errorf("producer: Queue Bind: %s", err)
		}
	}

	// Declare a Dead Letter Exchange
	if len(publishConfig.ExchangeNameDLX) > 0 {
		if err := channel.ExchangeDeclare(
			publishConfig.ExchangeNameDLX, // name
			"direct",                      // type
			true,                          // durable
			false,                         // auto-deleted
			false,                         // internal
			false,                         // noWait
			nil,                           // arguments
		); err != nil {
			return fmt.Errorf("failed exchange Declare DLX: %s", err)
		}
		rabbit.logger.Infof("declared Exchange DLX success %v", publishConfig.ExchangeNameDLX)
	}

	if len(publishConfig.QueueNameDLX) > 0 {

		// Declare a Dead Letter Queue
		queueDLX, err := channel.QueueDeclare(
			publishConfig.QueueNameDLX, // name
			true,                       // durable
			false,                      // delete when unused
			false,                      // exclusive
			false,                      // no-wait
			amqp.Table{
				"x-dead-letter-exchange":    publishConfig.ExchangeName,
				"x-dead-letter-routing-key": publishConfig.RoutingKey,
				"x-message-ttl":             10000,
			}, // arguments (DLX settings)
		)
		if err != nil {
			return fmt.Errorf("failed to declare a Dead Letter Queue: %s", err)
		}
		log.Printf("producer: declaring binding")
		if err := channel.QueueBind(queueDLX.Name, publishConfig.RoutingKeyDLX, publishConfig.ExchangeNameDLX, false, nil); err != nil {
			return fmt.Errorf("producer: Queue Bind: %s", err)
		}
	}

	rabbit.logger.Infof("declared Exchange, publishing messages")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := channel.PublishWithContext(ctx,
		publishConfig.ExchangeNameDLX, // publish to an exchange
		publishConfig.RoutingKeyDLX,   // routing to 0 or more queues
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{

			ContentType:     "text/plain",
			ContentEncoding: "",
			// Expiration:      "10000", // TTL in milliseconds (5 seconds),
			Body:         publishConfig.Data.EncodeToBytes(),
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:     0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %s", err)
	}
	rabbit.logger.Infof(" publishing messages success")
	return nil
}
func (rabbit *appRabbitMQ) sendMessageToRabbitMQDelay(ctx context.Context, channel *amqp.Channel, data *Message) error {

	queue_delay := "queue.delay.1"

	ex_dlx_delay := "ex.dlx.delay.1"
	key_dlx_delay := "key.dlx.delay.1"
	type_dlx_delay := "topic"

	queue_dlx_delay := "queue.dlx.delay.1"

	// SEND MESSAGE TO QUEUE MICROSERVICE
	queueMicroservice, err := channel.QueueDeclare(
		queue_delay, // name of the queue
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // noWait
		amqp.Table{
			"x-dead-letter-exchange":    ex_dlx_delay,
			"x-dead-letter-routing-key": key_dlx_delay,
			"x-message-ttl":             10000,
		},
	)
	if err != nil {
		rabbit.logger.Errorf("sendMessageToRabbitMQDelay: Queue Declare error: %s", err)
		return fmt.Errorf("sendMessageToRabbitMQDelay: Queue Declare error: %s", err)
	}

	if err := channel.PublishWithContext(
		context.Background(),
		"",                     // publish to an exchange
		queueMicroservice.Name, // routing to 0 or more queues
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{

			ContentType:     "text/plain",
			ContentEncoding: "",
			// Expiration:      "10000", // TTL in milliseconds (5 seconds),
			Body:         data.EncodeToBytes(),
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:     0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %s", err)
	}

	// dlx delay

	if err := channel.ExchangeDeclare(
		ex_dlx_delay,   // name
		type_dlx_delay, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("failed exchange Declare DLX: %s", err)
	}
	rabbit.logger.Infof("sendMessageToRabbitMQDelay: declared Exchange DLX success %v", type_dlx_delay)

	// Declare a Dead Letter Queue
	queueDLXDelay, err := channel.QueueDeclare(
		queue_dlx_delay, // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a Dead Letter Queue: %s", err)
	}
	log.Printf("sendMessageToRabbitMQDelay: declaring binding")
	if err := channel.QueueBind(queueDLXDelay.Name, key_dlx_delay, ex_dlx_delay, false, nil); err != nil {
		return fmt.Errorf("producer: Queue Bind: %s", err)
	}
	log.Printf("sendMessageToRabbitMQDelay: consumer")
	dlvDLXDelay, err := channel.Consume(
		queueDLXDelay.Name, // name
		"tag1",             // consumerTag,
		false,              // autoAck
		false,              // exclusive
		false,              // noLocal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		rabbit.logger.Errorf("Queue Consume: %s", err)

	}
	for d := range dlvDLXDelay {
		log.Printf("d %v", d)
		message := DecodeToMessage(d.Body)
		d.Ack(false)
		log.Printf("sendMessageToRabbitMQDelay: message %v", message.RetryCount)

	}

	return nil
}
