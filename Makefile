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
