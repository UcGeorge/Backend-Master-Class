<!-- Section 1 -->
brew install sqlc
docker pull postgres:18-alpine
docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine
docker exec -it postgres18 psql -U root
docker logs postgres18
brew install golang-migrate
migrate create -ext sql -dir db/migration -seq init_schema
docker exec -it postgres18 /bin/sh
createdb --username=root --owner=root simple_bank
psql simple_bank
dropdb simple_bank
exit

<!-- Section 2 -->
