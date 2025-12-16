postgres18:
	docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine

start:
	docker start postgres18

stop:
	docker stop postgres18

createdb:
	docker exec -it postgres18 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres18 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/sqlc Store

.PHONY: createdb dropdb migrateup migratedown start stop postgres18 sqlc test server mock