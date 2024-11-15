package main

import (
	"artifacts/api"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func main() {
	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		panic(errors.New("API_TOKEN ENV VAR not found"))
	}
	fmt.Println(token)

	url := "https://api.artifactsmmo.com"

	client := api.NewClient(url, token)
	//char, err := client.GetCharacter("Kristi")
	//if err != nil {
	//	panic(err)
	//}
	//
	//charBytes, err := json.MarshalIndent(char, " ", "    ")
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Char Resp: ", string(charBytes))

	resp, err := client.Fight("Kristi")
	if err != nil {
		panic(err)
	}

	bytes, err := json.MarshalIndent(resp, " ", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println("Fight resp: ", string(bytes))
}
