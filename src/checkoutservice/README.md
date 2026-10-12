# Checkout Service

The Checkout Service coordinates the order checkout process: it retrieves the user's cart, calculates items and shipping costs, processes payment, calls the shipping service, sends an order confirmation email, empties the cart, and asynchronously publishes order placed events to Apache Kafka.

## Dependency Management & Local Development

Dependencies are managed using Go modules (Go 1.23+).

To download and vendor dependencies:

```sh
go mod download
go mod vendor
```

To run tests:

```sh
go test ./...
```

To run locally:

```sh
export PORT="5050"
export PRODUCT_CATALOG_SERVICE_ADDR="localhost:3550"
export SHIPPING_SERVICE_ADDR="localhost:50051"
export PAYMENT_SERVICE_ADDR="localhost:50051"
export EMAIL_SERVICE_ADDR="localhost:5000"
export CURRENCY_SERVICE_ADDR="localhost:7000"
export CART_SERVICE_ADDR="localhost:7070"
export KAFKA_REST_URL="http://localhost:8082"
export KAFKA_ORDERS_TOPIC="orders"

go run .
```

## Environment Variables

| Variable | Description | Required | Default |
|---|---|---|---|
| `PORT` | Port for the gRPC service | No | `5050` |
| `PRODUCT_CATALOG_SERVICE_ADDR` | Host and port of the product catalog service | Yes | - |
| `SHIPPING_SERVICE_ADDR` | Host and port of the shipping service | Yes | - |
| `PAYMENT_SERVICE_ADDR` | Host and port of the payment service | Yes | - |
| `EMAIL_SERVICE_ADDR` | Host and port of the email service | Yes | - |
| `CURRENCY_SERVICE_ADDR` | Host and port of the currency service | Yes | - |
| `CART_SERVICE_ADDR` | Host and port of the cart service | Yes | - |
| `KAFKA_REST_URL` | Base URL of the Kafka REST Proxy (e.g. `http://kafka-rest:8082`) | No | `""` (skips publish if unset) |
| `KAFKA_ORDERS_TOPIC` | Kafka topic name for order events (e.g. `orders`) | No | `""` (skips publish if unset) |

## Kafka Integration

When `PlaceOrder` completes successfully, an order placed event is published to Kafka via the Kafka REST Proxy:

* **Endpoint:** `POST ${KAFKA_REST_URL}/topics/${KAFKA_ORDERS_TOPIC}`
* **Header:** `Content-Type: application/vnd.kafka.json.v2+json`
* **Message Key:** `orderId`
* **Payload Structure (`orderPlacedEvent`):**
  ```json
  {
    "orderId": "uuid-string",
    "userId": "string",
    "userCurrency": "USD",
    "total": {
      "currencyCode": "USD",
      "units": 100,
      "nanos": 0
    },
    "items": [
      {
        "item": { "productId": "...", "quantity": 1 },
        "cost": { "currencyCode": "USD", "units": 100, "nanos": 0 }
      }
    ],
    "shippingTrackingId": "uuid-string",
    "shippingAddress": {
      "streetAddress": "...",
      "city": "...",
      "state": "...",
      "country": "...",
      "zipCode": 12345
    }
  }
  ```

If `KAFKA_REST_URL` or `KAFKA_ORDERS_TOPIC` is not set or if publishing fails, an error is logged as a warning, but the order placement call itself still succeeds so as not to interrupt user purchases.
