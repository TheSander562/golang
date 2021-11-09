package main

// Import packages, and GitHub link as TMDB
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

const filePath = "settings.yml"

type Config struct {
	// Make config structure with all variabels from settings yml
	Api      string `yaml:"api-key"` // First specify config.Api, then load Api (string) with api-key
	Language string `yaml:"search-language"`
}

func parseConfig() {
	// YML-parser for the config
	fileName, _ := filepath.Abs(filePath)
	// Check if config file exists
	yamlFile, err := ioutil.ReadFile(fileName)
	// If config file does not exist
	if err != nil {
		errorCodes("noConfig")
	}

	// Create variable config and put the read config in the variable
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		errorCodes("invalidConfig")
	}

	// API key TMDB save in environment, if err is not empty, return
	err = os.Setenv("apiKey", config.Api)
	err = os.Setenv("searchLanguage", config.Language)
	if err != nil {
		errorCodes("")
	}
}

func searchTMDB(p int, s string) *tmdb.SearchMulti {
	// Get environment config keys
	searchLanguage := os.Getenv("searchLanguage")
	apiKey := os.Getenv("apiKey")

	// Initialise variable with the APIKey, if api error print error code
	tmdbClient, err := tmdb.Init(apiKey)
	if err != nil {
		errorCodes("apiError")
	}

	// Create variables
	var pageNumber = p
	var searchInput = s

	// Save options as Map (array), the map exists of strings
	options := map[string]string{
		"language": fmt.Sprintf("%s", searchLanguage),
		"page":     fmt.Sprintf("%d", pageNumber),
	}

	// Create the search result with the options specified
	searchResult, err := tmdbClient.GetSearchMulti(searchInput, options)
	if err != nil {
		errorCodes("searchError")
	}

	// return the search results
	return searchResult
}

func searchChosen(a [2][]string, n int) {
	searchArray := a
	userNumber := n - 1
	var (
		movie  *tmdb.MovieDetails
		tv     *tmdb.TVDetails
		person *tmdb.PersonDetails
		err    error
	)

	searchID, _ := strconv.Atoi(searchArray[0][userNumber])

	// Get environment config keys
	searchLanguage := os.Getenv("searchLanguage")
	apiKey := os.Getenv("apiKey")

	// Initialise variable with the APIKey, if api error print error code
	tmdbClient, err := tmdb.Init(apiKey)
	if err != nil {
		errorCodes("apiError")
	}

	// Save options as Map (array), the map exists of strings
	options := map[string]string{
		"language": fmt.Sprintf("%s", searchLanguage),
	}

	// Check if the searched ID is tv, movie or person
	switch searchArray[1][userNumber] {
	case "tv":
		tv, err = tmdbClient.GetTVDetails(searchID, options)
		fmt.Println(tv)
		fmt.Println(tv.Name)
	case "movie":
		movie, err = tmdbClient.GetMovieDetails(searchID, options)
		fmt.Println(movie)
		fmt.Println(movie.Title)
	case "person":
		person, err = tmdbClient.GetPersonDetails(searchID, options)
		fmt.Println(person)
		fmt.Println(person.Name)
	}

	if err != nil {
		errorCodes("searchError")
	}

}

func errorCodes(e string) {
	// If there is any error, check the errorCode string to see which error occurred
	switch e {
	case "noInput":
		log.Println("Do not leave it blank please! Choose what you want to do!")
	case "noSearch":
		log.Println("Do not leave it blank please! You need to search something!")
	case "wrongInput":
		log.Println("This number is not on the list! Try again.")
	case "wrongNumber":
		log.Println("No valid input! Please choose from above options.")
	case "noResults":
		log.Fatal("No results for this search!")
	case "noConfig":
		log.Fatal("Config file not found.")
	case "invalidConfig":
		log.Fatal("Config file is not valid. Please check for any issues.")
	case "apiError":
		log.Fatal("No valid API key given.")
	case "searchError":
		log.Fatal("Input text is not valid.")
	default:
		log.Fatal("Default Error Code")
	}
}

func inputSearch() string {
	for {
		// Input voor zoekopdracht (met Bufio is het mogelijk voor spaties)
		input := bufio.NewReader(os.Stdin)
		inputText, _ := input.ReadString('\n')

		if inputText == "\n" {
			errorCodes("noSearch")
		} else {
			fmt.Println("")
			return inputText
		}
	}
}

func inputNavigation() (int, error) {
	for {
		// Input voor zoekopdracht (met Bufio is het mogelijk voor spaties)
		input := bufio.NewReader(os.Stdin)
		inputText, _ := input.ReadString('\n')

		if inputText == "\n" {
			errorCodes("noInput")
		} else {
			trimInput := strings.Trim(inputText, "\n")
			inputNumber, e := strconv.Atoi(trimInput)

			return inputNumber, e
		}
	}
}

