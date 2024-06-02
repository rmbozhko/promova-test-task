migratecreate:
	docker run --rm -v $(PWD)/db/migration:/migrations migrate/migrate create -ext=sql -dir=/migrations/ -seq init_schema

migrateup:
	docker run --rm -v $(shell pwd)/db/migration:/migration --network host migrate/migrate -path=/migration/ -database "postgresql://postgres:postgres@localhost:5432/promova_test_task?sslmode=disable" -verbose up

migratedown:
	docker run --rm -v $(shell pwd)/db/migration:/migration --network host migrate/migrate -path=/migration/ -database "postgresql://postgres:postgres@localhost:5432/promova_test_task?sslmode=disable" down -all

mock:
	docker run --rm -v $(PWD):/app -w /app ekofr/gomock:latest mockgen -package=mockdb -destination=db/mock/querier.go  -source=db/sqlc/querier.go

sqlc:
	docker run --rm -v $(PWD):/app -w /app kjconroy/sqlc generate

swag:
	${HOME}/go/bin/swag init -d cmd,api,db/sqlc

test:
	go test -v -cover ./...

server:
	go run cmd/main.go

deploy:
	docker compose up -d

.PHONY: migratecreate migrateup migratedown mock sqlc swag server deploy