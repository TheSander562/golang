package main

// Import pakketjes, en github link als tmdb
import (
	"fmt"
	tmdb "github.com/cyruzin/golang-tmdb"
	"log"
	"os"
)

func main() {
	// Instellen van error code prefix en verberg de tijdstempel
	log.SetPrefix("Foutmelding: ")
	log.SetFlags(0)

	// Print tekst
	fmt.Println("Zoek naar films/series/acteurs: ")

	// Input voor zoekopdracht
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return
	}

	if input == "" {
		log.Fatal("U moet een zoekopdracht invoeren!")
	}

	// API key TMDB als environment opslaan in err en als err niet nul dan terug
	err = os.Setenv("APIKey", "980006f7a03fa02226232d54d8a76314")
	if err != nil {
		return
	}

	search(input, 1)

}

func search(input string, pagina int) {
	// Initialiseer met APIKey, als err niet leeg is, dan print error en exit
	tmdbClient, err := tmdb.Init(os.Getenv("APIKey"))
	if err != nil {
		log.Fatal(err)
	}

	zoekString := input

	// Opties opgeslagen als map (array), in de map komen strings
	opties := map[string]string{
		"language": "nl-NL",
		"page":     fmt.Sprintf("%d", pagina),
	}

	// Zoekopdracht maken
	zoekopdracht, err := tmdbClient.GetSearchMulti(input, opties)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Resultaten: ")

	// Print resultaten in de range van de resultaten
	for i, v := range zoekopdracht.Results {
		if v.MediaType == "movie" {
			fmt.Println(i+1, "- Film: ", v.Title)
		} else if v.MediaType == "tv" {
			fmt.Println(i+1, "- Serie: ", v.Name)
		} else if v.MediaType == "person" {
			fmt.Println(i+1, "- Persoon: ", v.Name)
		}
	}
	fmt.Println("=====")

	// Print tekst
	fmt.Println("Kies: 1 voor huidige pagina - 2 voor volgende pagina")

	// Input voor zoekopdracht
	_, err = fmt.Scanln(&input)
	if err != nil {
		return
	}
	if input == "1" {
		fmt.Println("Nummer 1 gekozen!")
	} else if input == "2" {
		pagina++
		search(zoekString, pagina)
	} else {
		fmt.Println("U moet getal 1 of 2 invoeren!")
	}

}
