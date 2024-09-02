package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) int {
	return func(playingState routing.PlayingState) int {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
		return pubsub.Ack
	}

}
