package api

import "fmt"

func (c *Svc) Gatherer(itemCode string, quantity int) error {
	item, err := c.Client.GetItem(itemCode)
	if err != nil {
		return fmt.Errorf("getting item: %w", err)
	}

	resourceData, err := c.Client.GetResource(item.Code)
	if err != nil {
		return fmt.Errorf("getting resource data: %w", err)
	}

	// find location of item
	contentType := "resource"
	mapResp, err := c.Client.GetMaps(&resourceData[0].Code, &contentType)
	if err != nil {
		return fmt.Errorf("finding item: %s: %w", resourceData[0].Code, err)
	}

	// move to resource
	if err := c.MoveCharacter("Woodcutter", mapResp.Data[0].X, mapResp.Data[0].Y); err != nil {
		return fmt.Errorf("moving to item: %w", err)
	}

	for i := 0; i < quantity; i++ {
		// gather item
		if err := c.gather(); err != nil {
			return fmt.Errorf("attempting to gather %s #%d: %w", item.Name, i, err)
		}
	}

	invSlot := InventorySlot{
		Code:     item.Code,
		Quantity: quantity,
	}
	if err := c.DepositBank("Woodcutter", invSlot); err != nil {
		return fmt.Errorf("depositing %s: %w", item.Code, err)
	}

	c.ResourceChannels[resourceData[0].Skill] <- item.Code
	return nil
}
