export COMPOSE_BAKE=true
.PHONY: build-prod build-qa test-docker docs lint

test-docker:
	docker build -f deployment/Dockerfile --target test -t account-test .
	docker run --rm account-test

build-qa:
	docker build -f deployment/Dockerfile --target qa -t account:qa .

build-prod:
	docker build -f deployment/Dockerfile --target prod -t account:latest .

docs:
	swag init -g internal/server/server.go

lint:
	go fmt ./... && golangci-lint run ./...