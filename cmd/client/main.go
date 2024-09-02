package main

import (
	"time"
	"fmt"
	"strconv"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
)

func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("could not connect!")
	}

	defer conn.Close()

	fmt.Println("connection successful to " + connectionString)

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("no username, help!")
	}

	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, pubsub.Transient)

	fmt.Println("Starting Peril client...")
	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, "pause."+username, routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if err != nil {
		fmt.Println(err)
		log.Fatal("subscribe json error!")
	}

	armyMovesCh, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, "army_moves."+username, "army_moves.*", pubsub.Transient)
	if err != nil {
		fmt.Println(err)
		log.Fatal("declare and bind error!!")
	}

	pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, "army_moves."+username, "army_moves.*", pubsub.Transient, handlerArmyMove(gameState, armyMovesCh))
	logsCh, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+"."+username, pubsub.Durable)
	if err != nil {
		fmt.Println(err)
		log.Fatal("declare and bind error!!")
	}

	pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, "war", "war.*", pubsub.Durable, handlerWar(gameState, logsCh))
	for {
		words := gamelogic.GetInput()
		if len(words) > 0 {
			if words[0] == "spawn" {
				err := gameState.CommandSpawn(words)
				if err != nil {
					fmt.Println("Spawn Error!")
					fmt.Println(err)
				}
			} else if words[0] == "move" {
				armyMove, err := gameState.CommandMove(words)
				if err != nil {
					fmt.Println("Bad Move!")
					fmt.Println(err)
				} else {
					fmt.Println("Good move! 😊")
					fmt.Println(armyMove)

					pubsub.PublishJSON(armyMovesCh, routing.ExchangePerilTopic, "army_moves."+username, armyMove)
					fmt.Println("army move was published successfully:")
				}
			} else if words[0] == "status" {
				gameState.CommandStatus()
			} else if words[0] == "help" {
				gamelogic.PrintClientHelp()
			} else if words[0] == "spam" {
				if len(words) < 2 {
					fmt.Println("you have to say how many spams!")
				} else {
					n, err := strconv.Atoi(words[1])
					if err != nil {
						fmt.Println("you have to use an integer for the number of spams!")
					} else {
						for range(n) {
							maliciousLog := gamelogic.GetMaliciousLog()
							msg := routing.GameLog{
								CurrentTime: time.Now(),
								Message: maliciousLog,
								Username: gameState.GetUsername(),
							}

							pubsub.PublishGob(logsCh, routing.ExchangePerilTopic, routing.GameLogSlug + "." + gameState.GetUsername(), msg)
						}
					}
				}
			} else if words[0] == "quit" {
				fmt.Println("quitting...")
				break
			} else {
				fmt.Println("unrecognized input. doing nothing...")
			}

		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("the program is shutting down")
}
