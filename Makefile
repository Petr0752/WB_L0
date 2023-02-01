LOCAL_DB_NAME:=orders_local
LOCAL_DB_DSN:=host=127.0.0.1 port=54320 dbname=orders_local user=orders_user password=orders_password sslmode=disable

docker-up-local:
	docker-compose --env-file ./configs/docker_local.env up -d --build
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" up

docker-down-local:
	docker-compose down

db-reset-local:
	psql -c "drop database if exists $(LOCAL_DB_NAME)"
	psql -c "create database $(LOCAL_DB_NAME)"
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" up
	make db-gen-structure

db-create-migration:
	goose -dir db/migrations create name sql

db-migrate:
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" up
	make db-gen-structure

db-migrate-down:
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" down
	make db-gen-structure

db-gen-structure:
	pg_dump "$(LOCAL_DB_DSN)" --schema-only --no-owner --no-privileges --no-tablespaces --no-security-labels --no-comments > db/structure.sql

run:
	go run cmd/main.go

send-test-message:
	go run internal/nats/test-publisher/main.go
