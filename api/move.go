package api

type MoveRequestBody struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MoveResponse struct {
	Data  MoveData     `json:"data"`
	Error ErrorMessage `json:"error"`
}

type MoveData struct {
	Cooldown    Cooldown    `json:"cooldown"`
	Destination Destination `json:"destination"`
	Character   Character   `json:"character"`
}

type Destination struct {
	Name    string  `json:"name"`
	Skin    string  `json:"skin"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
	Content Content `json:"content"`
}
