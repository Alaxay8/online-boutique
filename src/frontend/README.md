# Frontend Service

The Frontend service exposes an HTTP server written in Go that serves the Online Boutique e-commerce website. It communicates with all backend microservices over gRPC and handles session management, product rendering, cart actions, checkout, currency conversion, and theme switching.

## Features

* **Session Management:** Automatically creates and manages anonymous user sessions via cookies.
* **Currency Conversion:** Allows switching currency via `POST /setCurrency`. Supported currencies: `USD`, `EUR`, `JPY`, `GBP`, `TRY`, `RUB`. Persisted in `shop_currency` cookie.
* **Theme Switching:** Supports UI themes (`auto`, `light`, `dark`) selectable via a dropdown in the header or `POST /setTheme`. The setting is stored in the `shop_theme` cookie. When set to `auto`, the interface responds to the browser's `prefers-color-scheme` media query.
* **Assistant / Chatbot:** Optional AI assistant endpoint (`/bot`, `/assistant`) when enabled.

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
export PORT="8080"
export PRODUCT_CATALOG_SERVICE_ADDR="localhost:3550"
export CURRENCY_SERVICE_ADDR="localhost:7000"
export CART_SERVICE_ADDR="localhost:7070"
export RECOMMENDATION_SERVICE_ADDR="localhost:8080"
export SHIPPING_SERVICE_ADDR="localhost:50051"
export CHECKOUT_SERVICE_ADDR="localhost:5050"
export AD_SERVICE_ADDR="localhost:9555"

go run .
```

## Environment Variables

| Variable | Description | Required | Default |
|---|---|---|---|
| `PORT` | Listening port for the HTTP server | No | `8080` |
| `PRODUCT_CATALOG_SERVICE_ADDR` | Host and port of the product catalog service | Yes | - |
| `CURRENCY_SERVICE_ADDR` | Host and port of the currency service | Yes | - |
| `CART_SERVICE_ADDR` | Host and port of the cart service | Yes | - |
| `RECOMMENDATION_SERVICE_ADDR` | Host and port of the recommendation service | Yes | - |
| `SHIPPING_SERVICE_ADDR` | Host and port of the shipping service | Yes | - |
| `CHECKOUT_SERVICE_ADDR` | Host and port of the checkout service | Yes | - |
| `AD_SERVICE_ADDR` | Host and port of the ad service | Yes | - |
| `ENV_PLATFORM` | Deployment platform tag (`local`, `gcp`, `azure`, `aws`, `onprem`, `alibaba`) | No | `local` |
| `ENABLE_PROFILER` | Enable Google Cloud Profiler (`1` to enable) | No | `0` |
| `FRONTEND_MESSAGE` | Custom notification message displayed on banner | No | `""` |
| `CYMBAL_BRANDING` | Enable Cymbal Shops theme branding (`true`/`false`) | No | `false` |
| `ENABLE_ASSISTANT` | Enable shopping assistant AI integration | No | `false` |
| `SHOPPING_ASSISTANT_SERVICE_ADDR` | Address for shopping assistant service | If assistant enabled | - |

## HTTP Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Home page |
| `GET` | `/product/{id}` | Product details page |
| `GET` | `/cart` | View current shopping cart |
| `POST` | `/cart` | Add product to cart |
| `POST` | `/cart/empty` | Empty shopping cart |
| `POST` | `/cart/checkout` | Process checkout and place order |
| `POST` | `/setCurrency` | Change currency preference (form field `currency`) |
| `POST` | `/setTheme` | Change UI theme (form field `theme`: `auto`, `light`, `dark`) |
| `GET` | `/logout` | Clear session and redirect to home |
| `GET` | `/_healthz` | Health check probe endpoint (returns `ok`) |
| `POST` | `/bot` | Chat assistant endpoint |
