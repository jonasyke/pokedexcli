package main

import (
	"fmt"
)

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
