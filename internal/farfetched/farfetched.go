package farfetched

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

type FetchPath string

const POKEMAP FetchPath = "location-area"

const pokeURL string = "https://pokeapi.co/"

func Farfetch[T any](obj T, path FetchPath) error {
	// Fetches the data at pokeURL/path and streams it to obj if
	// T is an implemented struct, otherwise it returns an error

	if objType := reflect.TypeOf(obj); objType.Kind() != reflect.Struct {
		return fmt.Errorf("Farfetch requires obj to be an implemented struct like PokePage, not %v", objType)
	}

	res, err := http.Get(pokeURL + string(path))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&obj)
	if err != nil {
		return err
	}

	return nil
}
