package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"bytes"
	"encoding/gob"
)

func gobUnmarshaller[T any](data []byte) (v T, err error) {
	buf := bytes.NewBuffer(data)
	err = gob.NewDecoder(buf).Decode(&v)
	return
}

	

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, simpleQueueType int, handler func(T) int) (err error) {
	err = subscribe(conn, exchange, queueName, key, simpleQueueType, handler, gobUnmarshaller)
	return
}
