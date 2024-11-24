# GO MICRO

- [GO MICRO](#go-micro)
  - [Tools](#tools)
  - [Packages](#packages)
  - [Development](#development)
  - [TODO](#todo)

## Tools

- [Go](https://go.dev/doc/install) - For building and running the microservice.
- [Docker](https://docs.docker.com/desktop) - To containerize the application for deployment.
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) The Kubernetes command-line tool.
- [minikube](https://minikube.sigs.k8s.io/docs/start) - Local Kubernetes cluster for testing.

## Packages

- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - Generates Go server boilerplate from OpenAPI specs.
- [chi](https://github.com/go-chi/chi) - Lightweight and idiomatic routing for HTTP services.
- [ent](https://github.com/ent/ent) - ORM framework for data modeling and querying in Go.
- [pg](https://github.com/lib/pq) - PostgreSQL driver for database connectivity.
- [dotenv](https://github.com/joho/godotenv) - For loading environment variables from `.env` files.

## Development

- [x] Generate Service Definition from OAS
  - `make gen-openapi SERVICE=<service>`
    - ex: [petstore/api/api.gen.go](./petstore/api/api.gen.go)
- [x] Add ent and Create Schema
  - `cd <service>`
  - `go run -mod=mod entgo.io/ent/cmd/ent new <entity>`
    - ex: [petstore/ent/schema/pet.go](./petstore/ent/schema/pet.go)
  - `go generate ./ent`
- [x] Implement Service Endpoints
  - ex: [petstore/api/impl.go](./petstore/api/impl.go)
- [X] Deploy Locally to minikube
  - Start minikube
    - `minikube start`
  - Add `DB_PASS` (any value) to a [.env](./.env.example) file
    - `make create-db-secret`
  - Build + Deploy + Expose Port
    - `make minikube-docker-build`
    - `make minikube-deploy`
    - `make minikube-expose`

## TODO

- [ ] Unit Tests for service implementation
- [ ] Another Micro-Service generated from OAS
- [ ] Deploy to a live environment (probably AWS)
- [ ] Add Helm and Argo to manage config and deployments
