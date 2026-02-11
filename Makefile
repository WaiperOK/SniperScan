.PHONY: fmt test lint run-scan run-api docker-up docker-down

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

lint:
	go test ./...

run-scan:
	go run ./cmd/sniperscan scan --target scanme.nmap.org --ports 22,80,443

run-api:
	go run ./cmd/sniperscan serve --addr :8097

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v
