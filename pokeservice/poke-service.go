// service.go

package pokeservice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SpriteData struct {
	FrontDefault string `json:"front_default"`
	Other        struct {
		DreamWorld struct {
			FrontDefault string `json:"front_default"`
			FrontFemale  string `json:"front_female"`
		} `json:"dream_world"`
		Home struct {
			FrontDefault     string `json:"front_default"`
			FrontFemale      string `json:"front_female"`
			FrontShiny       string `json:"front_shiny"`
			FrontShinyFemale string `json:"front_shiny_female"`
		} `json:"home"`
		OfficialArtwork struct {
			FrontDefault string `json:"front_default"`
			FrontShiny   string `json:"front_shiny"`
		} `json:"official-artwork"`
	} `json:"other"`
}
type AbilityObj struct {
	Ability Ability `json:"ability"`
}

type Ability struct {
	Name string `json:"name"`
}

type TypeObj struct {
	Type Type `json:"type"`
}

type Type struct {
	Name string `json:"name"`
}

// Pokemon represents a Pokémon.
type Pokemon struct {
	Name        string       `json:"name"`
	ID          int          `json:"id"`
	Description int          `json:"description"`
	Abilities   []AbilityObj `json:"abilities"`
	Types       []TypeObj    `json:"types"`
	Sprites     SpriteData   `json:"sprites"`
}

// GetAllPokemon fetches a list of Pokémon from the PokeAPI with complete data including sprites, abilities, and types.
func GetAllPokemon(page int) ([]*Pokemon, error) {
	offset := (page - 1) * 20

	apiUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon?offset=%d&limit=%d", offset, 20)

	response, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned a non-OK status code: %v", response.Status)
	}

	var result struct {
		Results []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	var pokemonList []*Pokemon

	for _, pokemonInfo := range result.Results {
		pokemonData, err := GetPokemonByName(pokemonInfo.Name)
		if err != nil {
			return nil, err
		}
		pokemonList = append(pokemonList, pokemonData)
	}

	return pokemonList, nil
}

func GetPokemonByName(name string) (*Pokemon, error) {
	apiUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

	response, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned a non-OK status code: %v", response.Status)
	}

	var pokemon Pokemon
	if err := json.NewDecoder(response.Body).Decode(&pokemon); err != nil {
		return nil, err
	}

	return &pokemon, nil
}
