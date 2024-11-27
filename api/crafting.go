package api

import (
	"fmt"
)

func (c *Svc) CraftItem(code string, quantity int) (*CraftableItem, error) {
	fmt.Printf("Attempting to craft item %s\n", code)

	item, err := c.Client.GetItem(code)
	if err != nil {
		return nil, fmt.Errorf("getting item: %w", err)
	}

	charName := "Kristi"

	// verify character can craft item
	if !c.Characters[charName].AbleToCraft(item.Craft.Skill, item.Craft.Level) {
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
		if c.Characters[charName].IsEquipped(*craftable) {
			if err := c.unequip(*craftable); err != nil {
				return nil, fmt.Errorf("unequipping item for crafting: %w", err)
			}
			continue
		}

		// check if item in inventory
		found, inventoryQuantity := c.Characters[charName].FindItemInInventory(subItem.Code)
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

		if !c.Characters[charName].AbleToCraft(craftable.Craft.Skill, craftable.Craft.Level) {
			return nil, fmt.Errorf("unable to craft subitem: %s: needs %s level: %d", craftable.Name, craftable.Craft.Skill, craftable.Craft.Level)
		}

		requiredItems[craftable] = remainingQuantity
	}

	// let's assume we're gathering for now
	fmt.Printf("Gathering required items to craft %s\n", code)
	for reqItem, q := range requiredItems {
		if len(reqItem.Craft.Items) == 0 {
			go c.Gatherer(reqItem.Code, q)
			if err := c.WaitForMaterials(reqItem); err != nil {
				return nil, fmt.Errorf("waiting for materials: %w", err)
			}
		}
		if _, err := c.CraftItem(reqItem.Code, q); err != nil {
			return nil, fmt.Errorf("crafting item %s: %w", reqItem.Code, err)
		}
	}

	fmt.Println("Ready to craft item...")

	contentType := "workshop"
	contentCode := item.Craft.Skill
	mapResp, err := c.Client.GetMaps(&contentCode, &contentType)
	if err != nil {
		return nil, fmt.Errorf("getting location of workshop: %w", err)
	}
	fmt.Println(mapResp)

	if err := c.MoveCharacter(charName, mapResp.Data[0].X, mapResp.Data[0].Y); err != nil {
		return nil, fmt.Errorf("moving to workshop location: %w", err)
	}

	if err := c.craft(code, quantity); err != nil {
		return nil, fmt.Errorf("crafting final item: %w", err)
	}

	fmt.Println("Successfully crafted item!")
	return item, nil
}

func (c *Svc) WaitForMaterials(mat *CraftableItem) error {
	<-c.ResourceChannels[mat.Craft.Skill]
	if _, err := c.CraftItem(mat.Code, 1); err != nil {
		return fmt.Errorf("crafting subitem %s: %w", mat.Code, err)
	}
	return nil
}

func (c *Svc) craft(code string, quantity int) error {
	fmt.Printf("Crafting %s!\n", code)
	charName := "Kristi"
	craftingResp, err := c.Client.CraftItem(charName, code, quantity)
	if err != nil {
		return fmt.Errorf("crafting item: %w", err)
	}
	fmt.Printf("received %v", craftingResp.Details.Items)
	c.Characters[charName] = &craftingResp.Character
	c.Characters[charName].WaitForCooldown()
	return nil
}

func (c *Svc) MoveCharacter(charName string, x, y int) error {
	if c.Characters[charName].X == x && c.Characters[charName].Y == y {
		fmt.Printf("character already at %d, %d\n", x, y)
		return nil
	}

	moveResp, err := c.Client.MoveCharacter(charName, x, y)
	if err != nil {
		return fmt.Errorf("moving character: %w", err)
	}

	c.Characters[charName] = &moveResp.Data.Character
	c.Characters[charName].WaitForCooldown()

	return nil
}

func (c *Svc) unequip(item CraftableItem) error {
	fmt.Printf("Unquipping item: %s\n", item.Name)
	charName := "Kristi"
	unequipResp, err := c.Client.Unequip(charName, item)
	if err != nil {
		return fmt.Errorf("unequipping item: %w", err)
	}
	c.Characters[charName] = &unequipResp.Character
	c.Characters[charName].WaitForCooldown()
	return nil
}

func (c *Svc) equip(item CraftableItem) error {
	fmt.Printf("Equipping item: %s\n", item.Name)
	charName := "Kristi"
	equipResp, err := c.Client.Equip(charName, item)
	if err != nil {
		return fmt.Errorf("equipping item: %w", err)
	}
	c.Characters[charName] = &equipResp.Character
	c.Characters[charName].WaitForCooldown()
	return nil
}

//func (c *Svc) Gather(item CraftableItem, quantity int) error {
//	resourceData, err := c.Client.GetResource(item.Code)
//	if err != nil {
//		return fmt.Errorf("getting resource data: %w", err)
//	}
//
//	fmt.Printf("Gathering %v\n", item)
//	// find location of item
//	contentType := "resource"
//	mapResp, err := c.Client.GetMaps(&resourceData[0].Code, &contentType)
//	if err != nil {
//		return fmt.Errorf("finding item: %s: %w", resourceData[0].Code, err)
//	}
//
//	// move to item
//	if _, err := c.MoveCharacter(mapResp.Data[0].X, mapResp.Data[0].Y); err != nil {
//		return fmt.Errorf("moving to item: %w", err)
//	}
//
//	for i := 0; i < quantity; i++ {
//		// gather item
//		if err := c.gather(); err != nil {
//			return fmt.Errorf("attempting to gather %s #%d: %w", item.Name, i, err)
//		}
//	}
//
//	return nil
//}

func (c *Svc) gather() error {
	charName := "Woodcutter"
	gatherResp, err := c.Client.Gather(charName)
	if err != nil {
		return fmt.Errorf("gathering: %w", err)
	}
	fmt.Printf("received %v", gatherResp.Details.Items)

	c.Characters[charName] = &gatherResp.Character
	c.Characters[charName].WaitForCooldown()

	return nil
}
