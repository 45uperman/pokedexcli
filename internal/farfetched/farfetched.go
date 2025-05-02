package farfetched

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/45uperman/pokedexcli/internal/pokecache"
)

type PokePage struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

const PokeURL string = "https://pokeapi.co/api/v2/"

const LocationArea string = "location-area"

func PokeGet(url string, c *pokecache.Cache) ([]byte, error) {
	// Fetches the data at url, reads it, and
	// returns it as bytes.

	data, ok := c.Get(url)
	if ok {
		return data, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	c.Add(url, data)

	return data, nil
}

func BuildPage(data []byte) (PokePage, error) {
	var page PokePage
	err := json.Unmarshal(data, &page)
	if err != nil {
		return PokePage{}, fmt.Errorf("error unmarshaling data: %w", err)
	}
	return page, nil
}
