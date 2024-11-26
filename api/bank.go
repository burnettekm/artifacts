package api

type BankResponse struct {
	Data  BankData     `json:"data"`
	Error ErrorMessage `json:"error"`
}

type BankData struct {
	Cooldown  Cooldown      `json:"cooldown"`
	Item      CraftableItem `json:"item"`
	Bank      []SimpleItem  `json:"bank"`
	Character Character     `json:"character"`
}

func (c *Svc) DepositBank(characterName string, itemDetail SimpleItem) error {

}
