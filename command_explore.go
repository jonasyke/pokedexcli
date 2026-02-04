package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func commandExplore(conf *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a location area name")
	}
	areaName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	fmt.Printf("Exploring %s...\n", areaName)

	var data []byte
	if val, ok := conf.pokeapiClient.Get(url); ok {
		data = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		conf.pokeapiClient.Add(url, data)
	}

	dest := RespLocationArea{}
	err := json.Unmarshal(data, &dest)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range dest.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
