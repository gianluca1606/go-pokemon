package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gianluca1606/go-pokemon/pokeservice"
)

type Film struct {
	Title    string
	Director string
}

func getSinglePoke(w http.ResponseWriter, r *http.Request) {
	pokeName := r.URL.Path[len("/poke/"):]

	pokemon, err := pokeservice.GetPokemonByName(pokeName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	layoutData := struct {
		Title       string
		MainContent template.HTML
	}{
		Title:       "Pokémon Details",
		MainContent: renderPokemonCard(pokemon),
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html"))
	if err := tmpl.Execute(w, layoutData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderPokemonCard(pokemon *pokeservice.Pokemon) template.HTML {
	tmpl := template.Must(template.ParseFiles("templates/single-poke.html"))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, pokemon); err != nil {
		// Handle the error
		return template.HTML("<p>Error rendering Pokémon card</p>")
	}
	return template.HTML(buf.String())
}

func loadMorePokemon(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	pokemonList, err := pokeservice.GetAllPokemon(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/poke-list.html"))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		PokemonList []*pokeservice.Pokemon
		NextPage    int
	}{
		PokemonList: pokemonList,
		NextPage:    page + 1,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the rendered HTML as the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(buf.String()))
}

func getHomePage(w http.ResponseWriter, r *http.Request) {
	pokemonList, err := pokeservice.GetAllPokemon(1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/poke-list.html"))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		PokemonList []*pokeservice.Pokemon
		NextPage    int
	}{
		PokemonList: pokemonList,
		NextPage:    2,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title       string
		MainContent template.HTML
	}{
		Title:       "Home Page",
		MainContent: template.HTML(buf.String()),
	}

	renderLayoutTemplate(w, "layout.html", data)
}

func main() {
	fmt.Println("Go app...")

	http.HandleFunc("/poke/", getSinglePoke)
	http.HandleFunc("/load-more/", loadMorePokemon)
	http.HandleFunc("/", getHomePage)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func renderLayoutTemplate(w http.ResponseWriter, tmplPath string, data interface{}) {
	tmpl := template.Must(template.ParseFiles("templates/" + tmplPath))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
