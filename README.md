# Base Media Cloud Test

## Introduction
This test/training piece will demonstrate the necessary knowledge to work at Base.

## User Story
We have an electronics store which sells products online. We have an array of products in the `products.json` file. Each product contains a `name`, `weight` and `price`. 

We frequently change between delivery companies for pricing reasons. To ensure accurate pricing, we need you to create a new endpoint which will return the up-to-date prices. This will enable us to display the correct price on our website and send orders to the appropriate delivery company. We would prefer to manage pricing and the selection of delivery companies by environment variables.

We need to be able to make a request to a `/products` endpoint, which will return an array of products with their prices. Each product should include the `name`, `product_price`, `delivery_price`, and `total_price`.

We would like you to deploy the API in a container with Docker Compose. We expect appropriate logging/tracing throughout so that we can monitor errors and see what's happening in the container. In the future, we will be introducing MongoDB to store our products in a database rather than a JSON file. We want you to spin up a Mongo container which serves on port `27017` and gets this value from a `.env` file. We want to ensure Mongo starts up before the API, but we don't need it to do anything else.

We need a Makefile with commands to `start` and `stop` the service. We also need a Postman collection/environment which contains the request. Our team believes strongly in writing maintainable and testable code. We expect the app to be adequately unit-tested.

Feel free to import any packages which will make your life easier. Be sure to keep a list, as you will be asked to justify your reasoning for each.

## Examples

### Response
The example below is using UPS as the delivery service and has the price set as 1p.
`GET` - `localhost:8000/products`.

```
[
    {
        "name": "Phone",
        "delivery_price": "2.21",
        "product_price": "1000.00",
        "total_price": "1002.21"
    },
    {
        "name": "TV",
        "delivery_price": "100.00",
        "product_price": "800.00",
        "total_price": "900.00"
    },
    {
        "name": "Laptop",
        "delivery_price": "14.00",
        "product_price": "5000.00",
        "total_price": "5014.00"
    },
    {
        "name": "Speaker",
        "delivery_price": "6.80",
        "product_price": "200.00",
        "total_price": "206.80"
    },
    {
        "name": "Keyboard",
        "delivery_price": "7.50",
        "product_price": "150.00",
        "total_price": "157.50"
    },
    {
        "name": "Mouse",
        "delivery_price": "1.50",
        "product_price": "100.00",
        "total_price": "101.50"
    },
    {
        "name": "Headphones",
        "delivery_price": "3.75",
        "product_price": "350.00",
        "total_price": "353.75"
    },
    {
        "name": "Microphone",
        "delivery_price": "0.05",
        "product_price": "120.00",
        "total_price": "120.05"
    },
    {
        "name": "Tablet",
        "delivery_price": "6.50",
        "product_price": "900.00",
        "total_price": "906.50"
    },
    {
        "name": "Webcam",
        "delivery_price": "0.63",
        "product_price": "120.00",
        "total_price": "120.63"
    }
]
```

### Success Log
```
{"level":"info","provider":"ups","time":1695987270,"message":"successfully got the prices of the products"}
```

## Out Of scope
- Front end
- Graceful shutdowns
- Saving the records to the DB
- End-To-End tests
- Go Docs
- Allowing multiple delivery companies at once
- Anything related to sending orders

## Acceptance Criteria 
- [x] Create an endpoint `GET` `/products` which returns the same response as the `Examples` section
- [x] The prices must be calculated at request time
- [x] The project must follow Hexagonal Architecture
- [x] The numbers in the response, must all be to 2 decimal places
- [x] Able to switch between delivery companies by an environment variable and it logs which delivery service was used. TIP: we expect an implementation for each delivery service
- [x] Able to set which port the app runs on via an environment variable
- [x] The price can be changed from an environment variable
- [x] Adequate unit testing
- [ ] API and Mongo running in a container with Docker Compose with Mongo starting before the app on port `27017` with the port coming from a `.env` file
- [x] Appropriate logging/tracing, including a log which states which delivery service was used just like the `Examples` section
- [ ] Postman collection and environment provided
- [ ] Makefile with commands to `stop` and `start` the service

---

## Installation Guide

This step-by-step guide will help you recreate and understand the project architecture, highlighting the advantages of each implementation decision.

### Prerequisites
- Go 1.19 or higher
- Git
- Text editor/IDE

### Step 1: Project Structure Setup
```bash
mkdir base_techtest
cd base_techtest
go mod init github.com/PythonAkoto/base_techtest
```

**Advantages of this structure:**
- Clean Go module initialization
- Proper package naming convention
- Scalable foundation for hexagonal architecture

