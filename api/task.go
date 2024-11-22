package api

type AcceptTaskResponse struct {
	Data  AcceptTaskData `json:"data"`
	Error ErrorMessage   `json:"error"`
}

type CompleteTaskResponse struct {
	Data  CompleteTaskData `json:"data"`
	Error ErrorMessage     `json:"error"`
}

type AcceptTaskData struct {
	Cooldown  Cooldown  `json:"cooldown"`
	Character Character `json:"character"`
	Task      Task      `json:"task"`
}

type CompleteTaskData struct {
	Cooldown  Cooldown    `json:"cooldown"`
	Character Character   `json:"character"`
	Rewards   TaskRewards `json:"rewards"`
}

type Task struct {
	Code    string      `json:"code"`
	Type    string      `json:"type"`
	Total   int         `json:"total"`
	Rewards TaskRewards `json:"rewards"`
}

type TaskRewards struct {
	Gold  int          `json:"gold"`
	Items []SimpleItem `json:"items"`
}
