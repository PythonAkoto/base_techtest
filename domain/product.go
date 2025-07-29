package domain

type Product struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	Price  float64 `json:"price"`
}

type PricedProduct struct {
	Name string `json:"name"`
	// convert to strings in order to keep trailing zeros (.00) in JSON response
	ProductPrice    string `json:"product_price"`
	DeliveryPrice   string `json:"delivery_price"`
	TotalPrice      string `json:"total_price"`
	DeliveryService string `json:"delivery_service"`
}
