package domain

type Product struct {
	Name string `json:"name"`
	Weight float64    `json:"weight"`
	Price float64    `json:"price"`
}

type PricedProduct struct {
	Name string  `json:"name"`
	ProductPrice float64 `json:"product_price"`
	DeliveryPrice float64 `json:"delivery_price"`
	TotalPrice float64 `json:"total_price"`
	DeliveryService string `json:"delivery_service"`
}