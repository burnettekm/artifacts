package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ActionBankResponse struct {
	Data ActionBankData `json:"data"`
}

type GetBankResponse struct {
	Data []SimpleItem `json:"data"`
}

type ActionBankData struct {
	Cooldown  Cooldown      `json:"cooldown"`
	Item      CraftableItem `json:"item"`
	Bank      []SimpleItem  `json:"bank"`
	Character Character     `json:"character"`
}

func (c *Svc) WithdrawFromBankIfFound(itemCode string, quantity int) (int, error) {
	fmt.Printf("Checking bank for item %s", itemCode)
	path := "/my/bank/items"
	params := map[string]string{
		"item_code": itemCode,
	}

	respBytes, err := c.Client.Do(http.MethodGet, path, params, nil)
	if err != nil {
		return 0, fmt.Errorf("executing deposit bank request: %w", err)
	}
	bankResp := GetBankResponse{}
	if err := json.Unmarshal(respBytes, &bankResp); err != nil {
		return 0, fmt.Errorf("unmarshalling bank response: %w", err)
	}

	// withdraw the lesser of requested or available quantity
	// requested quantity 3 bank quantity 7, withdraw 3
	// requested 3 bank quantity 2, withdraw 2
	// check for 0 found in bank below
	var minQuantity int
	if quantity <= bankResp.Data[0].Quantity {
		minQuantity = quantity
	} else {
		minQuantity = bankResp.Data[0].Quantity
	}

	if minQuantity == 0 {
		return 0, nil
	}

	if err := c.WithdrawBankItem(bankResp.Data[0].Code, minQuantity); err != nil {
		return 0, fmt.Errorf("withdrawing %s from bank: %w", bankResp.Data[0].Code, err)
	}

	return minQuantity, nil
}

func (c *Svc) WithdrawBankItem(itemCode string, quantity int) error {
	fmt.Printf("Withdrawing %d %s\n", quantity, itemCode)
	path := fmt.Sprintf("/my/%s/action/bank/withdraw", c.Character.Name)

	body := SimpleItem{
		Code:     itemCode,
		Quantity: quantity,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling body: %w", err)
	}
	resp, err := c.Client.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return fmt.Errorf("executing withdraw request: %w", err)
	}

	withdrawResp := ActionBankResponse{}
	if err := json.Unmarshal(resp, &withdrawResp); err != nil {
		return fmt.Errorf("unmarshalling action bank response: %w", err)
	}
	c.Character = &withdrawResp.Data.Character
	c.Character.WaitForCooldown()

	fmt.Println("Withdraw complete")

	return nil
}

func (c *Svc) DepositBank(inventoryItem InventorySlot) error {
	fmt.Printf("Depositing item %s in the bank\n", inventoryItem.Code)

	contentType := "bank"
	mapResp, err := c.Client.GetMaps(nil, &contentType)
	if err != nil {
		return fmt.Errorf("getting bank location: %w", err)
	}

	if _, err := c.MoveCharacter(mapResp.Data[0].X, mapResp.Data[0].Y); err != nil {
		return fmt.Errorf("moving to bank: %w", err)
	}

	path := fmt.Sprintf("/my/%s/action/bank/deposit", c.Character.Name)
	body := SimpleItem{
		Code:     inventoryItem.Code,
		Quantity: inventoryItem.Quantity,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling body: %w", err)
	}
	respBytes, err := c.Client.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return fmt.Errorf("executing deposit bank request: %w", err)
	}
	bankResp := ActionBankResponse{}
	if err := json.Unmarshal(respBytes, &bankResp); err != nil {
		return fmt.Errorf("unmarshalling bank response: %w", err)
	}
	fmt.Println("Deposit complete")
	c.Character = &bankResp.Data.Character
	c.Character.WaitForCooldown()

	return nil
}
