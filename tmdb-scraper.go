package main

// Import pakketjes, en github link als tmdb
import (
	"bufio"
	"fmt"
	tmdb "github.com/cyruzin/golang-tmdb"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	// Maak een config structuur met alle variabelen uit de settings yml
	Api string `yaml:"api-key"` // Eerst de config.Api specificeren, dan wat de variabele is (string) daarna inladen
}

func parseConfig(configFile string) Config {
	// YML-parser voor de config
	filename, _ := filepath.Abs(configFile)
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func search(input string, pagina int) *tmdb.SearchMulti {
	// Initialiseer met APIKey, als err niet leeg is, dan print error en exit
	tmdbClient, err := tmdb.Init(os.Getenv("APIKey"))
	if err != nil {
		log.Fatal(err)
	}

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

	if zoekopdracht.TotalResults == 0 {
		log.Fatal("Geen resultaten gevonden voor deze zoekopdracht!")
	}

	// Print resultaten in de range van de resultaten
	fmt.Println("Resultaten - Pagina", zoekopdracht.Page, "-", zoekopdracht.TotalPages, ":")
	for i, v := range zoekopdracht.Results {
		if v.MediaType == "movie" {
			fmt.Println(i+1, "- Film:", v.Title)
		} else if v.MediaType == "tv" {
			fmt.Println(i+1, "- Serie:", v.Name)
		} else if v.MediaType == "person" {
			fmt.Println(i+1, "- Persoon:", v.Name)
		}
	}
	fmt.Println("=====")

	return zoekopdracht
}

func main() {
	// Instellen van error code prefix en verberg de tijdstempel
	log.SetPrefix("Foutmelding: ")
	log.SetFlags(0)

	// Laad de config file in de config parser functie en zet de return in de variabele config
	config := parseConfig("settings.yml")

	// API key TMDB als environment opslaan in err en als err niet nul dan terug
	err := os.Setenv("APIKey", config.Api)
	if err != nil {
		return
	}

	// Print tekst
	fmt.Println("Zoek naar films/series/acteurs:")

	// Input voor zoekopdracht (met Bufio is het mogelijk voor spaties)
	input := bufio.NewReader(os.Stdin)
	text, _ := input.ReadString('\n')

	if text == "\n" {
		log.Fatal("U moet een zoekopdracht invoeren!")
	}

	// Krijg de waarden terug van de zoekopdracht
	zoekopdracht := search(text, 1)

	// Voor elke pagina, tussen 1 en totaal aantal van de zoekopdracht doe het volgende
	for pagina := int64(1); pagina <= zoekopdracht.TotalPages; {
		// Als de zoekopdracht 1 pagina heeft dan:
		if zoekopdracht.TotalPages == 1 {
			fmt.Println("Voer het nummer in van de rij die je wilt kiezen:")

			// Loop voor als er iets anders gekozen wordt
			for {
				var input string
				_, err = fmt.Scanln(&input)

				// Als er een error wordt gegenereerd
				if err != nil {
					fmt.Println("Voer een getal in!")
				} else if input == "1" {
					fmt.Println("Nummer 1 gekozen!")
					return
				} else {
					fmt.Println("Voer een getal in!")
				}
			}
			//Als totaal pagina's groter is dan 1 en het totaal aantal pagina's niet overeenkomt met de huidige pagina
		} else if zoekopdracht.TotalPages > 1 && int(zoekopdracht.TotalPages) != zoekopdracht.Page {
			// Als je op de eerste pagina bent dan
			if zoekopdracht.Page == 1 {
				fmt.Println("Kies: 1 voor huidige pagina - 2 voor volgende pagina")

				for {
					var input string
					_, err = fmt.Scanln(&input)

					if err != nil {
						fmt.Println("Graag 1 of 2 invoeren!")
					} else if input == "1" {
						fmt.Println("Nummer 1 gekozen!")
						return
					} else if input == "2" {
						pagina++
						zoekopdracht = search(text, int(pagina))
						break
					} else {
						fmt.Println("Graag 1 of 2 invoeren!")
					}
				}
			} else {
				fmt.Println("Kies: 0 voor vorige pagina - 1 voor huidige pagina - 2 voor volgende pagina")

				for {
					var input string
					_, err = fmt.Scanln(&input)

					if err != nil {
						fmt.Println("Graag 0, 1 of 2 invoeren!")
					} else if input == "0" {
						pagina--
						zoekopdracht = search(text, int(pagina))
						break
					} else if input == "1" {
						fmt.Println("Nummer 1 gekozen!")
						return
					} else if input == "2" {
						pagina++
						zoekopdracht = search(text, int(pagina))
						break
					} else {
						fmt.Println("Graag 0, 1 of 2 invoeren!")
					}
				}
			}
			// Als de pagina overeenkomt met de laatste pagina dan
		} else if int64(zoekopdracht.Page) == zoekopdracht.TotalPages {
			fmt.Println("Kies: 0 voor vorige pagina - 1 voor huidige pagina")

			for {
				var input string
				_, err = fmt.Scanln(&input)

				if err != nil {
					fmt.Println("Graag 0 of 1 invoeren!")
				} else if input == "0" {
					pagina--
					zoekopdracht = search(text, int(pagina))
					break
				} else if input == "1" {
					fmt.Println("Nummer 1 gekozen!")
					return
				} else {
					fmt.Println("Graag 0 of 1 invoeren!")
				}
			}
		}
	}
}
