SERVICES := $(shell find . -type f -name 'main.go' -exec dirname {} \; | sed 's|^\./||')

ifneq ("$(wildcard .env)","")
	include .env
endif

# dev cmds
.PHONY: tidy format lint test

format:
	go fmt ./...

lint:
	golangci-lint run ./...

test:
	go test -v ./...

# gen cmds
.PHONY: gen-openapi

gen-openapi:
	@for service in $(SERVICES); do \
		echo "Generating OpenAPI code for $$service..."; \
		oapi-codegen --config=$$service/oapi-codegen.yaml -o $$service/api/api.gen.go $$service/oas/openapi.yaml; \
	done

# build cmds
.PHONY: build start docker-build docker-run

build:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		go build -o $$service/bin/main ./$$service; \
	done

start:
	@for service in $(SERVICES); do \
		echo "Starting $$service..."; \
		go run $$service/main.go; \
	done

docker-build:
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		docker build -t $$service:latest -f $$service/Dockerfile .; \
	done

docker-run:
	@for service in $(SERVICES); do \
		echo "Running Docker container for $$service..."; \
		docker run -p 8080:8080 --env-file .env $$service; \
	done

# k8s cmds
.PHONY: minikube-create-db-secret minikube-docker-build minikube-deploy minikube-expose

minikube-create-db-secret:
	kubectl create secret generic db-secret --from-literal=DB_PASS=$(DB_PASS)

minikube-docker-build:
	@for service in $(SERVICES); do \
		echo "Building Minikube Docker image for $$service..."; \
		eval $$(minikube -p minikube docker-env); \
		docker build -t $$service:latest -f $$service/Dockerfile .; \
	done

minikube-deploy:
	@for service in $(SERVICES); do \
		echo "Deploying $$service to Minikube..."; \
		kubectl apply -f $$service/k8s/app-deployment.yaml; \
		kubectl apply -f $$service/k8s/app-service.yaml; \
		kubectl apply -f $$service/k8s/pg-deployment.yaml; \
		kubectl apply -f $$service/k8s/pg-service.yaml; \
	done

minikube-expose:
	@for service in $(SERVICES); do \
		echo "Exposing $$service from Minikube..."; \
		minikube service $$service-service --url; \
	done
