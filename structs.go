package main

type Goal struct {
	Team string
	Time float64
}

type Round struct {
	Winner string
	ScoreA int
	ScoreB int
	Time   float64
}

type Game struct {
	Winner string
	Time   float64
}

type BackendMatch struct {
	ID int `json:"id"`
}

type BackendMatchWrap struct {
	Match BackendMatch `json:"match"`
}

type BackendRound struct {
	Winner   string `json:"winner"`
	ScoreA   int    `json:"team_a_score"`
	ScoreB   int    `json:"team_b_score"`
	Duration string `json:"duration"`
	Match    int    `json:"match_id"`
}

type BackendRoundWrap struct {
	Round BackendRound `json:"round"`
}
