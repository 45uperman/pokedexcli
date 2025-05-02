package farfetched

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func PokeGet(path string) ([]byte, error) {
	// Fetches the data at path, reads it, and
	// returns it as bytes.

	res, err := http.Get(path)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

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
