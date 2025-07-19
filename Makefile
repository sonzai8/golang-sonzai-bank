-include ./.github/workflows/.env
export
MIGRATE_DIR = ./db/migrations

CONN_STRING = postgres://$(DB_DRIVER.DB_USER):$(DB_DRIVER.DB_PASSWORD)@$(DB_DRIVER.DB_HOST):$(DB_DRIVER.DB_PORT)/$(DB_DRIVER.DB_NAME)?sslmode=$(DB_DRIVER.DB_SSLMODE)

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)
#Run all pening migrations (make migrate-up)
migrate-up:
	migrate -path $(MIGRATE_DIR) -database $(CONN_STRING) up
# apply specific migrations version make migrate-togo version=1
migrate-goto:
	migrate -path $(MIGRATE_DIR) -database $(CONN_STRING) goto $(version)

# Rollback the last migration
migrate-down:
	migrate -path $(MIGRATE_DIR) -database $(CONN_STRING) down 1
# Rollback n migrations
migrate-down-n:
	migrate -path $(MIGRATE_DIR) -database $(CONN_STRING) down $(n)

# force migrate version (user with caution example : make migrate-force version=1)
migrate-force:
	migrate -path $(MIGRATE_DIR) -database $(CONN_STRING) force $(version)
#
# Drop everything (include schema migration)
migrate-drop:
	migrate -path $(MIGRATE_DIR) -database $(CONN_STRING) drop

migrateup:
	migrate -path $(MIGRATE_DIR) -database $(CONN_STRING) -verbose up

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/sonzai8/golang-sonzai-bank/db/sqlc Store

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        proto/*.proto
server:
	go run main.go

evans:
	evans --host localhost --port 9090 -r repl
.PHONY: migrate-create migrate-up migrate-down migrate-down-n migrate-goto\
 		migrate-force migrate-drop sqlc test migrateup mock proto server evans