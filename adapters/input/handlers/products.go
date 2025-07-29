package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
	"github.com/PythonAkoto/base_techtest/adapters/output/storage"
	"github.com/PythonAkoto/base_techtest/domain"
)

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// get delivery provider from environment variable
	provider := os.Getenv("DELIVERY_PROVIDER")
	if provider == "" {
		logs.Logs(3, "DELIVERY_PROVIDER environment variable not set", provider)
		http.Error(w, "Delivery provider not set", http.StatusInternalServerError)
		return
	}

	// load products from storage
	products, err := storage.LoadProcuts()
	if err != nil {
		logs.Logs(3, "Failed to load products: "+err.Error(), provider)
		http.Error(w, "Failed to load products", http.StatusInternalServerError)
		return
	}

	// calculate prices for products
	productPrices, err := domain.PriceProducts(products)
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
