package main

import (
	"artifacts/api"
	"errors"
	"os"
	"sync"
)

func main() {
	// process flags
	//itemPtr := flag.String("item", "", "provide item code that you wish to craft")
	////fightMonsterPtr := flag.String("monster", "", "provide the monster you wish to fight")
	//flag.Parse()

	// set up app dependencies
	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		panic(errors.New("API_TOKEN ENV VAR not found"))
	}

	service, err := api.NewSvc(token)
	if err != nil {
		panic(err)
	}
	//if err := service.RecycleItems("Kristi"); err != nil {
	//	panic(err)
	//}

	wg := sync.WaitGroup{}
	for _, character := range service.GetAllCharacters() {
		wg.Add(1)
		go func(character *api.Character) {
			defer wg.Done()
			for _, invItem := range character.Inventory {
				if invItem.Code == "" {
					continue
				}
				if err := service.DepositBank(character.Name, invItem); err != nil {
					panic(err)
				}
				service.GetCharacterByName(character.Name).WaitForCooldown()
			}
		}(character)
	}
	wg.Wait()

	//	wg2 := sync.WaitGroup{}
	//	for _, character := range service.GetAllCharacters() {
	//		wg2.Add(1)
	//		if character.Name == "Kristi" {
	//			go func(characterName string) {
	//				defer wg2.Done()
	//				for service.GetCharacterByName(characterName).JewelrycraftingLevel < 10 {
	//					item, err := service.CraftItem(characterName, "copper_ring", 1)
	//					if err != nil {
	//						panic(err)
	//					}
	//					invSlot := api.InventorySlot{
	//						Code:     item.Code,
	//						Quantity: 1,
	//					}
	//					if err := service.DepositBank(characterName, invSlot); err != nil {
	//						panic(err)
	//					}
	//					service.GetCharacterByName(characterName).WaitForCooldown()
	//				}
	//				if err := service.GatherLoop(characterName, "iron_ore"); err != nil {
	//					panic(err)
	//				}
	//			}(character.Name)
	//			continue
	//		}
	//		go func(characterName string) {
	//			defer wg2.Done()
	//			if err := service.GatherLoop(characterName, "iron_ore"); err != nil {
	//				panic(err)
	//			}
	//			//
	//			//if service.GetCharacterByName(characterName).HelmetSlot == "" {
	//			//	if err := service.Equip(characterName, *item); err != nil {
	//			//		panic(err)
	//			//	}
	//			//} else {
	//			//	if err := service.RecycleItems(characterName); err != nil {
	//			//		panic(err)
	//			//	}
	//			//	if err := service.DepositAllItems(characterName); err != nil {
	//			//		panic(err)
	//			//	}
	//			//}
	//		}(character.Name)
	//	}
	//	wg2.Wait()
	//}
	wg2 := sync.WaitGroup{}
	for _, character := range service.GetAllCharacters() {
		wg2.Add(1)
		if character.Name == "Kristi" {
			go func(characterName string) {
				defer wg2.Done()
				if err := service.FightForCrafting(characterName, "cowhide", nil); err != nil {
					panic(err)
				}
			}(character.Name)
			continue
		}
		go func(characterName string) {
			defer wg2.Done()
			for service.GetCharacterByName(characterName).WoodcuttingLevel < 10 {
				if _, err := service.CraftItem(characterName, "ash_plank", 5); err != nil {
					panic(err)
				}
			}

			//if service.GetCharacterByName(characterName).HelmetSlot == "" {
			//	if err := service.Equip(characterName, *item); err != nil {
			//		panic(err)
			//	}
			//} else {
			//	if err := service.RecycleItems(characterName); err != nil {
			//		panic(err)
			//	}
			//	if err := service.DepositAllItems(characterName); err != nil {
			//		panic(err)
			//	}
			//}
		}(character.Name)
	}
	wg2.Wait()
}
