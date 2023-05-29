postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.2-alpine

postgres-createdb:
	docker exec -it postgres15 createdb --username=root --owner=root finances

postgres-dropdb:
	docker exec -it postgres15 dropdb finances

postgres-migrate-up:
	migrate -path internal/repository/postgres/migration -database "postgresql://root:secret@localhost:5432/finances?sslmode=disable" -verbose up 

postgres-migrate-up1:
	migrate -path internal/repository/postgres/migration -database "postgresql://root:secret@localhost:5432/finances?sslmode=disable" -verbose up 1

postgres-migrate-down:
	migrate -path internal/repository/postgres/migration -database "postgresql://root:secret@localhost:5432/finances?sslmode=disable" -verbose down 

postgres-migrate-down1:
	migrate -path internal/repository/postgres/migration -database "postgresql://root:secret@localhost:5432/finances?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/main.go

mock:
	mockgen -package mockdb -destination internal/repository/postgres/sqlc/mock/store.go github.com/saver89/finance-management/internal/repository/postgres/sqlc Store

.PHONY: postgres postgres-createdb postgres-dropdb postgres-migrate-up postgres-migrate-down sqlc test server mock postgres-migrate-down1 postgres-migrate-up1