package storage

import (
	"encoding/json"
	"os"

	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
	"github.com/PythonAkoto/base_techtest/domain"
)

func LoadProcuts() ([]domain.Product, error) {
	path := os.Getenv("PRODUCTS_FILE_PATH") // Get the path to the products file from environment variable
	if path == "" {
		logs.Logs(3, "PRODUCTS_FILE_PATH environment variable not set", "")
		return nil, os.ErrNotExist
	}

	file, err := os.Open(path) // Open the products file
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure the file is closed after reading

	var products []domain.Product
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&products) // Decode the JSON data into the products slice
	if err != nil {
		return nil, err
	}

	return products, nil // Return the loaded products
}
