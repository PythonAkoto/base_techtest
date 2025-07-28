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
- [ ] Adequate unit testing
- [ ] API and Mongo running in a container with Docker Compose with Mongo starting before the app on port `27017` with the port coming from a `.env` file
- [x] Appropriate logging/tracing, including a log which states which delivery service was used just like the `Examples` section
- [ ] Postman collection and environment provided
- [ ] Makefile with commands to `stop` and `start` the service

