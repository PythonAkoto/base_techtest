package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
	"github.com/PythonAkoto/base_techtest/adapters/output/storage"
	"github.com/PythonAkoto/base_techtest/domain"
)

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// get the default provider from environment variable
	defaultProvider := strings.ToUpper(os.Getenv("DELIVERY_PROVIDER"))
	if defaultProvider == "" {
		logs.Logs(3, "DELIVERY_PROVIDER environment variable not set", "")
		http.Error(w, "Delivery provider not set", http.StatusInternalServerError)
		return
	}

	// check for provider query parameter
	queryProvider := strings.ToUpper(r.URL.Query().Get("provider"))
	var provider string

	if queryProvider == "" {
		// use default from env
		provider = defaultProvider
		logs.Logs(1, "No provider specified in URL, using default from environment", provider)
	} else {
		provider = queryProvider
		if provider != defaultProvider {
			logs.Logs(2, "Query provider differs from env provider", provider)
		}
	}

	// load products from storage
	products, err := storage.LoadProductsFunc()
	if err != nil {
		logs.Logs(3, "Failed to load products: "+err.Error(), provider)
		http.Error(w, "Failed to load products", http.StatusInternalServerError)
		return
	}

	// calculate prices for products
	productPrices, err := domain.PriceProductsFunc(products, provider)
	if err != nil {
		logs.Logs(3, "Failed to price products: "+err.Error(), provider)
		http.Error(w, "Failed to price products", http.StatusInternalServerError)
		return
	}

	// set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// encode products to JSON and write to response
	err = json.NewEncoder(w).Encode(productPrices)
	if err != nil {
		logs.Logs(3, "Failed to write response: "+err.Error(), provider)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// log success message
	logs.Logs(1, "successfully got the prices of the products", provider)
}
