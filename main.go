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

	//wg2 := sync.WaitGroup{}
	//for _, character := range service.GetAllCharacters() {
	//	wg2.Add(1)
	//	if character.Name == "Kristi" {
	//		go func(characterName string) {
	//			defer wg2.Done()
	//			if err := service.FightForCrafting(characterName, "red_slimeball", nil); err != nil {
	//				panic(err)
	//			}
	//		}(character.Name)
	//		continue
	//	}
	//	go func(characterName string) {
	//		defer wg2.Done()
	//		for service.GetCharacterByName(characterName).GearcraftingLevel < 5 {
	//			_, err := service.CraftItem(characterName, "copper_helmet", 1)
	//			if err != nil {
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
	//		}
	//	}(character.Name)
	//}
	//
	//wg2.Wait()

	//if _, err := service.CraftItem("Gatherer1", "copper", 1); err != nil {
	//	panic(err)
	//}

	wg3 := sync.WaitGroup{}
	for _, character := range service.GetAllCharacters() {
		wg3.Add(1)
		go func(characterName string) {
			defer wg3.Done()
			if err := service.GatherLoop(characterName, "copper_ore"); err != nil {
				panic(err)
			}
		}(character.Name)
	}
	wg3.Wait()

	//for service.Characters["Kristi"].WeaponcraftingLevel < 5 {
	//	_, err := service.CraftItem("Kristi", *itemPtr, 1)
	//	if err != nil {
	//		panic(err)
	//	}
	//	if err := service.RecycleItems("Kristi"); err != nil {
	//		panic(err)
	//	}
	//}

	//_, err = service.CraftItem("Kristi", *itemPtr, 1)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// equip item
	//if err := service.Equip("Kristi", *item); err != nil {
	//	panic(err)
	//}

	// find chicken
	//contentCode := "chicken"
	//contentType := "monster"

	// find task master
	//contentCode := "chicken"
	//contentType := "tasks_master"
	//maps, err := client.GetMaps(nil, &contentType)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("Found %d %s\n", len(maps.Data), contentType)
	//fmt.Println(maps.Data)
	//
	//x := maps.Data[0].X
	//y := maps.Data[0].Y
	//
	//// move
	//_, err = charSvc.MoveCharacter(x, y)
	//if err != nil {
	//	panic(err)
	//}

	//if _, err := charSvc.AcceptTask(); err != nil {
	//	panic(err)
	//}

	//if _, err := charSvc.CompleteTask(); err != nil {
	//	panic(err)
	//}

	//if err := charSvc.ContinuousFightLoop(); err != nil {
	//	panic(err)
	//}
	//err = charSvc.Rest()
	//if err != nil {
	//	panic(err)
	//}
}
