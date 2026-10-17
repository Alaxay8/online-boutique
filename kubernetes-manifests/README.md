# Kubernetes Manifests

This directory contains Kubernetes YAML manifests for deploying Online Boutique along with the Apache Kafka event streaming infrastructure.

## Components

* **Microservices (11 services):**
  * `adservice.yaml`, `cartservice.yaml`, `checkoutservice.yaml`, `currencyservice.yaml`, `emailservice.yaml`, `frontend.yaml`, `loadgenerator.yaml`, `paymentservice.yaml`, `productcatalogservice.yaml`, `recommendationservice.yaml`, `shippingservice.yaml`
* **Event Streaming & Messaging:**
  * `kafka.yaml`: Deploys Apache Kafka (broker/controller in KRaft mode), Confluent Schema Registry, Confluent Kafka REST Proxy, and a one-time setup Job (`kafka-init-orders-topic`) creating the `orders` topic with 3 partitions and a 7-day retention period.
* **Kustomize configuration:**
  * `kustomization.yaml`: Bundles all microservices and `kafka.yaml` into a single Kustomize target.

## Container Registry & Credentials

The manifests are configured with images hosted on GitHub Container Registry (`ghcr.io/alaxay8/*`).

Each service deployment and service account specifies `imagePullSecrets: ghcr-cred`. Before deploying, you must create this secret in your cluster:

```sh
kubectl create secret docker-registry ghcr-cred \
  --docker-server=ghcr.io \
  --docker-username=<GITHUB_USERNAME> \
  --docker-password=<GITHUB_TOKEN> \
  -n <NAMESPACE>
```

> **Note:** The GitHub personal access token (PAT) must have `read:packages` permission.

## Namespaces & Deployment

Resources in `kafka.yaml` are explicitly set to `namespace: shop`. For proper internal DNS resolution (`kafka-rest:8082` and `kafka:9092`), it is recommended to deploy the application into the `shop` namespace:

```sh
# 1. Create the namespace
kubectl create namespace shop

# 2. Create the image pull secret in the namespace
kubectl create secret docker-registry ghcr-cred \
  --docker-server=ghcr.io \
  --docker-username=<GITHUB_USERNAME> \
  --docker-password=<GITHUB_TOKEN> \
  -n shop

# 3. Apply manifests using Kustomize
kubectl apply -k kubernetes-manifests -n shop
```

## Local Development with Skaffold

If you are developing locally and rebuilding images using [Skaffold](https://skaffold.dev), Skaffold will build the images defined in `skaffold.yaml` and deploy them to your target cluster.
