package pubsub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
	"strings"
)

const Ack = 1
const NackRequeue = 2
const NackDiscard = 3

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int,
	handler func(T) int,
) (err error) {
	fmt.Println("somebody has subscribed json to " + exchange)
	fmt.Println("queue name is " + queueName)
	fmt.Println("key is " + key)
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return
	}

	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return
	}

	go func() {
		fmt.Println("will consume deliveries ")
		for delivery := range deliveryCh {
			var unmarshalled T
			fmt.Println(reflect.TypeOf(unmarshalled))
			body := strings.Trim(string(delivery.Body), "\"")
			decoded, err := base64.StdEncoding.WithPadding(base64.StdPadding).DecodeString(body)
			if err != nil {
				fmt.Println(err)
				fmt.Println("(message said " + string(delivery.Body) + ")")

				err = json.Unmarshal([]byte(body), &unmarshalled)
				if err != nil {
					fmt.Println(err)
					fmt.Println("(message said " + string(delivery.Body) + ")")
					continue
				}

				fmt.Println(unmarshalled)
				acktype := handler(unmarshalled)
				if acktype == Ack {
					fmt.Println("ack")
					err = delivery.Ack(false)
				} else if acktype == NackRequeue {
					fmt.Println("nack requeue")
					err = delivery.Nack(false, true)
				} else if acktype == NackDiscard {
					fmt.Println("nack discard")
					err = delivery.Nack(false, false)
				}
				if err != nil {
					fmt.Println(err)
					fmt.Println("(message said " + string(delivery.Body) + ")")
				}
				continue

			}

			fmt.Println("decoded message said " + string(decoded))
			err = json.Unmarshal(decoded, &unmarshalled)
			if err != nil {
				fmt.Println(err)
				fmt.Println("(message said " + string(delivery.Body) + ")")
				continue
			}

			fmt.Println(unmarshalled)
			handler(unmarshalled)
			err = delivery.Ack(false)
			if err != nil {
				fmt.Println(err)
				fmt.Println("(message said " + string(delivery.Body) + ")")
				continue
			}
		}

		fmt.Println("no longer consuming deliveries")

		return
	}()

	return
}
