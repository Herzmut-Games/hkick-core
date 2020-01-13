package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	scoreRed   = 0
	scoreWhite = 0

	gameStartTime  = time.Now()
	roundStartTime = time.Now()

	gameIsRunning = false

	// The teams will be fixed in the first round as follows:
	// team A will start as red
	// team B will start as white
	// after one round teams will be swapped. hkick-core needs to take care of that when counting win.
	teamAWinCount = 0
	teamBWinCount = 0

	winHistory  = []Round{}
	goalHistory = []Goal{}
	lastGoalHistory = []Goal{}

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
	if (!gameIsRunning) { return }
	goalHistory = append(goalHistory, Goal{Team: team, Time: time.Since(roundStartTime).Seconds()})
	playSound("goal")

	updateScore()
}

func undoScore() {
	if (!gameIsRunning) { return }

	if (len(goalHistory) == 0 && len(winHistory) >= 1) {
		goalHistory = lastGoalHistory
		winHistory = winHistory[:len(winHistory)-1]
	}

	if len(goalHistory) > 0 {
		goalHistory = goalHistory[:len(goalHistory)-1]
	}

	updateScore()
}

func resetScore() {
	if (!gameIsRunning) { return }
	goalHistory = []Goal{}
	roundStartTime = time.Now()
	updateScore()
}

func updateScore() {
  debug()

	scoreRed = 0
	scoreWhite = 0
	for _, goal := range goalHistory {
		switch goal.Team {
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
	publish("round/current", strconv.Itoa(currentRound()), true)

	goals, _ := json.Marshal(goalHistory)
	publish("round/goals", string(goals), true)

	if distance >= 2 {
		if (scoreRed >= 5) || (scoreWhite >= 5) {
			roundEnd()
		}
	} else if (scoreRed >= 8) || (scoreWhite >= 8) {
		roundEnd()
	}
}

func startGame() {
	clearAll()
	gameIsRunning = true
	publish("sound/play", "start", false)
	updateScore()

}

func currentRound() int {
	return len(winHistory) + 1
}

func teamsAreSwapped() bool {
	return currentRound() == 2
}

func nextRound() {
	publish("round/end", "end", false)
	publish("round/current", strconv.Itoa(currentRound()), true)
	
	rounds, _ := json.Marshal(winHistory)
	fmt.Printf(string(rounds))
	lastGoalHistory = goalHistory
	resetScore()
}

func roundEnd() {
	if scoreRed >= scoreWhite {
		if teamsAreSwapped() {
			winHistory = append(winHistory, Round{Winner: "b", Time: time.Since(roundStartTime).Seconds()})
		} else {
			winHistory = append(winHistory, Round{Winner: "a", Time: time.Since(roundStartTime).Seconds()})
		}
	} else {
		if teamsAreSwapped() {
			winHistory = append(winHistory, Round{Winner: "a", Time: time.Since(roundStartTime).Seconds()})
		} else {
			winHistory = append(winHistory, Round{Winner: "b", Time: time.Since(roundStartTime).Seconds()})
		}
	}

	teamAWinCount = 0
	teamBWinCount = 0
	for _, round := range winHistory {
		switch round.Winner {
		case "a":
				teamAWinCount++

		case "b":
				teamBWinCount++
		}
	}

	if teamAWinCount == 2 {
		gameEnd("a")
	} else if teamBWinCount == 2 {
		gameEnd("b")
	} else {
		nextRound()
	}
}

func gameEnd(winner string) {
	fmt.Println("game is over")
	rounds, _ := json.Marshal(winHistory)
	fmt.Printf(string(rounds))

	fmt.Printf("%s is the winner \n", winner)

	game, _ := json.Marshal(Game{Winner: winner, Time: time.Since(gameStartTime).Seconds()})
	publish("game/end", string(game), false)
	publish("sound/play", "end", false)

	resetScore()
	clearAll()
}

func clearAll() {
	gameStartTime = time.Now()
	roundStartTime = time.Now()
	teamAWinCount = 0
	teamBWinCount = 0
	scoreRed = 0
	scoreWhite = 0
	winHistory = []Round{}
	goalHistory = []Goal{}
	gameIsRunning = false

	publish("round/current", strconv.Itoa(currentRound()), true)
	
	debug()
}



func debug() {
	publish("debug/teamAWinCount", strconv.Itoa(teamAWinCount), false)
	publish("debug/teamBWinCount", strconv.Itoa(teamBWinCount), false)
}