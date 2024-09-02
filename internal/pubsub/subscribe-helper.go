package pubsub

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
)

func subscribe[T any](conn *amqp.Connection, exchange, queueName, key string, simpleQueueType int, handler func(T) int, unmarshaller func([]byte) (T, error)) (err error) {
	fmt.Println("somebody has subscribed to " + exchange)
	fmt.Println("queue name is " + queueName)
	fmt.Println("key is " + key)
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return
	}

	err = ch.Qos(10, 0, true)
	if err != nil {
		return
	}

	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return
	}

	var unmarshalled T
	go func() {
		fmt.Println("will consume deliveries ")
		for delivery := range deliveryCh {
			fmt.Println("a delivery came in!")
			unmarshalled, err = unmarshaller(delivery.Body)
			fmt.Println(unmarshalled)
			fmt.Println(err)
			if err != nil {
				return
			}

			fmt.Println(reflect.TypeOf(unmarshalled))
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
				return
			}
		}

		fmt.Println("no longer consuming deliveries")
	}()

	return
}
