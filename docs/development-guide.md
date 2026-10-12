# Development Guide 

This doc explains how to build and run the Online Boutique source code locally using the `skaffold` command-line tool or `kubectl`.

## Prerequisites

- [Docker for Desktop](https://www.docker.com/products/docker-desktop)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) (can be installed via `gcloud components install kubectl` for Option 1 - GKE)
- [skaffold **2.0.2+**](https://skaffold.dev/docs/install/) (latest version recommended), a tool that builds and deploys Docker images in bulk.
- Clone the repository:
    ```sh
    git clone https://github.com/Alaxay8/online-boutique.git
    cd online-boutique/
    ```
- [GitHub Personal Access Token](https://github.com/settings/tokens) with `read:packages` permission to pull images from GitHub Container Registry (`ghcr.io/alaxay8/*`).
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) or [Kind](https://kind.sigs.k8s.io/) (for Option 2 - Local Cluster).

## Setup Container Registry Credentials

Before deploying pre-built images or running manifests, create the `ghcr-cred` secret in your cluster:

```sh
# Create namespace
kubectl create namespace shop

# Create registry secret
kubectl create secret docker-registry ghcr-cred \
  --docker-server=ghcr.io \
  --docker-username=<GITHUB_USERNAME> \
  --docker-password=<GITHUB_TOKEN> \
  -n shop
```

## Option 1: Google Kubernetes Engine (GKE)

> 💡 Recommended if you're using Google Cloud and want to try it on
> a realistic cluster. **Note**: If your cluster has Workload Identity enabled, 
> [see these instructions](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable).

1.  Create a Google Kubernetes Engine cluster and make sure `kubectl` is pointing
    to the cluster:

    ```sh
    gcloud services enable container.googleapis.com
    ```

    ```sh
    gcloud container clusters create-auto demo --region=us-central1
    ```

    ```sh
    kubectl get nodes
    ```

2.  Enable Artifact Registry (AR) on your GCP project (if building and pushing custom images) and configure the `docker` CLI:

    ```sh
    gcloud services enable artifactregistry.googleapis.com
    gcloud artifacts repositories create microservices-demo \
      --repository-format=docker \
      --location=us
    gcloud auth configure-docker -q
    ```

3.  Configure namespace and credentials:

    ```sh
    kubectl create namespace shop
    kubectl create secret docker-registry ghcr-cred \
      --docker-server=ghcr.io \
      --docker-username=<GITHUB_USERNAME> \
      --docker-password=<GITHUB_TOKEN> \
      -n shop
    ```

4.  In the root of this repository, run:

    ```sh
    skaffold run --default-repo=us-docker.pkg.dev/PROJECT_ID/microservices-demo -n shop
    ```

    Where `PROJECT_ID` is replaced by your Google Cloud project ID.

    This command:
    - Builds the container images.
    - Pushes them to Artifact Registry.
    - Applies `./kubernetes-manifests` deploying the application and Kafka to Kubernetes.

5.  Find the IP address of your application, then visit the application in your browser:

    ```sh
    kubectl get service frontend-external -n shop
    ```

6.  Navigate to `http://EXTERNAL-IP` to access the web frontend.

## Option 2 - Local Cluster

1. Launch a local Kubernetes cluster with one of the following tools:

    - To launch **Minikube** (tested with Ubuntu Linux / macOS). Ensure the local cluster has at least:
      - 4 CPUs
      - 4.0 GiB memory (recommended 6.0 GiB for running Kafka alongside services)
      - 32 GB disk space

      ```shell
      minikube start --cpus=4 --memory 6144 --disk-size 32g
      ```

    - To launch **Docker for Desktop** (tested with Mac/Windows). Go to Preferences:
      - choose “Enable Kubernetes”,
      - set CPUs to at least 4, and Memory to at least 6.0 GiB
      - on the "Disk" tab, set at least 32 GB disk space

    - To launch a **Kind** cluster:

      ```shell
      kind create cluster
      ```

2. Create the `shop` namespace and GHCR credentials secret:

    ```shell
    kubectl create namespace shop
    kubectl create secret docker-registry ghcr-cred \
      --docker-server=ghcr.io \
      --docker-username=<GITHUB_USERNAME> \
      --docker-password=<GITHUB_TOKEN> \
      -n shop
    ```

3. Verify connection:

    ```shell
    kubectl get nodes
    ```

4. Deploy the application:

    - **Using pre-built images with Kustomize:**
      ```shell
      kubectl apply -k kubernetes-manifests -n shop
      ```

    - **Or building locally with Skaffold:**
      ```shell
      skaffold run -n shop
      ```
      If you need to rebuild images automatically on code changes, run `skaffold dev -n shop`.

5. Verify all pods (including Kafka, Schema Registry, and microservices) are running:

    ```shell
    kubectl get pods -n shop
    ```

6. Forward port to the frontend service:

    ```shell
    kubectl port-forward deployment/frontend 8080:8080 -n shop
    ```

7. Navigate to `http://localhost:8080` to access the web frontend.

## Adding a new microservice

In general, the set of core microservices for Online Boutique is fairly complete, but optional microservices can be added.

See the [Adding a new microservice](adding-new-microservice.md) guide for instructions on how to add a new microservice.

## Cleanup

If you deployed with Skaffold:

```shell
skaffold delete -n shop
```

If you deployed with `kubectl`:

```shell
kubectl delete -k kubernetes-manifests -n shop
```
