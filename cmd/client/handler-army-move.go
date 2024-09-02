package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerArmyMove(gs *gamelogic.GameState, armyMovesCh *amqp.Channel) func(gamelogic.ArmyMove) int {
	return func(armyMove gamelogic.ArmyMove) int {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(armyMove)
		fmt.Println("got a move outcome ... ...")
		fmt.Println(moveOutcome)
		if moveOutcome == gamelogic.MoveOutComeSafe {
			return pubsub.Ack
		}

		if moveOutcome == gamelogic.MoveOutcomeMakeWar {
			/*
			1. Publish a message to the "topic" exchange with the routing key $WARPREFIX.$USERNAME. 
				routing.WarRecognitionsPrefix contains the $WARPREFIX constant, 
				and the $USERNAME should be the name of the player consuming the move.
			2. NackRequeue the message... Might seem crazy, but it will be fun.
			*/
			fmt.Println("now to publish json with move outcome make war")
			fmt.Println(armyMovesCh)
			fmt.Println(routing.ExchangePerilTopic)
			routingKey := routing.WarRecognitionsPrefix + "." + gs.GetUsername()
			fmt.Println(routingKey)
			message := gamelogic.RecognitionOfWar{
				Attacker: gs.Player,
				Defender: armyMove.Player,
			}

			fmt.Println(message)
			err := pubsub.PublishJSON(armyMovesCh, routing.ExchangePerilTopic, routingKey, message)
			if err != nil {
				fmt.Println("got error!")
				fmt.Println(err)
				fmt.Println("publishing has failed! nack requeuing the message...")
				return pubsub.NackRequeue
			}

			return pubsub.Ack
		}


		return pubsub.NackDiscard
	}

}
