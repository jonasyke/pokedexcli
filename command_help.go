package main

import (
	"fmt"
)

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
		"inspect",
		"exit",
	}

	commands := getCommands()

	for _, name := range keys {
		cmd := commands[name]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil

}
