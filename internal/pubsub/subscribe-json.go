package pubsub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

const Ack = 1
const NackRequeue = 2
const NackDiscard = 3

func jsonUnmarshaller[T any](data []byte) (v T, err error) {
	body := strings.Trim(string(data), "\"")
	decoded, err := base64.StdEncoding.WithPadding(base64.StdPadding).DecodeString(body)
	if err != nil {
		fmt.Println(err)
		fmt.Println("(message said " + string(data) + ")")

		err = json.Unmarshal([]byte(body), &v)
		if err != nil {
			fmt.Println(err)
			fmt.Println("(message said " + string(data) + ")")
		}

		return
	}

	fmt.Println("decoded message said " + string(decoded))
	err = json.Unmarshal(decoded, &v)
	return
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int,
	handler func(T) int,
) (err error) {

	err = subscribe(conn, exchange, queueName, key, simpleQueueType, handler, jsonUnmarshaller)
	return
}
