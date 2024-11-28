package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MonsterResponse struct {
	Data  []MonsterData `json:"data"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
	Pages int           `json:"pages"`
}

type MonsterData struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Level       int    `json:"level"`
	Hp          int    `json:"hp"`
	AttackFire  int    `json:"attack_fire"`
	AttackEarth int    `json:"attack_earth"`
	AttackWater int    `json:"attack_water"`
	AttackAir   int    `json:"attack_air"`
	ResFire     int    `json:"res_fire"`
	ResEarth    int    `json:"res_earth"`
	ResWater    int    `json:"res_water"`
	ResAir      int    `json:"res_air"`
	MinGold     int    `json:"min_gold"`
	MaxGold     int    `json:"max_gold"`
	Drops       []Drop `json:"drops"`
}

func (c *ArtifactsClient) GetMonsters(pageNumber int) ([]MonsterData, error) {
	path := "/monsters"
	params := map[string]string{
		"page": strconv.Itoa(pageNumber),
		"size": strconv.Itoa(100),
	}

	resp, err := c.Do(http.MethodGet, path, params, nil)
	if err != nil {
		return nil, fmt.Errorf("executing list monsters request: %w", err)
	}

	monsterResp := MonsterResponse{}
	if err := json.Unmarshal(resp, &monsterResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp: %w", err)
	}
	return monsterResp.Data, nil
}
