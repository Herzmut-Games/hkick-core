package main

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"os/signal"
	"strconv"
)

var (
	scoreRed            = 0
	scoreWhite          = 0
	currentRound        = 0
	goalHistory         = []string{}
	availableSoundModes = []string{"default", "meme", "quake", "techno"}
	currentSoundMode    = "random"
)

func mqttURI() *url.URL {
	uri, err := url.Parse("mqtt://172.30.1.32:1883")
	if err != nil {
		log.Fatal(err)
	}

	return uri
}

func playSound(event string) {
	if currentSoundMode == "random" {
		publish("sound/play", event, false)

	} else {
		publish("sound/play", event+"/"+currentSoundMode, false)

	}
}

func main() {
	connect("hkick-core", mqttURI())
	go subscribe(mqttURI())

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)
	select {
	case <-quit:
	}
}

func leadingTeam() string {
	if scoreRed > scoreWhite {
		return "red"
	}

	return "white"
}

func increaseScore(team string) {
	goalHistory = append(goalHistory, team)
	playSound("goal")

	updateScore()
}

func undoScore() {
	if len(goalHistory) > 0 {
		goalHistory = goalHistory[:len(goalHistory)-1]
	}

	updateScore()
}

func resetScore() {
	goalHistory = []string{}
	updateScore()
}

func updateScore() {
	scoreRed = 0
	scoreWhite = 0
	for _, team := range goalHistory {
		switch team {
		case "red":
			scoreRed++
		case "white":
			scoreWhite++
		}
	}

	distance := int(math.Abs(float64(scoreRed - scoreWhite)))

	fmt.Printf("red is %d and white is %d (distance %d)\n", scoreRed, scoreWhite, distance)

	publish("score/red", strconv.Itoa(scoreRed), true)
	publish("score/white", strconv.Itoa(scoreWhite), true)

	if distance >= 2 {
		if (scoreRed >= 5) || (scoreWhite >= 5) {
			roundEnd()
		}
	} else if (scoreRed >= 8) || (scoreWhite >= 8) {
		roundEnd()
	}
}

func startGame() {
	currentRound = 1
	publish("game/round", strconv.Itoa(currentRound), true)
	publish("sound/play", "start", false)

}

func nextRound() {
	currentRound++
	resetScore()
}

func roundEnd() {
	if currentRound < 3 {
		nextRound()
	} else {
		gameEnd()
	}
}

func gameEnd() {
	fmt.Println("game is over")

	winner := leadingTeam()
	fmt.Printf("%s is the winner \n", winner)

	publish("game/end", winner, false)

	resetScore()
}
