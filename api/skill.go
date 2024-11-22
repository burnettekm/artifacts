package api

type SkillResponse struct {
	Data  SkillData    `json:"data"`
	Error ErrorMessage `json:"error"`
}

type SkillData struct {
	Cooldown  Cooldown     `json:"cooldown"`
	Details   SkillDetails `json:"details"`
	Character Character    `json:"character"`
}
