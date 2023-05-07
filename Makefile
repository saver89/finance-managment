postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.2-alpine

postgres-createdb:
	docker exec -it postgres15 createdb --username=root --owner=root finances

postgres-dropdb:
	docker exec -it postgres15 dropdb finances

postgres-migrate-up:
	migrate -path internal/repository/postgres/migration -database "postgresql://root:secret@localhost:5432/finances?sslmode=disable" -verbose up 

postgres-migrate-down:
	migrate -path internal/repository/postgres/migration -database "postgresql://root:secret@localhost:5432/finances?sslmode=disable" -verbose down 

.PHONY: postgres postgres-createdb postgres-dropdb postgres-migrate-up postgres-migrate-down