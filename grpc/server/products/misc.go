package products

import (
	"encoding/csv"
	"net/http"
	"os"
)

func readCSVFromURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'

	return reader.ReadAll()
}

func readCSVFromFile(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.Comma = ';'

	return reader.ReadAll()
}
