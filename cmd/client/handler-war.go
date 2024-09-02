package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"fmt"
)

func handlerWar(gs *gamelogic.GameState) func (gamelogic.RecognitionOfWar) int {
	return func(rec gamelogic.RecognitionOfWar) int {
		defer fmt.Print("> ")
		outcome, _, _ := gs.HandleWar(rec)
		if outcome == gamelogic.WarOutcomeNotInvolved {
			return pubsub.NackRequeue
		}

		if outcome == gamelogic.WarOutcomeNoUnits {
			return pubsub.NackDiscard
		}

		if outcome == gamelogic.WarOutcomeOpponentWon {
			return pubsub.Ack
		}

		if outcome == gamelogic.WarOutcomeYouWon {
			return pubsub.Ack
		}

		if outcome == gamelogic.WarOutcomeDraw {
			return pubsub.Ack
		}

		fmt.Printf("error. the game logic wasn't recognized: %v\n", rec)
		return pubsub.NackDiscard
	}
}

