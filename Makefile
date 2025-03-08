test:
	go test ./... -v

lint:
	golangci-lint run

build:
	go build -o fabrik main.go

schema-create:
	migrate create --ext=sql --dir=db/migrations/ -seq $(name)

start-test-db:
	docker compose up -d

stop-test-db:
	docker compose down -v

psql:
	docker exec -it fabrik-db-1 psql -U fabrik

run-migrations:
	go run main.go migrate -c ./creds/fabrik.creds.yaml

.PHONY: test lint build schema-create start-test-db stop-test-db run-migrations