### Step 2: Create Hexagonal Architecture Directories
```bash
mkdir -p domain
mkdir -p adapters/input/handlers
mkdir -p adapters/input/static
mkdir -p adapters/output/storage
mkdir -p adapters/output/logs
mkdir -p adapters/output/delivery
mkdir -p env
```

**Advantages of hexagonal architecture:**
- **Separation of Concerns**: Business logic isolated in domain layer
- **Testability**: Each layer can be tested independently
- **Flexibility**: Easy to swap implementations (e.g., JSON file to database)
- **Maintainability**: Clear boundaries between layers reduce coupling

### Step 3: Domain Layer Implementation
Create the core business entities and logic:

**File: `domain/product.go`**
- Defines `Product` and `PricedProduct` structs
- Uses JSON tags with snake_case for API consistency
- String formatting for prices to maintain trailing zeros

**File: `domain/pricing.go`**
- Contains `PriceProducts` function with provider parameter support
- Implements delivery price calculation logic
- Uses function variables (`PriceProductsFunc`) for testability

**Advantages of domain-first approach:**
- **Business Logic Isolation**: No external dependencies in core logic
- **Type Safety**: Strong typing prevents runtime errors
- **Mockable Design**: Function variables enable comprehensive testing
- **Provider Flexibility**: Dynamic provider selection via parameters

### Step 4: Output Adapters (Infrastructure Layer)
**File: `adapters/output/storage/loadProducts.go`**
- Handles JSON file loading with environment variable configuration
- Implements error handling for missing files/environment variables
- Uses function variables for testing isolation

**File: `adapters/output/logs/logs.go`**
- Centralized logging system with different log levels
- Provider-aware logging for delivery service tracking
- Asynchronous log processing for performance

**Advantages of output adapters:**
- **External Dependency Isolation**: File I/O and logging separated from business logic
- **Configuration Flexibility**: Environment variable driven configuration
- **Error Resilience**: Proper error handling and logging
- **Performance**: Asynchronous logging doesn't block main execution

### Step 5: Input Adapters (Presentation Layer)
**File: `adapters/input/handlers/products.go`**
- HTTP handler with query parameter support (`?provider=UPS`)
- Environment variable fallback mechanism
- Proper HTTP status codes and JSON responses
- Case-insensitive provider handling

**File: `adapters/input/handlers/server.go`**
- HTTP server setup with environment-driven port configuration
- Route registration and static file serving capability
- Graceful error handling and logging integration

**Advantages of input adapters:**
- **Protocol Independence**: Business logic doesn't know about HTTP
- **Query Parameter Flexibility**: Runtime provider override capability
- **Fallback Mechanisms**: Environment variable defaults with query override
- **HTTP Best Practices**: Proper status codes and content types

### Step 6: Environment Configuration
**⚠️ IMPORTANT: Create `.env` file**

