postgres-container:
	docker run --name web-dev-tools -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres

migrate-create:
	echo "Creating migration... $(name)"
	migrate create -seq -ext=.sql -dir=./migrations "$(name)"

migrate-up:
	echo "Migrating up..."
	migrate -path ./migrations -database "$(DB_URL)"  up

migrate-up-test:
	echo "Migrating up..."
	migrate -path ./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable  up

test:
	go test -v -cover ./...

lint:
	golangci-lint run

.PHONY: postgres-container postgres-start postgres-stop postgres-remove migrate-create migrate-up test-data lint
