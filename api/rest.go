package api

type RestResponse struct {
	Rest  Rest         `json:"data"`
	Error ErrorMessage `json:"error"`
}

type Rest struct {
	Cooldown   Cooldown  `json:"cooldown"`
	HpRestored int       `json:"hp_restored"`
	Character  Character `json:"character"`
}