You must create a `.env` file in the `env/` directory with the following contents (this file is in `.gitignore` and won't be included in the repository):

```env
# Application Configuration
APP_PORT=9000
PRODUCTS_FILE_PATH=adapters/output/storage/products.json

# Delivery Provider Pricing (price per unit weight)
DHL_DELIVERY_PRICE=0.15
UPS_DELIVERY_PRICE=0.12
AMAZON_DELIVERY_PRICE=0.18
ROYAL_MAIL_DELIVERY_PRICE=0.10
DPD_DELIVERY_PRICE=0.14
YODEL_DELIVERY_PRICE=0.16

# Default Delivery Provider
DELIVERY_PROVIDER=UPS
```

**Advantages of environment configuration:**
- **Security**: Sensitive configuration kept out of source control
- **Flexibility**: Easy deployment across different environments
- **Runtime Configuration**: No code changes needed for different setups
- **Twelve-Factor App Compliance**: Configuration through environment variables

### Step 7: Products Data Setup
**File: `adapters/output/storage/products.json`**
- Contains sample product data with name, weight, and price
- JSON format for easy parsing and modification
- Realistic product data for testing

### Step 8: Main Application Entry Point
**File: `main.go`**
- Goroutine-based concurrent execution
- Environment loading and HTTP server startup
- Log processing initialization
- Blocking main goroutine to keep application alive

**Advantages of concurrent design:**
- **Performance**: Non-blocking log processing
- **Responsiveness**: HTTP server runs independently
- **Resource Efficiency**: Goroutines are lightweight
- **Scalability**: Ready for additional concurrent services

### Step 9: Comprehensive Testing Implementation
**File: `adapters/input/handlers/products_test.go`**
- 12 comprehensive test cases covering all scenarios
- Mock implementations for external dependencies
- Environment variable testing and cleanup
- Query parameter validation testing
- Error condition testing

**Advantages of comprehensive testing:**
- **TDD Compliance**: Test-driven development approach
- **High Coverage**: 54.0% code coverage achieved
- **Isolation**: Proper mocking prevents external dependencies
- **Regression Prevention**: Comprehensive scenarios prevent future bugs
- **Documentation**: Tests serve as usage examples

### Step 10: Build and Run
```bash
# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run the application
go run main.go
```

### Step 11: Test the API
```bash
# Test with default provider (from environment)
curl http://localhost:9000/products

# Test with query parameter override
curl "http://localhost:9000/products?provider=DHL"

# Test case insensitive
curl "http://localhost:9000/products?provider=ups"
```

---

## Design Patterns Used

### Factory Pattern Implementation
The project implements the **Factory Pattern** in the delivery price calculation system:

**Location**: `domain/pricing.go` - `calculateDeliveryPrice()` function

**Implementation Details:**
- Switch statement acts as a factory for delivery price calculations
- Each case creates a specific pricing strategy based on provider
- Environment variable configuration drives factory decisions
- Default case provides fallback behavior

**Benefits of Factory Pattern:**
- **Encapsulation**: Object creation logic is centralized
- **Flexibility**: Easy to add new delivery providers
- **Consistency**: Uniform interface for all delivery calculations
- **Maintainability**: Single location for delivery provider logic
- **Configuration-Driven**: Runtime provider selection through environment variables

**Factory Pattern Advantages in This Context:**
1. **Provider Abstraction**: Client code doesn't need to know specific provider implementations
2. **Easy Extension**: Adding new providers requires only adding new cases
3. **Runtime Selection**: Provider choice determined at runtime via environment/query parameters
4. **Consistent Interface**: All providers return the same data structure
5. **Error Handling**: Centralized error handling for all provider types

---

## Future Refactoring Opportunities

If I were to refactor this project for improved efficiency and maintainability, here are the key areas I would address:

### 1. Interface-Based Design
**Current**: Function variables for mocking
**Improvement**: Explicit interfaces for better abstraction
```go
type ProductRepository interface {
    LoadProducts() ([]Product, error)
}

type PricingService interface {
    PriceProducts(products []Product, provider string) ([]PricedProduct, error)
}
```
**Benefits**: Better abstraction, clearer contracts, improved testability

### 2. Configuration Injection
**Current**: Environment variables accessed directly in domain layer
**Improvement**: Configuration struct injection
```go
type DeliveryConfig struct {
    Providers map[string]float64
    Default   string
}
```
**Benefits**: Reduced coupling, easier testing, better configuration management

### 3. Provider Strategy Pattern
**Current**: Switch statement for provider selection
**Improvement**: Strategy pattern with provider interfaces
```go
type DeliveryProvider interface {
    CalculatePrice(weight float64) (float64, error)
    Name() string
}
```
**Benefits**: Better extensibility, cleaner code, easier testing of individual providers

### 4. Caching Layer
**Current**: File read on every request
**Improvement**: In-memory caching with file watching
**Benefits**: Improved performance, reduced I/O operations, better scalability

### 5. Structured Logging
**Current**: Custom logging system
**Improvement**: Structured logging with libraries like logrus or zap
**Benefits**: Better log parsing, structured data, improved observability

### 6. Error Handling Enhancement
**Current**: Basic error returns
**Improvement**: Custom error types with context
```go
type DeliveryError struct {
    Provider string
    Reason   string
    Err      error
}
```
**Benefits**: Better error context, improved debugging, cleaner error handling

### 7. Validation Layer
**Current**: Basic provider validation
**Improvement**: Comprehensive input validation with validator library
**Benefits**: Better data integrity, clearer error messages, reduced runtime errors

### 8. Database Abstraction
**Current**: Direct file I/O
**Improvement**: Repository pattern with database interface
**Benefits**: Database agnostic, easier testing, better data management

### 9. HTTP Middleware
**Current**: Basic HTTP handlers
**Improvement**: Middleware for logging, authentication, rate limiting
**Benefits**: Cross-cutting concerns separation, better request handling, improved security

### 10. Configuration Management
**Current**: Environment variables only
**Improvement**: Hierarchical configuration (env vars, config files, defaults)
**Benefits**: Better configuration management, environment-specific settings, easier deployment

These refactoring opportunities would transform the project into a more enterprise-ready application while maintaining the solid architectural foundation already established.

