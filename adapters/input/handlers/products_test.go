package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
	"github.com/PythonAkoto/base_techtest/adapters/output/storage"
	"github.com/PythonAkoto/base_techtest/domain"
)

// TestGetProductsHandler tests the GetProductsHandler function with various scenarios
// following TDD principles with comprehensive test coverage

func TestGetProductsHandler(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "delivery provider not set",
			setupMocks: func() {
				os.Unsetenv("DELIVERY_PROVIDER")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Delivery provider not set\n",
		},
		{
			name: "failed to load products",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("DHL_DELIVERY_PRICE", "2.00")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return nil, errors.New("failed to load")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Failed to load products\n",
		},
		{
			name: "failed to price products",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("DHL_DELIVERY_PRICE", "2.00")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Test", Weight: 2, Price: 10},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product) ([]domain.PricedProduct, error) {
					return nil, errors.New("pricing failed")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Failed to price products\n",
		},
		{
			name: "successfully priced products - single item",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("DHL_DELIVERY_PRICE", "2.00")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item A", Weight: 1.5, Price: 20},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product) ([]domain.PricedProduct, error) {
					return []domain.PricedProduct{
						{
							Name:            "Item A",
							ProductPrice:    "20.00",
							DeliveryPrice:   "3.00",
							TotalPrice:      "23.00",
							DeliveryService: "DHL",
						},
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"name":"Item A","product_price":"20.00","delivery_price":"3.00","total_price":"23.00","delivery_service":"DHL"}]` + "\n",
		},
		{
			name: "successfully priced products - multiple items",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "UPS")
				os.Setenv("UPS_DELIVERY_PRICE", "1.50")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item A", Weight: 1.0, Price: 15.99},
						{Name: "Item B", Weight: 2.5, Price: 25.50},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product) ([]domain.PricedProduct, error) {
					return []domain.PricedProduct{
						{
							Name:            "Item A",
							ProductPrice:    "15.99",
							DeliveryPrice:   "1.50",
							TotalPrice:      "17.49",
							DeliveryService: "UPS",
						},
						{
							Name:            "Item B",
							ProductPrice:    "25.50",
							DeliveryPrice:   "3.75",
							TotalPrice:      "29.25",
							DeliveryService: "UPS",
						},
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"name":"Item A","product_price":"15.99","delivery_price":"1.50","total_price":"17.49","delivery_service":"UPS"},{"name":"Item B","product_price":"25.50","delivery_price":"3.75","total_price":"29.25","delivery_service":"UPS"}]` + "\n",
		},
		{
			name: "empty products list",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("DHL_DELIVERY_PRICE", "2.00")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product) ([]domain.PricedProduct, error) {
					return []domain.PricedProduct{}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: "[]\n",
		},
	}

	for _, tc := range tests {
		runGetProductsTest(t, tc.name, tc.setupMocks, tc.expectedCode, tc.expectedBody)
	}
}

func runGetProductsTest(t *testing.T, name string, setupMocks func(), expectedCode int, expectedBody string) {
	t.Run(name, func(t *testing.T) {
		// Store original functions for restoration
		originalLoadProductsFunc := storage.LoadProductsFunc
		originalPriceProductsFunc := domain.PriceProductsFunc

		// Reset to real implementations before the test
		storage.LoadProductsFunc = storage.LoadProducts
		domain.PriceProductsFunc = domain.PriceProducts

		log.SetFlags(0) // Disable log timestamps for cleaner test output

		// Start the log processor if needed
		go logs.ProcessLogs()

		// Setup mocks for the test
		setupMocks()

		// Create request and response recorder
		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()

		// Execute the handler
		GetProductsHandler(w, req)

		// Assertions
		if w.Code != expectedCode {
			t.Errorf("[%s] expected code %d, got %d", name, expectedCode, w.Code)
		}
		if w.Body.String() != expectedBody {
			t.Errorf("[%s] expected body %q, got %q", name, expectedBody, w.Body.String())
		}

		// Cleanup - restore original functions and clear environment variables
		storage.LoadProductsFunc = originalLoadProductsFunc
		domain.PriceProductsFunc = originalPriceProductsFunc
		os.Unsetenv("DELIVERY_PROVIDER")
		os.Unsetenv("DHL_DELIVERY_PRICE")
		os.Unsetenv("UPS_DELIVERY_PRICE")
		os.Unsetenv("PRODUCTS_FILE_PATH")
	})
}

// BenchmarkGetProductsHandler benchmarks the GetProductsHandler function
func BenchmarkGetProductsHandler(b *testing.B) {
	// Setup
	os.Setenv("DELIVERY_PROVIDER", "DHL")
	os.Setenv("DHL_DELIVERY_PRICE", "2.00")
	
	storage.LoadProductsFunc = func() ([]domain.Product, error) {
		return []domain.Product{
			{Name: "Item A", Weight: 1.5, Price: 20},
			{Name: "Item B", Weight: 2.0, Price: 15},
		}, nil
	}
	
	domain.PriceProductsFunc = func(products []domain.Product) ([]domain.PricedProduct, error) {
		return []domain.PricedProduct{
			{Name: "Item A", ProductPrice: "20.00", DeliveryPrice: "3.00", TotalPrice: "23.00", DeliveryService: "DHL"},
			{Name: "Item B", ProductPrice: "15.00", DeliveryPrice: "4.00", TotalPrice: "19.00", DeliveryService: "DHL"},
		}, nil
	}

	// Benchmark
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()
		GetProductsHandler(w, req)
	}

	// Cleanup
	os.Unsetenv("DELIVERY_PROVIDER")
	os.Unsetenv("DHL_DELIVERY_PRICE")
}
