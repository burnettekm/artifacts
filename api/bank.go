package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (c *Svc) WithdrawFromBankIfFound(characterName, itemCode string, quantity int) (int, error) {
	fmt.Printf("%s searching bank for %d %s\n", characterName, quantity, itemCode)
	c.BankMutex.Lock()
	defer c.BankMutex.Unlock()

	item, found := c.GetBankItemsByCode(itemCode)
	if !found {
		return 0, nil
	}
	foundQuantity := item.Quantity

	// withdraw the lesser of requested or available quantity
	// requested quantity 3 bank quantity 7, withdraw 3
	// requested 3 bank quantity 2, withdraw 2
	// check for 0 found in bank below
	var minQuantity int
	if quantity <= foundQuantity {
		minQuantity = quantity
	} else {
		minQuantity = foundQuantity
	}

	if minQuantity > c.Characters[characterName].InventoryMaxItems {
		minQuantity = c.Characters[characterName].InventoryMaxItems
	}

	if err := c.WithdrawBankItem(characterName, itemCode, minQuantity); err != nil {
		return 0, fmt.Errorf("withdrawing %s from bank: %w", itemCode, err)
	}

	return minQuantity, nil
}

func (c *Svc) GetBankItems() ([]SimpleItem, error) {
	fmt.Println("Getting bank items")
	path := "/my/bank/items"
	params := map[string]string{
		"size": strconv.Itoa(100),
	}

	respBytes, err := c.Client.Do(http.MethodGet, path, params, nil)
	if err != nil {
		return nil, fmt.Errorf("executing deposit bank request: %w", err)
	}
	bankResp := GetBankResponse{}
	if err := json.Unmarshal(respBytes, &bankResp); err != nil {
		return nil, fmt.Errorf("unmarshalling bank response: %w", err)
	}

	return bankResp.Data, nil
}

func (c *Svc) WithdrawBankItem(characterName, itemCode string, quantity int) error {
	fmt.Printf("%s withdrawing %d %s\n", characterName, quantity, itemCode)

	// find location of bank
	coords := c.GetCoordinatesByCode("bank")
	if _, err := c.MoveCharacter(characterName, coords[0].X, coords[0].Y); err != nil {
		return fmt.Errorf("moving to bank: %w", err)
	}

	path := fmt.Sprintf("/my/%s/action/bank/withdraw", characterName)

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

	c.updateBank(withdrawResp.Data.Bank)

	c.Characters[characterName] = &withdrawResp.Data.Character
	c.Characters[characterName].WaitForCooldown()

	fmt.Println("Withdraw complete")

	return nil
}

func (c *Svc) DepositBank(characterName string, inventoryItem InventorySlot) error {
	c.BankMutex.Lock()
	defer c.BankMutex.Unlock()

	fmt.Printf("%s depositing item %s in the bank\n", characterName, inventoryItem.Code)
	coords := c.GetCoordinatesByCode("bank")
	if _, err := c.MoveCharacter(characterName, coords[0].X, coords[0].Y); err != nil {
		return fmt.Errorf("moving to bank: %w", err)
	}

	path := fmt.Sprintf("/my/%s/action/bank/deposit", characterName)
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
	c.updateBank(bankResp.Data.Bank)
	fmt.Println("Deposit complete")

	c.Characters[characterName] = &bankResp.Data.Character
	c.Characters[characterName].WaitForCooldown()

	return nil
}

func (c *Svc) DepositAllItems(characterName string) error {
	for _, item := range c.Characters[characterName].Inventory {
		if item.Code == "" {
			continue
		}
		if err := c.DepositBank(characterName, item); err != nil {
			return fmt.Errorf("depositing item %s: %w", item.Code, err)
		}
	}

	return nil
}
