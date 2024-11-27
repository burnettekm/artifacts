package main

import (
	"artifacts/api"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	// process flags
	itemPtr := flag.String("item", "", "provide item code that you want to craft")
	////fightMonsterPtr := flag.String("monster", "", "provide the monster you wish to fight")
	//skillPtr := flag.String("skill", "", "provide skill you want to level up")
	//skillLvlPtr := flag.Int("skillLevel", 1, "provide level you want to achieve")
	flag.Parse()

	// set up app dependencies
	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		panic(errors.New("API_TOKEN ENV VAR not found"))
	}

	client := api.NewClient(token)
	characters, err := client.GetCharacters()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", characters)

	service := api.NewSvc(client, characters)
	for _, character := range characters {
		for _, invItem := range character.Inventory {
			if invItem.Code == "" {
				continue
			}
			if err := service.DepositBank(character.Name, invItem); err != nil {
				panic(err)
			}
		}
	}
	//
	if *itemPtr != "" {
		_, err := service.CraftItem(*itemPtr, 1)
		if err != nil {
			panic(err)

		}
	}
	//
	//	// equip item
	//	if err := service.Equip(*item); err != nil {
	//		panic(err)
	//	}
	//}
	//
	//if *skillPtr != "" {
	//	level := characters.Character.Level + 1
	//	if *skillLvlPtr != 1 {
	//		level = *skillLvlPtr
	//	}
	//	// train skill
	//	err := service.LevelUpSkill(*skillPtr, level)
	//	if err != nil {
	//		panic(err)
	//	}
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
