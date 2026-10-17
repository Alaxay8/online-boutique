<!-- <p align="center">
<img src="/src/frontend/static/icons/Hipster_HeroLogoMaroon.svg" width="300" alt="Online Boutique" />
</p> -->
![Continuous Integration](https://github.com/GoogleCloudPlatform/microservices-demo/workflows/Continuous%20Integration%20-%20Main/Release/badge.svg)

> [!NOTE]
> **Educational Fork & DevOps Sandbox**  
> This repository is a learning-focused fork of the original [GoogleCloudPlatform/microservices-demo](https://github.com/GoogleCloudPlatform/microservices-demo).  
> It was created to study and practice modern Kubernetes and DevOps patterns, including:
> - **Event-driven messaging:** Integrated Apache Kafka (KRaft mode), Confluent Schema Registry, and Kafka REST Proxy.
> - **Microservice enhancements:** Extended `checkoutservice` (Go) to publish order events to Kafka; added dark/light UI themes and currency support to `frontend` (Go).
> - **Container registry & CI/CD workflows:** Configured deployment with GitHub Container Registry (`ghcr.io/alaxay8/*`) and Kubernetes `imagePullSecrets`.
> - **Kubernetes management:** Advanced orchestration using Kustomize and dedicated namespaces.

**Online Boutique** is a cloud-first microservices demo application. The application is a web-based e-commerce app where users can browse items, add them to the cart, and purchase them.

This application demonstrates modern microservice architecture patterns using Kubernetes, gRPC, and asynchronous event-driven messaging with Apache Kafka. It works on any Kubernetes cluster (GKE, Minikube, Kind, etc.).

## Architecture

**Online Boutique** is composed of 11 microservices written in different languages communicating over gRPC, paired with an **Apache Kafka** event streaming infrastructure:

[![Architecture of microservices](/docs/img/architecture-diagram.png)](/docs/img/architecture-diagram.png)

Find **Protocol Buffers Descriptions** in the [`./protos` directory](/protos).

### Microservices

| Service | Language | Description |
|---|---|---|
| [frontend](/src/frontend) | Go | Exposes an HTTP server to serve the website. Features anonymous session handling, multi-currency support (including USD, EUR, RUB), and dark/light UI theme switcher (`auto`, `light`, `dark`). |
| [cartservice](/src/cartservice) | C# | Stores items in the user's shopping cart in Redis and retrieves them. |
| [productcatalogservice](/src/productcatalogservice) | Go | Provides the list of products from JSON and supports search and individual product lookup. |
| [currencyservice](/src/currencyservice) | Node.js | Converts money amounts to requested currencies using real exchange rates from the ECB. |
| [paymentservice](/src/paymentservice) | Node.js | Charges the given credit card info (mock) with the given amount and returns a transaction ID. |
| [shippingservice](/src/shippingservice) | Go | Provides shipping cost estimates and tracking IDs for order fulfillment (mock). |
| [emailservice](/src/emailservice) | Python | Sends users an order confirmation email (mock). |
| [checkoutservice](/src/checkoutservice) | Go | Retrieves user cart, prepares order, coordinates payment/shipping/email, and publishes order placed events to Kafka. |
| [recommendationservice](/src/recommendationservice) | Python | Recommends other products based on items in the cart. |
| [adservice](/src/adservice) | Java | Provides text ads based on contextual keywords. |
| [loadgenerator](/src/loadgenerator) | Python/Locust | Continuously generates simulated user shopping traffic. |

### Event Streaming (Apache Kafka)

* **Kafka Broker:** Apache Kafka 3.7.0 deployed in KRaft mode (KRaft broker/controller) for asynchronous event streaming.
* **Schema Registry:** Confluent Schema Registry (7.5.0) for schema governance.
* **Kafka REST Proxy:** Confluent Kafka REST (7.5.0) exposing an HTTP interface (`kafka-rest:8082`) for publishing and consuming events.
* **Topic Initializer:** `kafka-init-orders-topic` Job that creates the `orders` topic (3 partitions, 7 days retention).
* **Event Publishing:** When an order is completed, `checkoutservice` sends an `orderPlacedEvent` payload to `${KAFKA_REST_URL}/topics/orders`.

## Screenshots

| Home Page | Checkout Screen |
|---|---|
| [![Screenshot of store homepage](/docs/img/online-boutique-frontend-1.png)](/docs/img/online-boutique-frontend-1.png) | [![Screenshot of checkout screen](/docs/img/online-boutique-frontend-2.png)](/docs/img/online-boutique-frontend-2.png) |

## Quickstart

### Prerequisites
- Kubernetes cluster (GKE, Minikube, Kind, or Docker Desktop)
- `kubectl` configured with cluster access
- GitHub account and Personal Access Token (PAT) with `read:packages` permission (to pull images from `ghcr.io/alaxay8/*`)

### 1. Clone the repository

```sh
git clone https://github.com/Alaxay8/online-boutique.git
cd online-boutique
```

### 2. Create the namespace and registry credentials

The application manifests and Kafka stack are configured for the `shop` namespace:

```sh
kubectl create namespace shop

kubectl create secret docker-registry ghcr-cred \
  --docker-server=ghcr.io \
  --docker-username=<GITHUB_USERNAME> \
  --docker-password=<GITHUB_TOKEN> \
  -n shop
```

### 3. Deploy the application

Deploy using Kustomize:

```sh
kubectl apply -k kubernetes-manifests -n shop
```

### 4. Wait for the pods to be ready

```sh
kubectl get pods -n shop
```

After a few minutes, all pods should be in the `Running` state:

```
NAME                                       READY   STATUS      RESTARTS   AGE
adservice-76bdd69666-ckc5j                 1/1     Running     0          2m
cartservice-66d497c6b7-dp5jr               1/1     Running     0          2m
checkoutservice-666c784bd6-4jd22           1/1     Running     0          2m
currencyservice-5d5d496984-4jmd7           1/1     Running     0          2m
emailservice-667457d9d6-75jcq              1/1     Running     0          2m
frontend-6b8d69b9fb-wjqdg                  1/1     Running     0          2m
kafka-5d8f6d76d4-8jkmz                     1/1     Running     0          2m
kafka-init-orders-topic-66d49              0/1     Completed   0          1m
kafka-rest-7f64875cb4-z5k8l                1/1     Running     0          2m
loadgenerator-665b5cd444-gwqdq             1/1     Running     0          2m
paymentservice-68596d6dd6-bf6bv            1/1     Running     0          2m
productcatalogservice-557d474574-888kr     1/1     Running     0          2m
recommendationservice-69c56b74d4-7z8r5     1/1     Running     0          2m
redis-cart-5f59546cdd-5jnqf                1/1     Running     0          2m
schema-registry-64c96b659c-l9c8x           1/1     Running     0          2m
shippingservice-6ccc89f8fd-v686r           1/1     Running     0          2m
```

### 5. Access the web frontend

- **On cloud clusters (GKE/EKS/AKS):**
  ```sh
  kubectl get service frontend-external -n shop
  ```
  Open `http://<EXTERNAL-IP>` in your browser.

- **On local clusters (Minikube / Kind):**
  ```sh
  kubectl port-forward service/frontend 8080:80 -n shop
  ```
  Open `http://localhost:8080` in your browser.

## Additional Documentation

- [Development Guide](docs/development-guide.md): Local development, testing, and running with Skaffold.
- [Kubernetes Manifests Guide](kubernetes-manifests/README.md): Details about YAML definitions, Kafka components, and GHCR configuration.
- [Checkout Service Documentation](src/checkoutservice/README.md): Details on gRPC endpoints, Kafka event streaming, and environment variables.
- [Frontend Service Documentation](src/frontend/README.md): Theme switching, multi-currency support, and HTTP endpoints.
- [Adding a New Microservice](docs/adding-new-microservice.md): Guide to extending Online Boutique.
- [Kustomize Variations](kustomize/README.md): Customizing deployments with Kustomize components.

## Additional Deployment Options

- **Terraform**: [See instructions](/terraform) to deploy Online Boutique using [Terraform](https://www.terraform.io/intro).
- **Istio / Service Mesh**: [See instructions](/kustomize/components/service-mesh-istio/README.md) to deploy alongside an Istio-backed service mesh.
- **AI Shopping Assistant**: [See instructions](/kustomize/components/shopping-assistant/README.md) to deploy an assistant that suggests products based on an image.
- **Helm**: [See instructions](/helm-chart/README.md) to deploy using the experimental Helm chart.

