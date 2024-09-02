package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) (err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(val)
	if err != nil {
		return
	}

	ctx := context.Background()
	mandatory := false
	immediate := false
	msg := amqp.Publishing{
		ContentType: "application/gob",
		Body:        buf.Bytes(),
	}

	err = ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)

	return
}
