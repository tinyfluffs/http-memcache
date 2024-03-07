# http/memcache

HTTP memory cache. Microservice is built to the [12 Factor App](https://12factor.net/) guidelines.

## Environment

Just pull the latest dependencies.
```shell
go mod download
```

## Running Tests

### Tests

```shell
go test ./...
```

### Benchmarks

```shell
go test -bench=. -benchmem -benchtime=4s -timeout=30m ./...
```

## Building to Minikube

Assuming local builds for now, have Docker's new buildx tool ready to go:

## Container Build
#### macOS and Linux

From the project directory, execute:
```shell
# If using minikube:-
minikube start
eval $(minikube -p minikube docker-env)

# Finally
docker buildx build -t http-memcache:latest -f ./build/package/Dockerfile .
```

## Minikube Deploy
```shell
kubectl apply -f ./deployments/memcache-deployment.yaml
kubectl apply -f ./deployments/memcache-service.yaml
```

## Minikube Service URL

```shell
minikube service http-memcache-entry --url
```
