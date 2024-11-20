package api

type FightResponse struct {
	Data  FightData    `json:"data"`
	Error ErrorMessage `json:"error"`
}

type FightData struct {
	Cooldown  Cooldown  `json:"cooldown"`
	Fight     Fight     `json:"fight"`
	Character Character `json:"character"`
}

type Fight struct {
	Xp                 int         `json:"xp"`
	Gold               int         `json:"gold"`
	Drops              []Item      `json:"drops"`
	Turns              int         `json:"turns"`
	MonsterBlockedHits BlockedHits `json:"monster_blocked_hits"`
	PlayerBlockedHits  BlockedHits `json:"player_blocked_hits"`
	Logs               []string    `json:"logs"`
	Result             string      `json:"result"`
}

type BlockedHits struct {
	Fire  int `json:"fire"`
	Earth int `json:"earth"`
	Water int `json:"water"`
	Air   int `json:"air"`
	Total int `json:"total"`
}
