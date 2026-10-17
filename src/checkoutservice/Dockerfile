# Builder stage
FROM golang:1.23.4-alpine AS builder

ARG TARGETARCH
ARG TARGETOS
ARG SKAFFOLD_GO_GCFLAGS

WORKDIR /src

# Restore dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /checkoutservice .

# Final stage
FROM scratch

WORKDIR /src
COPY --from=builder /checkoutservice /src/checkoutservice

ENV GOTRACEBACK=single

EXPOSE 5050
ENTRYPOINT ["/src/checkoutservice"]
