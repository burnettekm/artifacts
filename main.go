package main

import (
	"artifacts/api"
	"errors"
	"flag"
	"os"
	"sync"
)

func main() {
	// process flags
	itemPtr := flag.String("item", "", "provide item code that you wish to craft")
	////fightMonsterPtr := flag.String("monster", "", "provide the monster you wish to fight")
	flag.Parse()

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
			}
		}(character)
	}
	wg.Wait()

	for _, character := range service.GetAllCharacters() {
		if character.Name == "Kristi" {
			continue
		}
		go func(characterName string) {
			if err := service.FightForCrafting(characterName, "feather", nil); err != nil {
				panic(err)
			}
		}(character.Name)
	}

	//for service.Characters["Kristi"].WeaponcraftingLevel < 5 {
	//	_, err := service.CraftItem("Kristi", *itemPtr, 1)
	//	if err != nil {
	//		panic(err)
	//	}
	//	if err := service.RecycleItems("Kristi"); err != nil {
	//		panic(err)
	//	}
	//}

	_, err = service.CraftItem("Kristi", *itemPtr, 20)
	if err != nil {
		panic(err)
	}
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
