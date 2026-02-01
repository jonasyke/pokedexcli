package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type config struct {
	nextURL *string
	prevURL *string
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
	name	string
	description string
	callback	func(*config) error
}

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	words := strings.Fields(lowercase)
	return words
}

func commandHelp(conf *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:" )
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(conf *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func fetchLocations(url string) (RespShallowLocations, error) {
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

    locationResp := RespShallowLocations{}
    err = json.Unmarshal(data, &locationResp)
    if err != nil {
        return RespShallowLocations{}, err
    }

    return locationResp, nil
}

func commandMap(conf *config) error {
    url := "https://pokeapi.co/api/v2/location-area"
    if conf.nextURL != nil {
        url = *conf.nextURL
    }

    locationResp, err := fetchLocations(url)
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

func commandMapB(conf *config) error {
    if conf.prevURL == nil {
        fmt.Println("you're on the first page")
        return nil
    }

    locationResp, err := fetchLocations(*conf.prevURL)
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

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "displays the name of 20 location areas in the Pokemon World",
			callback:commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays the previous 20 location areas",
			callback: commandMapB,
		},
	}
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		nextURL: nil,
		prevURL: nil,
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
		command, exists := commands[commandName]
		if exists {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}
