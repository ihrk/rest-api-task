test-local: 
	go test -v ./internal/app/...

test:
	sh -c "trap 'docker-compose down' EXIT; \
		docker-compose up --build --force-recreate \
		--remove-orphans --abort-on-container-exit \
		--exit-code-from rest_api_test rest_api_test"

lint:
	sh -c "trap 'docker-compose down' EXIT; \
		docker-compose up --build --force-recreate \
		--remove-orphans --abort-on-container-exit \
		--exit-code-from rest_api_lint rest_api_lint"

lint-local:
	golangci-lint run --timeout 3m ./...

ci: lint test

run:
	docker-compose up --build --force-recreate -d rest_api

migrate-up:
	go run ./cmd/rest-api-task/main.go -migrate up