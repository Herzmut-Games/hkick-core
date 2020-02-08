package main

import (
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func reportGame(winHistory []Round) {
	url := "https://dashboard.kickr.me"
	client := resty.New()
	var match BackendMatch

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&match).
		SetBody(BackendMatchWrap{}).
		Post(url + "/matches")

	fmt.Println("Response Info:")
	fmt.Println("Error      :", err)
	fmt.Println("Status     :", resp.Status())
	fmt.Println("Time       :", resp.Time())

	for _, round := range winHistory {

		backendRound := BackendRound{Duration: strconv.FormatFloat(round.Time, 'f', 6, 64), Winner: round.Winner, ScoreA: round.ScoreA, ScoreB: round.ScoreB, Match: match.ID}
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(BackendRoundWrap{Round: backendRound}).
			Post(url + "/rounds")

		fmt.Println("Response Info:")
		fmt.Println("Error      :", err)
		fmt.Println("Status     :", resp.Status())
		fmt.Println("Time       :", resp.Time())
		fmt.Println("Body       :", resp)
	}
}
