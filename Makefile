# Makefile
run-service:
	docker compose up --build -d
stop-service:
	docker-compose stop
rm-service:
	docker-compose down
run-tests:
	go test -v -cover ./...
test-cov:
	go test -coverprofile=coverage.out ./... & go tool cover -html=coverage.out
test-load:
	k6 run ./src/test/load_test.js

# Линтинг
lint:
	golangci-lint run ./... --config .golangci.yml

lint-fix:
	golangci-lint run ./... --config .golangci.yml --fix