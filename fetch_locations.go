package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jonasyke/pokedexcli/internal/pokecache"
)

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
