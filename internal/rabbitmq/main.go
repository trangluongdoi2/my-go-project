package rabbitmq

import (
	"errors"
	"fmt"
	"go-backend-project/internal/config"
	"go-backend-project/utils"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var AmqpQueue *AmqpQueueService

const (
	DELAY           = 5
	RECONNECT_DELAY = 5
	MAX_WAIT_NUM    = 50
)

type AmqpQueueService struct {
	Connection  *amqp.Connection
	Channel     *amqp.Channel
	retry       int
	ctag        string
	isConnected bool
	notifyClose chan *amqp.Error
}

func NewAmqpQueueService(cfg *config.Config) (*AmqpQueueService, error) {
	if AmqpQueue != nil {
		return AmqpQueue, nil
	}

	r := AmqpQueueService{}
	result, err := r.Connect(cfg)
	if err != nil {
		return nil, err
	}

	AmqpQueue = result
	go AmqpQueue.Reconnect(cfg)

	return AmqpQueue, nil
}

func (s *AmqpQueueService) Connect(cfg *config.Config) (*AmqpQueueService, error) {
	schema := cfg.RabbitMQConfig.Schema
	username := cfg.RabbitMQConfig.Username
	password := cfg.RabbitMQConfig.Password
	host := cfg.RabbitMQConfig.Host
	port := cfg.RabbitMQConfig.Port
	vhost := cfg.RabbitMQConfig.VHost
	ssl := cfg.RabbitMQConfig.SSL
	ctag := cfg.RabbitMQConfig.CTag
	uri := cfg.RabbitMQConfig.Uri
	fmt.Println(schema, username, password, host, port, vhost, ssl, "Debug cfg")

	connectionName := "conn-" + ctag

	// if ssl == true {
	// 	uri = fmt.Sprintf("%s://%s:%s@%s/%s", schema, username, password, host, vhost)
	// } else {
	// 	uri = fmt.Sprintf("%s://%s:%s@%s:%s/%s", schema, username, password, host, port, vhost)
	// }

	connection, err := amqp.DialConfig(uri, amqp.Config{Properties: amqp.Table{"connection_name": connectionName}})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	notifyClose := make(chan *amqp.Error)
	result := &AmqpQueueService{
		Connection:  connection,
		Channel:     channel,
		retry:       0,
		ctag:        ctag,
		isConnected: true,
		notifyClose: notifyClose,
	}
	result.Connection.NotifyClose(notifyClose)

	return result, nil
}

func (s *AmqpQueueService) Reconnect(cfg *config.Config) {
	retryCount := 0
	for {
		if s.isConnected == false || s.Connection.IsClosed() {
			retryCount++
			fmt.Printf("Reconnect... #%d \n", retryCount)

			amqpQueue, errConn := s.Connect(cfg)
			if errConn != nil {
				fmt.Println(errConn)
			} else {
				s = amqpQueue
				AmqpQueue = amqpQueue
				retryCount = 0
			}
		}
		time.Sleep(RECONNECT_DELAY * time.Second)
		select {
		case <-s.notifyClose:
			s.Close()
		}
	}
}

func (s *AmqpQueueService) Close() {
	fmt.Println("Close..")
}

func (r *AmqpQueueService) waitConn() (bool, error) {
	var _flag = false
	var i = 1

	for {
		if i == MAX_WAIT_NUM {
			break
		}
		if !AmqpQueue.Connection.IsClosed() {
			_flag = true
			break
		}

		fmt.Printf("Wait... #%d \n", i)
		i++
		time.Sleep(5 * time.Second)
	}
	if AmqpQueue.Connection.IsClosed() {
		return _flag, errors.New("Not connected to the producer")
	}
	return _flag, nil
}

func (s *AmqpQueueService) Send(queueName string, message []byte) error {
	if ok, err := s.waitConn(); !ok {
		logrus.WithField("Err", err.Error())
		fmt.Println("Start Queue err:", err)
		return err
	}
	if AmqpQueue != nil && AmqpQueue.isConnected {
		s = AmqpQueue
	}

	cfg := config.GetConfig()
	ch, err := s.Connection.Channel()
	defer ch.Close()
	if err != nil {
		logrus.WithField("Err", err.Error()).Error("Failed to open a channel")
		fmt.Println("Open Queue err:", err)
		return err
	}

	args := make(amqp.Table)
	if _, err := ch.QueueDeclare(queueName, true, false, false, false, args); err != nil {
		fmt.Println("QueueDeclare err:", err)
		return err
	}

	msg := amqp.Publishing{
		MessageId:    uuid.New().String(),
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         message,
	}

	fmt.Println("[DEBUG]-msg")
	utils.DebugJson(msg)
	fmt.Println(cfg.RabbitMQConfig.Retry, "cfg.RabbitMQConfig.Retry")

	if err := ch.Publish("", queueName, false, false, msg); err != nil {
		if err != nil && s.retry <= cfg.RabbitMQConfig.Retry {
			s.retry++
			time.Sleep(DELAY * time.Second)
			s.Send(queueName, message)
		}
		fmt.Println("Queue err:", err)
		return err
	}

	return nil
}

func (s *AmqpQueueService) GetOne(queueName string) (*amqp.Delivery, error) {
	if ok, err := s.waitConn(); !ok {
		return nil, err
	}

	if AmqpQueue != nil && AmqpQueue.isConnected {
		s = AmqpQueue
	}

	ch, err := s.Connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot open channel: %w", err)
	}

	autoAck := false // Manual ack message

	msg, ok, err := ch.Get(queueName, autoAck)
	if err != nil {
		return nil, fmt.Errorf("cannot get message: %w", err)
	}

	if !ok {
		return nil, nil
	}

	return &msg, nil
}
