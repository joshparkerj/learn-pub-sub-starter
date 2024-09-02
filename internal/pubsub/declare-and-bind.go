package pubsub

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

const Durable = 1
const Transient = 2

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (ch *amqp.Channel, q amqp.Queue, err error) {
	ch, err = conn.Channel()
	if err != nil {
		return
	}

	var durable, autoDelete, exclusive bool
	noWait := false
	if simpleQueueType == Durable {
		durable = true
		autoDelete = false
		exclusive = false
	} else if simpleQueueType == Transient {
		durable = false
		autoDelete = true
		exclusive = true
	} else {
		err = errors.New("use pubsub.Durable or pubsub.Transient for the simpleQueueType please!")
		return
	}

	table := make(map[string]any)
	table["x-dead-letter-exchange"] = "peril_dlx"
	q, err = ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, table)
	if err != nil {
		return
	}

	err = ch.QueueBind(queueName, key, exchange, noWait, nil)
	return
}
