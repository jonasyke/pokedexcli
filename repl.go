package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jonasyke/pokedexcli/internal/pokecache"
)

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	words := strings.Fields(lowercase)
	return words
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
