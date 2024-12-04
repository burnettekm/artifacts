package api

import (
	"fmt"
)

func (c *Svc) CraftItem(characterName, code string, quantity int) (*CraftableItem, error) {
	fmt.Printf("%s attempting to craft item %s, quantity: %d\n", characterName, code, quantity)

	item := c.GetItem(code)
	if item.Craft == nil {
		return nil, nil
	}

	// verify character can craft item
	if !c.GetCharacterByName(characterName).AbleToCraft(item.Craft.Skill, item.Craft.Level) {
		return nil, fmt.Errorf("unable to craft item: required level: %d", item.Craft.Level)
	}

	// get dependent items
	for _, subItem := range item.Craft.Items {
		remainingQuantity := subItem.Quantity * quantity
		fmt.Printf("%s needs %d %s to craft %s\n", characterName, remainingQuantity, subItem.Code, code)
		craftable := c.GetItem(subItem.Code)
		// check if item equipped
		if c.GetCharacterByName(characterName).IsEquipped(craftable) {
			if err := c.Unequip(characterName, craftable); err != nil {
				return nil, fmt.Errorf("unequipping item for crafting: %w", err)
			}
			continue
		}

		// check if item in inventory
		found, inventoryQuantity := c.GetCharacterByName(characterName).FindItemInInventory(subItem.Code)
		if found && inventoryQuantity >= remainingQuantity {
			continue
		}
		remainingQuantity -= inventoryQuantity

		// check for item in bank
		bankQuantity, err := c.WithdrawFromBankIfFound(characterName, subItem.Code, remainingQuantity)
		if err != nil {
			return nil, fmt.Errorf("withdrawing %s from bank if found: %w", subItem.Code, err)
		}
		c.Characters[characterName].WaitForCooldown()
		remainingQuantity -= bankQuantity

		if remainingQuantity <= 0 {
			continue
		}

		if craftable.Craft != nil && !c.Characters[characterName].AbleToCraft(craftable.Craft.Skill, craftable.Craft.Level) {
			return nil, fmt.Errorf("unable to craft subitem: %s: needs %s level: %d", craftable.Name, craftable.Craft.Skill, craftable.Craft.Level)
		}

		fmt.Printf("%s gathering subitem: %s\n", characterName, craftable.Code)
		switch craftable.Subtype {
		case "mob":
			fightQty := subItem.Quantity * quantity
			if err := c.FightForCrafting(characterName, craftable.Code, &fightQty); err != nil {
				return nil, fmt.Errorf("%s fighting for required item %s: %w", characterName, craftable.Code, err)
			}
		default:
			if craftable.Craft == nil {
				fmt.Printf("%s needs to gather to craft %d %s\n", characterName, remainingQuantity, craftable.Code)
				if err := c.Gather(characterName, craftable, remainingQuantity); err != nil {
					return nil, fmt.Errorf("%s gathering required item: %s: %w", characterName, craftable.Code, err)
				}
			} else {
				if _, err := c.CraftItem(characterName, craftable.Code, remainingQuantity); err != nil {
					return nil, fmt.Errorf("%s crafting subitem %s: %w", characterName, craftable.Code, err)
				}
			}
		}
	}

	fmt.Println("Ready to craft item...")

	var contentCode string
	if item.Craft != nil {
		contentCode = item.Craft.Skill
	}

	coords := c.GetCoordinatesByCode(contentCode)
	if _, err := c.MoveCharacter(characterName, coords[0].X, coords[0].Y); err != nil {
		return nil, fmt.Errorf("moving to bank: %w", err)
	}

	if err := c.Craft(characterName, code, quantity); err != nil {
		return nil, fmt.Errorf("crafting final item: %w", err)
	}

	fmt.Println("Successfully crafted item!")

	return &item, nil
}

func (c *Svc) Craft(characterName, code string, quantity int) error {
	fmt.Printf("%s crafting %s!\n", characterName, code)
	craftingResp, err := c.Client.CraftItem(characterName, code, quantity)
	if err != nil {
		return fmt.Errorf("crafting item: %w", err)
	}
	fmt.Printf("received %v", craftingResp.Details.Items)
	c.Characters[characterName] = &craftingResp.Character
	c.Characters[characterName].WaitForCooldown()
	return nil
}

func (c *Svc) GatherLoop(characterName, code string, quantity int) error {
	fmt.Printf("%s gathering %s\n", characterName, code)
	item := c.GetItem(code)

	if err := c.Gather(characterName, item, quantity); err != nil { // 8 resources = 1 useful item
		return fmt.Errorf("gathering %s: %w", code, err)
	}

	_, q := c.Characters[characterName].FindItemInInventory(code)
	inventorySlot := InventorySlot{
		Code:     code,
		Quantity: q,
	}
	if err := c.DepositBank(characterName, inventorySlot); err != nil {
		return fmt.Errorf("depositing %d %s: %w", 8, code, err)
	}
	c.Characters[characterName].WaitForCooldown()

	//for i := 0; i < quantity; i++ {
	//	fmt.Printf("gather loop %d\n", i)
	//	if err := c.Gather(characterName, item, 8); err != nil { // 8 resources = 1 useful item
	//		return fmt.Errorf("gathering %s: %w", code, err)
	//	}
	//
	//	_, q := c.Characters[characterName].FindItemInInventory(code)
	//	inventorySlot := InventorySlot{
	//		Code:     code,
	//		Quantity: q,
	//	}
	//	if err := c.DepositBank(characterName, inventorySlot); err != nil {
	//		return fmt.Errorf("depositing %d %s: %w", 8, code, err)
	//	}
	//	c.Characters[characterName].WaitForCooldown()
	//}

	return nil
}

func (c *Svc) Gather(characterName string, item CraftableItem, quantity int) error {
	fmt.Printf("%s gathering %d %s\n", characterName, quantity, item.Name)
	resourceData := c.GetResourceByCode(item.Code)

	fmt.Printf("Gathering %v\n", item)
	// find location of item
	coords := c.GetCoordinatesByCode(resourceData[0].Code)
	if _, err := c.MoveCharacter(characterName, coords[0].X, coords[0].Y); err != nil {
		return fmt.Errorf("moving to bank: %w", err)
	}

	for i := 0; i < quantity; i++ {
		if c.GetCharacterByName(characterName).IsInventoryFull() {
			if err := c.DepositAllItems(characterName); err != nil {
				return fmt.Errorf("depositing inventory: %w", err)
			}

			// find location of item
			if _, err := c.MoveCharacter(characterName, coords[0].X, coords[0].Y); err != nil {
				return fmt.Errorf("moving to bank: %w", err)
			}
		}

		// gather item
		if err := c.gather(characterName); err != nil {
			return fmt.Errorf("attempting to gather %s #%d: %w", item.Name, i, err)
		}
	}

	return nil
}

func (c *Svc) gather(characterName string) error {
	gatherResp, err := c.Client.Gather(characterName)
	if err != nil {
		return fmt.Errorf("gathering: %w", err)
	}
	fmt.Printf("received %v", gatherResp.Details.Items)

	c.Characters[characterName] = &gatherResp.Character
	c.Characters[characterName].WaitForCooldown()

	return nil
}
