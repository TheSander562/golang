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
	"strconv"
	"strings"
)

var (
	pageNumber  = 1
	searchInput = ""
	inputNumber int
	filePath    = "settings.yml"
	searchSet   = false
)

type Config struct {
	// Maak een config structuur met alle variabelen uit de settings yml
	Api string `yaml:"api-key"` // Eerst de config.Api specificeren, dan wat de variabele is (string) daarna inladen
}

func pagination(result *tmdb.SearchMulti) {
	// Voor elke pagina, tussen 1 en totaal aantal van de zoekopdracht doe het volgende
	for page := int64(pageNumber); page <= result.TotalPages; {
		// Als de zoekopdracht 1 pagina heeft dan:
		if result.TotalPages == 1 {
			fmt.Println("Voer het nummer in van de rij die je wilt kiezen:")

			for {
				inputs()
				if 1 <= inputNumber && inputNumber <= 20 {
					fmt.Printf("Nummer %d gekozen! \n", inputNumber)
					return
				} else {
					log.Println("Voer een getal uit de lijst in!")
				}
			}
			//Als totaal pagina's groter is dan 1 en het totaal aantal pagina's niet overeenkomt met de huidige pagina
		} else if result.TotalPages > 1 && result.TotalPages != int64(pageNumber) {
			// Als je op de eerste pagina bent dan
			if pageNumber == 1 {
				fmt.Println("Kies: 1 voor huidige pagina - 2 voor volgende pagina")

				for {
					inputs()
					if inputNumber == 1 {
						fmt.Printf("Nummer %d gekozen! \n", inputNumber)
						return
					} else if inputNumber == 2 {
						pageNumber++
						outputs(search())
						break
					} else {
						log.Println("Graag 1 of 2 invoeren!")
					}
				}
			} else {
				fmt.Println("Kies: 0 voor vorige pagina - 1 voor huidige pagina - 2 voor volgende pagina")

				for {
					inputs()
					if inputNumber == 0 {
						pageNumber--
						outputs(search())
						break
					} else if inputNumber == 1 {
						fmt.Printf("Nummer %d gekozen! \n", inputNumber)
						return
					} else if inputNumber == 2 {
						pageNumber++
						outputs(search())
						break
					} else {
						log.Println("Graag 0, 1 of 2 invoeren!")
					}
				}
			}
			// Als de pagina overeenkomt met de laatste pagina dan
		} else if int64(pageNumber) == result.TotalPages {
			fmt.Println("Kies: 0 voor vorige pagina - 1 voor huidige pagina")

			for {
				inputs()
				if inputNumber == 0 {
					pageNumber--
					outputs(search())
					break
				} else if inputNumber == 1 {
					fmt.Printf("Nummer %d gekozen! \n", inputNumber)
					return
				} else {
					log.Println("Graag 0 of 1 invoeren!")
				}
			}
		}
	}
}

func parseConfig() {
	// YML-parser voor de config
	filename, _ := filepath.Abs(filePath)
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	// API key TMDB als environment opslaan in err en als err niet nul dan terug
	err = os.Setenv("APIKey", config.Api)
	if err != nil {
		return
	}
}

func search() *tmdb.SearchMulti {
	// Initialiseer met APIKey, als err niet leeg is, dan print error en exit
	tmdbClient, err := tmdb.Init(os.Getenv("APIKey"))
	if err != nil {
		log.Fatal(err)
	}

	// Opties opgeslagen als map (array), in de map komen strings
	options := map[string]string{
		"language": "nl-NL",
		"page":     fmt.Sprintf("%d", pageNumber),
	}

	// Zoekopdracht maken
	searchResult, err := tmdbClient.GetSearchMulti(searchInput, options)
	if err != nil {
		log.Fatal(err)
	}

	return searchResult
}

func inputs() {
	for {
		// Input voor zoekopdracht (met Bufio is het mogelijk voor spaties)
		input := bufio.NewReader(os.Stdin)
		inputText, _ := input.ReadString('\n')

		if inputText == "\n" {
			log.Println("U moet iets invoeren om door te gaan!")
		} else if searchSet == true {
			trimInput := strings.Trim(inputText, "\n")
			inputNumber, _ = strconv.Atoi(trimInput)
			break
		} else {
			searchInput = inputText
			searchSet = true
			break
		}
	}
}

func outputs(resultaat *tmdb.SearchMulti) {
	if resultaat.TotalResults == 0 {
		log.Fatal("Geen resultaten gevonden voor deze zoekopdracht!")
	}

	// Print resultaten in de range van de resultaten
	fmt.Println("Resultaten - Pagina", resultaat.Page, "-", resultaat.TotalPages, ":")
	for i, v := range resultaat.Results {
		if v.MediaType == "movie" {
			fmt.Println(i+1, "- Film:", v.Title)
		} else if v.MediaType == "tv" {
			fmt.Println(i+1, "- Serie:", v.Name)
		} else if v.MediaType == "person" {
			fmt.Println(i+1, "- Persoon:", v.Name)
		}
	}
	fmt.Println("=====")
}

func main() {
	// Instellen van error code prefix en verberg de tijdstempel
	log.SetPrefix("Foutmelding: ")
	log.SetFlags(0)

	parseConfig()

	// Print tekst
	fmt.Println("Zoek naar films/series/acteurs:")

	inputs()

	result := search()

	outputs(result)

	pagination(result)
}
