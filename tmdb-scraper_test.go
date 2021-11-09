package main

import (
	"testing"
)

func TestOnePage(t *testing.T) {
	parseConfig()
	results := searchTMDB(1, "Hallo K3")
	searchArray := printResults(results)
	searchChosen(searchArray, 1)
	if results.TotalResults != 8 {
		t.Errorf("Results number was incorrect, got: %d, want: %d.", results.TotalResults, 8)
	}
}
