package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func handlerWar(gs *gamelogic.GameState, logsCh *amqp.Channel) func(gamelogic.RecognitionOfWar) int {
	return func(rec gamelogic.RecognitionOfWar) int {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rec)

		if outcome == gamelogic.WarOutcomeNotInvolved {
			return pubsub.NackRequeue
		}

		if outcome == gamelogic.WarOutcomeNoUnits {
			return pubsub.NackDiscard
		}

		if outcome == gamelogic.WarOutcomeOpponentWon {
			msg := fmt.Sprintf("%v won a war against %v", winner, loser)
			gameLog := routing.GameLog{
				CurrentTime: time.Now(),
				Message:     msg,
				Username:    gs.GetUsername(),
			}

			err := pubsub.PublishGob(logsCh, routing.ExchangePerilTopic, routing.GameLogSlug+"."+gs.GetUsername(), gameLog)
			if err != nil {
				fmt.Println("error! could not publish gob")
				return pubsub.NackRequeue
			}

			return pubsub.Ack
		}

		if outcome == gamelogic.WarOutcomeYouWon {
			msg := fmt.Sprintf("%v won a war against %v", winner, loser)
			gameLog := routing.GameLog{
				CurrentTime: time.Now(),
				Message:     msg,
				Username:    gs.GetUsername(),
			}

			err := pubsub.PublishGob(logsCh, routing.ExchangePerilTopic, routing.GameLogSlug+"."+gs.GetUsername(), gameLog)
			if err != nil {
				fmt.Println("error! could not publish gob")
				return pubsub.NackRequeue
			}

			return pubsub.Ack
		}

		if outcome == gamelogic.WarOutcomeDraw {
			msg := fmt.Sprintf("A war between %v and %v resulted in a draw", winner, loser)
			gameLog := routing.GameLog{
				CurrentTime: time.Now(),
				Message:     msg,
				Username:    gs.GetUsername(),
			}

			err := pubsub.PublishGob(logsCh, routing.ExchangePerilTopic, routing.GameLogSlug+"."+gs.GetUsername(), gameLog)
			if err != nil {
				fmt.Println("error! could not publish gob")
				log.Fatal()
			}

			return pubsub.Ack
		}

		fmt.Printf("error. the game logic wasn't recognized: %v\n", rec)
		return pubsub.NackDiscard
	}
}
