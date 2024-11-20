package main

import (
	"artifacts/api"
	"errors"
	"fmt"
	"os"
)

func main() {
	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		panic(errors.New("API_TOKEN ENV VAR not found"))
	}

	client := api.NewClient(token)
	char, err := client.GetCharacter("Kristi")
	if err != nil {
		panic(err)
	}

	charSvc := api.NewCharacterSvc(client, &char.Character)

	// find chicken
	//contentCode := "chicken"
	//contentType := "monster"

	// find task master
	//contentCode := "chicken"
	contentType := "tasks_master"
	maps, err := client.GetMaps(nil, &contentType)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d %s\n", len(maps.Data), contentType)
	fmt.Println(maps.Data)

	x := maps.Data[0].X
	y := maps.Data[0].Y

	// move
	_, err = charSvc.MoveCharacter(x, y)
	if err != nil {
		panic(err)
	}

	//if _, err := charSvc.AcceptTask(); err != nil {
	//	panic(err)
	//}

	if _, err := charSvc.CompleteTask(); err != nil {
		panic(err)
	}

	//if err := charSvc.ContinuousFightLoop(); err != nil {
	//	panic(err)
	//}
	//err = charSvc.Rest()
	//if err != nil {
	//	panic(err)
	//}
}
