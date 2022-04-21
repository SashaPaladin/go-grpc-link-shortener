build:
	protoc --proto_path=proto proto/*.proto --go_out=proto
	protoc --proto_path=proto proto/*.proto --go-grpc_out=proto
	docker-compose build link-shortener

run_inmemory:
	docker-compose up link-shortener-inmemory

run_postgres:
	docker-compose up link-shortener-postgres

migrate:
	migrate -path ./schema/ -database "postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable" up