package api

import "fmt"

type SkillResponse struct {
	Data  SkillData    `json:"data"`
	Error ErrorMessage `json:"error"`
}

type SkillData struct {
	Cooldown  Cooldown     `json:"cooldown"`
	Details   SkillDetails `json:"details"`
	Character Character    `json:"character"`
}

func (c *Svc) LevelUpSkill(characterName, skill string, wantLevel int) error {
	character := c.Characters[characterName]
	currSkillLevel := character.FindSkillLevel(skill)
	// break out of loop when goal is reached
	if currSkillLevel == wantLevel {
		return nil
	}

	options, err := c.Client.GetItems(skill, wantLevel, currSkillLevel)
	if err != nil {
		return fmt.Errorf("getting list of items: %w", err)
	}

	// decide which item to start with
	minQuantity := 0
	itemCode := ""
	for _, item := range options {
		resourcesNeeded := sumResources(item)
		if !character.AbleToCraft(item.Craft.Skill, item.Craft.Level) {
			continue
		}
		if minQuantity == 0 || resourcesNeeded <= minQuantity {
			minQuantity = resourcesNeeded
			itemCode = item.Code
		}
	}

	// todo: will this work?
	for currSkillLevel <= wantLevel {
		item, err := c.CraftItem(itemCode, 1)
		if err != nil {
			return fmt.Errorf("crafting item %s: %w", itemCode, err)
		}

		slot := InventorySlot{
			Code:     item.Code,
			Quantity: 1,
		}

		if err := c.DepositBank(characterName, slot); err != nil {
			return fmt.Errorf("depositing crafted item: %w", err)
		}
	}

	return c.LevelUpSkill(characterName, skill, currSkillLevel+1)
}

func sumResources(item CraftableItem) int {
	out := 0
	for _, r := range item.Craft.Items {
		out += r.Quantity
	}
	return out
}
