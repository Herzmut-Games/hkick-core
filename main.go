package main

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"os/signal"
)

var (
	scoreRed   = 0
	scoreWhite = 0
)

func mqttURI() *url.URL {
	uri, err := url.Parse("mqtt://172.30.1.32:1883")
	if err != nil {
		log.Fatal(err)
	}

	return uri
}

func main() {

	go listen(mqttURI(), "goals")

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)
	select {
	case <-quit:
	}
}

func handleGoal(team string) {
	fmt.Println("goal")

	if team == "white" {
		scoreWhite++
	} else if team == "red" {
		scoreRed++
	}

	distance := int(math.Abs(float64(scoreRed - scoreWhite)))

	fmt.Printf("distance is %d\n", distance)
	fmt.Printf("scoreRed is %d and scoreWhite is %d\n", scoreRed, scoreWhite)

	if distance >= 2 {
		if (scoreRed >= 5) || (scoreWhite >= 5) {
			gameEnd()
		}
	} else if (scoreRed >= 8) || (scoreWhite >= 8) {
		gameEnd()
	}

}

func leadingTeam() string {
	if scoreRed > scoreWhite {
		return "red"
	}

	return "white"

}

func gameEnd() {
	fmt.Println("game is over")

	client := connect("pub", mqttURI())

	winner := leadingTeam()
	fmt.Printf("%s is the winner \n", winner)
	client.Publish("game", 0, false, winner)

	scoreWhite = 0
	scoreRed = 0
}
