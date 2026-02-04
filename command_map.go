package main

import (
	"fmt"
)

func commandMap(conf *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if conf.nextURL != nil {
		url = *conf.nextURL
	}

	locationResp, err := fetchLocations(url, &conf.pokeapiClient)
	if err != nil {
		return err
	}

	conf.nextURL = locationResp.Next
	conf.prevURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapB(conf *config, args ...string) error {
	if conf.prevURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationResp, err := fetchLocations(*conf.prevURL, &conf.pokeapiClient)
	if err != nil {
		return err
	}

	conf.nextURL = locationResp.Next
	conf.prevURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
