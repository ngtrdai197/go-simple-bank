include app.env

# Variables for binding arguments from cmd
MIGRATION_FILE_NAME?=default_migration_file_name

.PHONY: createdb dropdb migrateup migratedown sqlc test server mock createmigration proto producer consumer

consumer:
	go run ./cmd/kafka/consumer/main.go

producer:
	go run ./cmd/kafka/producer/main.go

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ngtrdai197/simple-bank/db/sqlc Store

createdb:
	docker exec -it c6af98ca2ee2 createdb --username=root --owner=root simple_bank_final

dropdb:
	docker exec -it c6af98ca2ee2 dropdb --username=root --owner=root simple_bank_final

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

createmigration:
	migrate create --ext sql --dir db/migration -seq $(MIGRATION_FILE_NAME)

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...