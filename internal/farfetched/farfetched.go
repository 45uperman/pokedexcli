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

type PokeArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []PokeEncounter `json:"pokemon_encounters"`
}

type PokeEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
}

type Pokemon struct {
	Name           string
	URL            string
	ID             int `json:"id"`
	Height         int `json:"height"`
	Weight         int `json:"weight"`
	BaseExperience int `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type Pokedex struct {
	NumPokemon int
	NumCaught  int
	Pokemon    map[string]Pokemon
}

func (pd Pokedex) Catch(pkmn Pokemon) {
	pd.Pokemon[pkmn.Name] = pkmn
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

func NewPokedex() Pokedex {
	dex := Pokedex{
		NumPokemon: 1025,
		NumCaught:  0,
		Pokemon:    make(map[string]Pokemon, 5),
	}
	return dex
}

func BuildPage(data []byte) (PokePage, error) {
	var page PokePage
	err := json.Unmarshal(data, &page)
	if err != nil {
		return PokePage{}, fmt.Errorf("error unmarshaling data: %w", err)
	}
	return page, nil
}

func BuildPokemon(name, url string, data []byte) (Pokemon, error) {
	pkmn := Pokemon{
		Name: name,
		URL:  url,
	}
	err := json.Unmarshal(data, &pkmn)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error unmarshaling data: %w", err)
	}
	return pkmn, nil
}

func ExploreArea(data []byte) ([]PokeEncounter, error) {
	var area PokeArea
	err := json.Unmarshal(data, &area)
	if err != nil {
		return []PokeEncounter{}, fmt.Errorf("error unmarshaling data: %w", err)
	}
	return area.PokemonEncounters, nil
}
