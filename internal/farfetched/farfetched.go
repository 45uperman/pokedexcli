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

const PokeURL string = "https://pokeapi.co/api/v2/"

const LocationArea string = "location-area"

func Farfetch[T any](objPtr *T, path string) error {
	// Fetches the data at pokeURL/path and streams it to obj if
	// T is an implemented struct, otherwise it returns an error
	if objPtr == nil {
		return fmt.Errorf("Farfetch received nil pointer")
	}

	if objType := reflect.TypeOf(*objPtr); objType.Kind() != reflect.Struct {
		return fmt.Errorf("Farfetch requires obj to be a pointer to an implemented struct like PokePage, not %v", objType)
	}

	res, err := http.Get(path)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(objPtr)
	if err != nil {
		return err
	}

	return nil
}
