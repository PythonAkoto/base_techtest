package handlers

import (
	"errors"
	"fmt"
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
// following TDD principles with comprehensive test coverage including query parameter support

func TestGetProductsHandler(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		setupMocks   func()
		expectedCode int
		expectedBody string
	}{
		{
			name:        "delivery provider not set in environment",
			queryParams: "",
			setupMocks: func() {
				os.Unsetenv("DELIVERY_PROVIDER")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Delivery provider not set\n",
		},
		{
			name:        "failed to load products - default provider",
			queryParams: "",
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
			name:        "failed to price products - default provider",
			queryParams: "",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("DHL_DELIVERY_PRICE", "2.00")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Test", Weight: 2, Price: 10},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					return nil, errors.New("pricing failed")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Failed to price products\n",
		},
		{
			name:        "successfully priced products - default provider",
			queryParams: "",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("DHL_DELIVERY_PRICE", "2.00")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item A", Weight: 1.5, Price: 20},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					if provider != "DHL" {
						t.Errorf("Expected provider DHL, got %s", provider)
					}
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
			name:        "successfully priced products - query provider overrides default",
			queryParams: "?provider=UPS",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("UPS_DELIVERY_PRICE", "1.50")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item B", Weight: 2.0, Price: 15},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					if provider != "UPS" {
						t.Errorf("Expected provider UPS, got %s", provider)
					}
					return []domain.PricedProduct{
						{
							Name:            "Item B",
							ProductPrice:    "15.00",
							DeliveryPrice:   "3.00",
							TotalPrice:      "18.00",
							DeliveryService: "UPS",
						},
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"name":"Item B","product_price":"15.00","delivery_price":"3.00","total_price":"18.00","delivery_service":"UPS"}]` + "\n",
		},
		{
			name:        "successfully priced products - multiple items with query provider",
			queryParams: "?provider=AMAZON",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("AMAZON_DELIVERY_PRICE", "1.25")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item A", Weight: 1.0, Price: 15.99},
						{Name: "Item B", Weight: 2.5, Price: 25.50},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					if provider != "AMAZON" {
						t.Errorf("Expected provider AMAZON, got %s", provider)
					}
					return []domain.PricedProduct{
						{
							Name:            "Item A",
							ProductPrice:    "15.99",
							DeliveryPrice:   "1.25",
							TotalPrice:      "17.24",
							DeliveryService: "AMAZON",
						},
						{
							Name:            "Item B",
							ProductPrice:    "25.50",
							DeliveryPrice:   "3.13",
							TotalPrice:      "28.63",
							DeliveryService: "AMAZON",
						},
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"name":"Item A","product_price":"15.99","delivery_price":"1.25","total_price":"17.24","delivery_service":"AMAZON"},{"name":"Item B","product_price":"25.50","delivery_price":"3.13","total_price":"28.63","delivery_service":"AMAZON"}]` + "\n",
		},
		{
			name:        "empty products list with query provider",
			queryParams: "?provider=DPD",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("DPD_DELIVERY_PRICE", "2.50")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					if provider != "DPD" {
						t.Errorf("Expected provider DPD, got %s", provider)
					}
					return []domain.PricedProduct{}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: "[]\n",
		},
		{
			name:        "case insensitive query provider",
			queryParams: "?provider=ups",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("UPS_DELIVERY_PRICE", "1.50")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item C", Weight: 1.0, Price: 10},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					if provider != "UPS" {
						t.Errorf("Expected provider UPS (uppercase), got %s", provider)
					}
					return []domain.PricedProduct{
						{
							Name:            "Item C",
							ProductPrice:    "10.00",
							DeliveryPrice:   "1.50",
							TotalPrice:      "11.50",
							DeliveryService: "UPS",
						},
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"name":"Item C","product_price":"10.00","delivery_price":"1.50","total_price":"11.50","delivery_service":"UPS"}]` + "\n",
		},
		{
			name:        "invalid provider in query parameter",
			queryParams: "?provider=INVALID",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item D", Weight: 1.0, Price: 10},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					if provider != "INVALID" {
						t.Errorf("Expected provider INVALID, got %s", provider)
					}
					return nil, errors.New("delivery provider not set")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Failed to price products\n",
		},
		{
			name:        "test ROYALMAIL provider",
			queryParams: "?provider=ROYALMAIL",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("ROYAL_MAIL_DELIVERY_PRICE", "3.00")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item E", Weight: 1.5, Price: 12.50},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					if provider != "ROYALMAIL" {
						t.Errorf("Expected provider ROYALMAIL, got %s", provider)
					}
					return []domain.PricedProduct{
						{
							Name:            "Item E",
							ProductPrice:    "12.50",
							DeliveryPrice:   "4.50",
							TotalPrice:      "17.00",
							DeliveryService: "ROYALMAIL",
						},
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"name":"Item E","product_price":"12.50","delivery_price":"4.50","total_price":"17.00","delivery_service":"ROYALMAIL"}]` + "\n",
		},
		{
			name:        "test YODEL provider",
			queryParams: "?provider=YODEL",
			setupMocks: func() {
				os.Setenv("DELIVERY_PROVIDER", "DHL")
				os.Setenv("YODEL_DELIVERY_PRICE", "2.75")
				storage.LoadProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Name: "Item F", Weight: 0.5, Price: 8.99},
					}, nil
				}
				domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
					if provider != "YODEL" {
						t.Errorf("Expected provider YODEL, got %s", provider)
					}
					return []domain.PricedProduct{
						{
							Name:            "Item F",
							ProductPrice:    "8.99",
							DeliveryPrice:   "1.38",
							TotalPrice:      "10.37",
							DeliveryService: "YODEL",
						},
					}, nil
				}
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"name":"Item F","product_price":"8.99","delivery_price":"1.38","total_price":"10.37","delivery_service":"YODEL"}]` + "\n",
		},
	}

	for _, tc := range tests {
		runGetProductsTest(t, tc.name, tc.queryParams, tc.setupMocks, tc.expectedCode, tc.expectedBody)
	}
}

func runGetProductsTest(t *testing.T, name string, queryParams string, setupMocks func(), expectedCode int, expectedBody string) {
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

		// Create request with query parameters
		url := "/products" + queryParams
		req := httptest.NewRequest("GET", url, nil)
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
		os.Unsetenv("AMAZON_DELIVERY_PRICE")
		os.Unsetenv("ROYAL_MAIL_DELIVERY_PRICE")
		os.Unsetenv("DPD_DELIVERY_PRICE")
		os.Unsetenv("YODEL_DELIVERY_PRICE")
		os.Unsetenv("PRODUCTS_FILE_PATH")
	})
}

// BenchmarkGetProductsHandler benchmarks the GetProductsHandler function with different scenarios
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
	
	domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
		return []domain.PricedProduct{
			{Name: "Item A", ProductPrice: "20.00", DeliveryPrice: "3.00", TotalPrice: "23.00", DeliveryService: provider},
			{Name: "Item B", ProductPrice: "15.00", DeliveryPrice: "4.00", TotalPrice: "19.00", DeliveryService: provider},
		}, nil
	}

	// Benchmark default provider
	b.Run("DefaultProvider", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("GET", "/products", nil)
			w := httptest.NewRecorder()
			GetProductsHandler(w, req)
		}
	})

	// Benchmark with query provider
	b.Run("QueryProvider", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("GET", "/products?provider=UPS", nil)
			w := httptest.NewRecorder()
			GetProductsHandler(w, req)
		}
	})

	// Cleanup
	os.Unsetenv("DELIVERY_PROVIDER")
	os.Unsetenv("DHL_DELIVERY_PRICE")
}

// BenchmarkGetProductsHandlerWithLargeDataset benchmarks with a larger dataset
func BenchmarkGetProductsHandlerWithLargeDataset(b *testing.B) {
	// Setup with larger dataset
	os.Setenv("DELIVERY_PROVIDER", "DHL")
	os.Setenv("DHL_DELIVERY_PRICE", "2.00")
	
	// Create a larger dataset
	largeProductSet := make([]domain.Product, 100)
	largePricedSet := make([]domain.PricedProduct, 100)
	
	for i := 0; i < 100; i++ {
		largeProductSet[i] = domain.Product{
			Name:   fmt.Sprintf("Item_%d", i),
			Weight: float64(i%5 + 1),
			Price:  float64(10 + i%20),
		}
		largePricedSet[i] = domain.PricedProduct{
			Name:            fmt.Sprintf("Item_%d", i),
			ProductPrice:    fmt.Sprintf("%.2f", float64(10+i%20)),
			DeliveryPrice:   fmt.Sprintf("%.2f", float64((i%5+1)*2)),
			TotalPrice:      fmt.Sprintf("%.2f", float64(10+i%20+(i%5+1)*2)),
			DeliveryService: "DHL",
		}
	}
	
	storage.LoadProductsFunc = func() ([]domain.Product, error) {
		return largeProductSet, nil
	}
	
	domain.PriceProductsFunc = func(products []domain.Product, provider string) ([]domain.PricedProduct, error) {
		return largePricedSet, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()
		GetProductsHandler(w, req)
	}

	// Cleanup
	os.Unsetenv("DELIVERY_PROVIDER")
	os.Unsetenv("DHL_DELIVERY_PRICE")
}
