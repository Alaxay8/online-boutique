# Shipping Service

The Shipping service provides price quotes, tracking IDs, and simulated order fulfillment & shipping processes.

## Local Development & Dependencies

Dependencies are managed using Go modules (Go 1.23+).

Run the following command to download and vendor dependencies:

```sh
go mod download
go mod vendor
```

## Build

From `src/shippingservice`, run:

```sh
docker build -t shippingservice .
```

## Test

```sh
go test ./...
```
