package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"math/rand"

	"github.com/jonasyke/pokedexcli/internal/pokecache"
)

type config struct {
	pokeapiClient pokecache.Cache
	pokedex map[string]RespPokemon
	nextURL       *string
	prevURL       *string
}

type RespPokemon struct {
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat struct {
			Name string `json:"name"`
		} `json:"stat"`
	}`json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		}`json:"type"`
	}`json:"types"`
}

type RespShallowLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type RespLocationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	words := strings.Fields(lowercase)
	return words
}

func commandPokedex(conf *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	if len(conf.pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. go catch some Pokemon!")
		return nil
	}
	for name := range conf.pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandHelp(conf *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	
	keys := []string{
		"help",
		"map",
		"mapb",
		"explore",
		"catch",
		"exit",
	}

	commands := getCommands()

	for _, name := range keys {
		cmd := commands[name]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil

}

func commandExit(conf *config, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func fetchLocations(url string, cache *pokecache.Cache) (RespShallowLocations, error) {
	if val, ok := cache.Get(url); ok {
		fmt.Println("(using cache)")
		locationResp := RespShallowLocations{}
		err := json.Unmarshal(val, &locationResp)
		if err != nil {
			return RespShallowLocations{}, err
		}
		return locationResp, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return RespShallowLocations{}, fmt.Errorf("bad status code: %v", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return RespShallowLocations{}, err
	}

	cache.Add(url, data)

	locationResp := RespShallowLocations{}
	err = json.Unmarshal(data, &locationResp)
	if err != nil {
		return RespShallowLocations{}, err
	}

	return locationResp, nil
}

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

func commandCatch(conf *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	name := args[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

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

	pokemon := RespPokemon{}
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return err
	}

	res := rand.Intn(pokemon.BaseExperience)

	if res > 40 {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}

	fmt.Printf("%s was caught!\n", name)
	conf.pokedex[name] = pokemon 

	return nil
}  

func commandInspect(conf *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	name := args[0]

	pokemon, ok := conf.pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "displays the name of 20 location areas in the Pokemon World",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapB,
		},
		"explore": {
			name: "explore",
			description: "see a list of all pokemon located in map",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "attempt to catch a pokemon and add it to your pokedex",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect <pokemon_name>",
			description: "View details of a caught pokemon",
			callback: commandInspect,
		},
	}
}

func startRepl() {
	cache := pokecache.NewCache(5 * time.Minute)
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		pokeapiClient: cache,
		pokedex: make(map[string]RespPokemon),
		nextURL:       nil,
		prevURL:       nil,
	}

	commands := getCommands()

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())

		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := words[1:]

		command, exists := commands[commandName]
		if exists {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}
