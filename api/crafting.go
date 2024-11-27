package api

import (
	"fmt"
)

func (c *Svc) CraftItem(code string, quantity int) (*CraftableItem, error) {
	fmt.Printf("Attempting to craft item %s", code)

	item, err := c.Client.GetItem(code)
	if err != nil {
		return nil, fmt.Errorf("getting item: %w", err)
	}

	// verify character can craft item
	if !c.Character.AbleToCraft(item.Craft.Skill, item.Craft.Level) {
		return nil, fmt.Errorf("unable to craft item: required level: %d", item.Craft.Level)
	}

	// get dependent items
	requiredItems := map[*CraftableItem]int{}
	for _, subItem := range item.Craft.Items {
		remainingQuantity := subItem.Quantity * quantity
		craftable, err := c.Client.GetItem(subItem.Code)
		if err != nil {
			return nil, fmt.Errorf("getting item: %s: %w", subItem.Code, err)
		}

		// check if item equipped
		if c.Character.IsEquipped(*craftable) {
			if err := c.unequip(*craftable); err != nil {
				return nil, fmt.Errorf("unequipping item for crafting: %w", err)
			}
			continue
		}

		// check if item in inventory
		found, inventoryQuantity := c.Character.FindItemInInventory(subItem.Code)
		if found && inventoryQuantity >= remainingQuantity {
			continue
		}
		remainingQuantity -= inventoryQuantity

		// check for item in bank
		// assume we're only looking for 1 item at a time
		bankQuantity, err := c.WithdrawFromBankIfFound(subItem.Code, remainingQuantity)
		if err != nil {
			return nil, fmt.Errorf("withdrawing %s from bank if found: %w", subItem.Code, err)
		}
		remainingQuantity -= bankQuantity

		if remainingQuantity <= 0 {
			continue
		}

		if !c.Character.AbleToCraft(craftable.Craft.Skill, craftable.Craft.Level) {
			return nil, fmt.Errorf("unable to craft subitem: %s: needs %s level: %d", craftable.Name, craftable.Craft.Skill, craftable.Craft.Level)
		}

		requiredItems[craftable] = remainingQuantity
	}

	// let's assume we're gathering for now
	fmt.Printf("Gathering required items to craft %s\n", code)
	for reqItem, q := range requiredItems {
		if len(reqItem.Craft.Items) == 0 {
			if err := c.Gather(*reqItem, q); err != nil {
				return nil, fmt.Errorf("gathering required item: %s: %w", reqItem.Code, err)
			}
		} else {
			if _, err := c.CraftItem(reqItem.Code, q); err != nil {
				return nil, fmt.Errorf("crafting subitem %s: %w", reqItem.Code, err)
			}
		}
	}

	fmt.Println("Ready to craft item...")

	contentType := "workshop"
	contentCode := item.Craft.Skill
	mapResp, err := c.Client.GetMaps(&contentCode, &contentType)
	if err != nil {
		return nil, fmt.Errorf("getting location of workshop: %w", err)
	}

	if _, err := c.moveCharacter(mapResp.Data[0].X, mapResp.Data[0].Y); err != nil {
		return nil, fmt.Errorf("moving to workshop location: %w", err)
	}

	if err := c.craft(code, quantity); err != nil {
		return nil, fmt.Errorf("crafting final item: %w", err)
	}

	//if err := c.equip(*item); err != nil {
	//	return fmt.Errorf("equipping crafted item: %w", err)
	//}

	fmt.Println("Successfully crafted item!")
	return item, nil
}

func (c *Svc) craft(code string, quantity int) error {
	fmt.Printf("Crafting %s!\n", code)
	craftingResp, err := c.Client.CraftItem(c.Character.Name, code, quantity)
	if err != nil {
		return fmt.Errorf("crafting item: %w", err)
	}
	fmt.Printf("received %v", craftingResp.Details.Items)
	c.Character = &craftingResp.Character
	c.Character.WaitForCooldown()
	return nil
}

func (c *Svc) moveCharacter(x, y int) (*MoveResponse, error) {
	if c.Character.X == x && c.Character.Y == y {
		fmt.Printf("character already at %d, %d\n", x, y)
		return nil, nil
	}

	moveResp, err := c.Client.MoveCharacter(c.Character.Name, x, y)
	if err != nil {
		return nil, fmt.Errorf("moving character: %w", err)
	}

	c.Character = &moveResp.Data.Character
	c.Character.WaitForCooldown()

	return moveResp, nil
}

func (c *Svc) unequip(item CraftableItem) error {
	fmt.Printf("Unquipping item: %s\n", item.Name)
	unequipResp, err := c.Client.Unequip(c.Character.Name, item)
	if err != nil {
		return fmt.Errorf("unequipping item: %w", err)
	}
	c.Character = &unequipResp.Character
	c.Character.WaitForCooldown()
	return nil
}

func (c *Svc) equip(item CraftableItem) error {
	fmt.Printf("Equipping item: %s\n", item.Name)
	equipResp, err := c.Client.Equip(c.Character.Name, item)
	if err != nil {
		return fmt.Errorf("equipping item: %w", err)
	}
	c.Character = &equipResp.Character
	c.Character.WaitForCooldown()
	return nil
}

func (c *Svc) Gather(item CraftableItem, quantity int) error {
	resourceData, err := c.Client.GetResource(item.Code)
	if err != nil {
		return fmt.Errorf("getting resource data: %w", err)
	}

	fmt.Printf("Gathering %v\n", item)
	// find location of item
	contentType := "resource"
	mapResp, err := c.Client.GetMaps(&resourceData[0].Code, &contentType)
	if err != nil {
		return fmt.Errorf("finding item: %s: %w", resourceData[0].Code, err)
	}

	// move to item
	if _, err := c.moveCharacter(mapResp.Data[0].X, mapResp.Data[0].Y); err != nil {
		return fmt.Errorf("moving to item: %w", err)
	}

	for i := 0; i < quantity; i++ {
		// gather item
		if err := c.gather(); err != nil {
			return fmt.Errorf("attempting to gather %s #%d: %w", item.Name, i, err)
		}
	}

	return nil
}

func (c *Svc) gather() error {
	gatherResp, err := c.Client.Gather(c.Character.Name)
	if err != nil {
		return fmt.Errorf("gathering: %w", err)
	}
	fmt.Printf("received %v", gatherResp.Details.Items)

	c.Character = &gatherResp.Character
	c.Character.WaitForCooldown()

	return nil
}
