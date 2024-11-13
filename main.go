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
	fmt.Println(token)

	url := "https://api.artifactsmmo.com"

	client := api.NewClient(url, token)
	char, err := client.GetCharacter("Kristi")
	if err != nil {
		panic(err)
	}

	fmt.Println(char)
}
