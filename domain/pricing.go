package domain

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
)

/*
PriceProducts calculates the delivery price and total price for a list of products
based on their weight and the delivery provider specified in the environment.
It returns a slice of PricedProduct containing the pricing details for each product,
or an error if the delivery provider is not set in the environment variables.
*/
func PriceProducts(products []Product) ([]PricedProduct, error) {
	provider := os.Getenv("DELIVERY_PROVIDER")
	// Check if the delivery provider is set in the environment variables
	allowedProviders := []string{"DHL", "UPS", "AMAZON", "ROYAL_MAIL", "DPD", "YODEL"}
	if !contains(allowedProviders, provider) {
		// Log an error if the delivery provider is not set
		logs.Logs(3, "DELIVERY_PROVIDER environment variable not set", provider)
		return nil, fmt.Errorf("delivery provider not set")
	}

	var result []PricedProduct

	for _, product := range products {
		deliveryPrice, err := calculateDeliveryPrice(product.Weight, provider)
		if err != nil {
			logs.Logs(3, fmt.Sprintf("failed to calculate delivery price for product %s: %s", product.Name, err.Error()), provider)
			return nil, err
		}

		// convert and calculate prices
		productPrincing := roundToTwoDecimalPlaces(product.Price)
		deliveryPricing := roundToTwoDecimalPlaces(deliveryPrice)
		total := productPrincing + deliveryPricing

		finalPrice := PricedProduct{
			Name:            product.Name,
			ProductPrice:    fmt.Sprintf("%.2f", productPrincing),
			DeliveryPrice:   fmt.Sprintf("%.2f", deliveryPricing),
			TotalPrice:      fmt.Sprintf("%.2f", total),
			DeliveryService: provider,
		}
		result = append(result, finalPrice)
		logs.Logs(1, "Product priced successfully: "+product.Name+" with total price: "+fmt.Sprintf("%.2f", total), provider)
	}

	return result, nil
}

/*
Contains checks if a string provider is present in a slice of strings.
It iterates through the slice and returns true if the provider is found,
otherwise it returns false.
*/
func contains(slice []string, provider string) bool {
	for _, item := range slice {
		if item == provider {
			return true
		}
	}
	return false
}

/*
roundToTwoDecimalPlaces rounds a given float64 value to two decimal places.

It multiplies the value by 100, converts it to an integer to remove any fractional part beyond two decimal places, and then divides by 100.0 to return the rounded result.

Parameters:
- value: The float64 number to be rounded.

Returns:
- A float64 number rounded to two decimal places.
*/
func roundToTwoDecimalPlaces(value float64) float64 {
	// Round the value to two decimal places
	return float64(int(value*100)) / 100.0
}

/*
CalculateDeliveryPrice returns the delivery price based on the weight of the product and the provider.
If the provider does not match any of the cases, it will default to 0.20 * weight.
If the environment variable for the provider is not set, it will return an error.
If the environment variable for the provider is set, but has an invalid value, it will return an error.
*/
func calculateDeliveryPrice(weight float64, provider string) (float64, error) {
	switch provider {
	case "DHL":
		DHL_DELIVERY_PRICE, err := strconv.ParseFloat(os.Getenv("DHL_DELIVERY_PRICE"), 64)
		if err != nil {
			if err == strconv.ErrSyntax {
				return 0, fmt.Errorf("DHL_DELIVERY_PRICE environment variable not set")
			}
			return 0, fmt.Errorf("invalid DHL_DELIVERY_PRICE: %s", err.Error())
		}
		return weight * DHL_DELIVERY_PRICE, nil

	case "UPS":
		UPS_DELIVERY_PRICE, err := strconv.ParseFloat(os.Getenv("UPS_DELIVERY_PRICE"), 64)
		if err != nil {
			if err == strconv.ErrSyntax {
				return 0, fmt.Errorf("UPS_DELIVERY_PRICE environment variable not set")
			}
			return 0, fmt.Errorf("invalid UPS_DELIVERY_PRICE: %s", err.Error())
		}
		return weight * UPS_DELIVERY_PRICE, nil

	case "AMAZON":
		AMAZON_DELIVERY_PRICE, err := strconv.ParseFloat(os.Getenv("AMAZON_DELIVERY_PRICE"), 64)
		if err != nil {
			if err == strconv.ErrSyntax {
				return 0, fmt.Errorf("AMAZON_DELIVERY_PRICE environment variable not set")
			}
			return 0, fmt.Errorf("invalid AMAZON_DELIVERY_PRICE: %s", err.Error())
		}
		return weight * AMAZON_DELIVERY_PRICE, nil

	case "ROYAL_MAIL":
		ROYAL_MAIL_DELIVERY_PRICE, err := strconv.ParseFloat(os.Getenv("ROYAL_MAIL_DELIVERY_PRICE"), 64)
		if err != nil {
			if err == strconv.ErrSyntax {
				return 0, fmt.Errorf("ROYAL_MAIL_DELIVERY_PRICE environment variable not set")
			}
			return 0, fmt.Errorf("invalid ROYAL_MAIL_DELIVERY_PRICE: %s", err.Error())
		}
		return weight * ROYAL_MAIL_DELIVERY_PRICE, nil

	case "DPD":
		DPD_DELIVERY_PRICE, err := strconv.ParseFloat(os.Getenv("DPD_DELIVERY_PRICE"), 64)
		if err != nil {
			if err == strconv.ErrSyntax {
				return 0, fmt.Errorf("DPD_DELIVERY_PRICE environment variable not set")
			}
			return 0, fmt.Errorf("invalid DPD_DELIVERY_PRICE: %s", err.Error())
		}
		return weight * DPD_DELIVERY_PRICE, nil

	case "YODEL":
		YODEL_DELIVERY_PRICE, err := strconv.ParseFloat(os.Getenv("YODEL_DELIVERY_PRICE"), 64)
		if err != nil {
			if err == strconv.ErrSyntax {
				return 0, fmt.Errorf("YODEL_DELIVERY_PRICE environment variable not set")
			}
			return 0, fmt.Errorf("invalid YODEL_DELIVERY_PRICE: %s", err.Error())
		}
		return weight * YODEL_DELIVERY_PRICE, nil

	default:
		return weight * 0.20, nil // Default case if no provider matches
	}
}
