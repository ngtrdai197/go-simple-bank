include app.env

# Variables for binding arguments from cmd
MIGRATION_FILE_NAME?=default_migration_file_name

.PHONY: createdb dropdb migrateup migratedown sqlc test server mock createmigration

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