package products

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromCSV(t *testing.T) {
	result := [][]string{
		{"\ufeffproduct1", "34.6"},
		{"product2", "34.7"},
		{"product3", "34.8"},
		{"product4", "34.9"},
		{"product5", "34.10"},
		{"product6", "34.11"},
	}
	data, err := readCSVFromFile("./product.csv")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, data, result)
}