func pageLogic(n int, r *tmdb.SearchMulti) (bool, bool, string) {
	input := n
	results := r
	page := results.Page
	maxPages := int(results.TotalPages)
	maxResults := int(results.TotalResults)

	if input > 20 || maxResults < input {
		errorCodes("wrongInput")
	} else {
		switch maxPages {
		case 1:
			searchSet, searchDetails := onePage(input)
			return searchSet, searchDetails, ""

		default:
			searchSet, searchDetails, searchAction := morePages(input, page, maxPages)
			return searchSet, searchDetails, searchAction
		}
	}

	return true, false, ""
}

func onePage(i int) (bool, bool) {
	input := i

	switch {
	case input >= 1:
		fmt.Println("You chose:", input)
		return true, true

	default:
		return false, false
	}
}

func morePages(i int, p int, m int) (bool, bool, string) {
	input := i
	page := p
	maxPages := m

	if page <= maxPages {
		switch {
		case maxPages > 1 && maxPages != page:
			if page == 1 {
				switch input {
				case 0:
					return false, false, ""

				case 1:
					return true, false, "current"

				case 2:
					return true, false, "next"

				default:
					return true, false, "error"
				}
			} else {
				switch input {
				case 0:
					return false, false, ""

				case 1:
					return true, false, "current"

				case 2:
					return true, false, "previous"

				case 3:
					return true, false, "next"

				default:
					return true, false, "error"
				}
			}
		case page == maxPages:
			switch input {
			case 0:
				return false, false, ""

			case 1:
				return true, false, "current"

			case 2:
				return true, false, "previous"

			default:
				return true, false, "error"
			}
		}
	}

	return true, false, "error"
}

func printResults(r *tmdb.SearchMulti) [2][]string {
	results := r

	// If there are no results, print error
	if results.TotalResults == 0 {
		errorCodes("noResults")
	}

	// Create 2D array, one for ID and one for mediaType
	var searchArray [2][]string

	// Print results in the range of the results
	fmt.Println("Results - Page", results.Page, "-", results.TotalPages, ":")
	for i, v := range results.Results {
		if v.MediaType == "movie" {
			fmt.Println(i+1, "- Movie:", v.Title)
		} else if v.MediaType == "tv" {
			fmt.Println(i+1, "- Series:", v.Name)
		} else if v.MediaType == "person" {
			fmt.Println(i+1, "- Person:", v.Name)
		}

		searchArray[0] = append(searchArray[0], strconv.FormatInt(v.ID, 10))
		searchArray[1] = append(searchArray[1], v.MediaType)
	}
	fmt.Println("======================")

	maxPages := int(results.TotalPages)

	if results.Page <= maxPages {
		switch {
		case maxPages == 1:
			fmt.Println("Enter number of the row to choose - 0 to search again:")
		case maxPages > 1 && maxPages != results.Page:
			if results.Page == 1 {
				fmt.Println("Choose: 0 to search again - 1 for current page - 2 for next page")
			} else {
				fmt.Println("Choose: 0 to search again - 1 for current page - 2 for previous page - 3 for next page")
			}
		case results.Page == maxPages:
			fmt.Println("Choose: 0 to search again - 1 for current page - 2 for previous page")
		}
	}

	return searchArray
}

func main() {
	// Set up prefix for error codes and hide the timestamp
	log.SetPrefix("Error code: ")
	log.SetFlags(0)

	// Parse the config and setup environment variables
	parseConfig()

	// Setup main variables
	var (
		userNumber    int
		userString    string
		pageNumber    = 1
		results       *tmdb.SearchMulti
		searchDetails bool
		searchSet     bool
		searchAction  string
		searchArray   [2][]string
		err           error
	)

	// Loop everything
	for {
		// Check if the search is set (when there is already searched for something)
		switch searchSet {
		case false:
			fmt.Println("Search for movies/series/actors:")
			// Get the search string
			userString = inputSearch()
			// Get the results from TMDB with the page number and search string
			results = searchTMDB(1, userString)
			// Set the search variable to true
			searchSet = true
			// Print the results in terminal and put in array
			searchArray = printResults(results)

			// Go to the case beneath here
			fallthrough

		case true:
			// Get the input number for navigating
			userNumber, err = inputNavigation()

			if err != nil {
				errorCodes("wrongNumber")
			} else {
				// Check how many pages in search result to check search again and to set searchDetails to true
				searchSet, searchDetails, searchAction = pageLogic(userNumber, results)
			}
		}

		if searchAction != "" && searchAction != "error" {
			switch searchAction {
			case "current":
				fmt.Println("Enter number of the row to choose - 0 to search again:")
				// Get the input number for navigating
				userNumber, err = inputNavigation()

				if err != nil {
					errorCodes("wrongNumber")
				} else {
					searchSet, searchDetails = onePage(userNumber)
				}
			case "next":
				pageNumber++
				// Get the results from TMDB with the page number and search string
				results = searchTMDB(pageNumber, userString)
				searchArray = printResults(results)
			case "previous":
				pageNumber--
				// Get the results from TMDB with the page number and search string
				results = searchTMDB(pageNumber, userString)
				searchArray = printResults(results)
			}

		} else if searchAction == "error" {
			errorCodes("wrongNumber")
		}

		if searchDetails == true {
			fmt.Println("Details of the search:")
			searchChosen(searchArray, userNumber)
			break
		}
	}
}
