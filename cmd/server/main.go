package main

import (
	"encoding/json"
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
)

func sendMessage(connChan *amqp.Channel, shouldPause bool) {
	msgStruct := routing.PlayingState{
		IsPaused: shouldPause,
	}

	msg, err := json.Marshal(msgStruct)
	if err != nil {
		log.Fatal("could not serialize message!")
	}

	fmt.Println("the key is " + routing.PauseKey)
	err = pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, msg)
	if err != nil {
		log.Fatal("could not publish json message!")
	}
}

func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("could not connect!")
	}

	defer conn.Close()

	fmt.Println("connection successful to " + connectionString)

	connChan, err := conn.Channel()
	if err != nil {
		log.Fatal("connection channel error!")
	}

	fmt.Println("Starting Peril server...")

	pubsub.SubscribeGob(conn, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", pubsub.Durable, func(gl routing.GameLog) int {
		defer fmt.Print("> ")
		fmt.Println("got gob I think")
		err = gamelogic.WriteLog(gl)
		if err != nil {
			fmt.Println("error reading logs! yikes...")
			return pubsub.NackDiscard
		}

		return pubsub.Ack
	})

	if err != nil {
		log.Fatal("could not declare and bind the durable queue!")
	}

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if len(words) > 0 {
			if words[0] == "pause" {
				fmt.Println("sending a pause message")
				sendMessage(connChan, true)
			} else if words[0] == "resume" {
				fmt.Println("Sending a resume message")
				sendMessage(connChan, false)
			} else if words[0] == "quit" {
				fmt.Println("ending loop")
				break
			} else {
				fmt.Println("command not recognized... doing nothing")
			}
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("the program is shutting down")
}
