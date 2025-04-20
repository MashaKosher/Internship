DOCS_DIR = ./docs
MAIN_FILE = ./cmd/core/main.go 


# Генерация документации для cmd/auth
generate-swagger:
	swag init --generalInfo $(MAIN_FILE)


SQLC_YAML_PATH = internal/config/sqlc.yaml

check-queries:
	sqlc vet --file $(SQLC_YAML_PATH)


generate-queries:
	sqlc generate --file $(SQLC_YAML_PATH)


MIGARTIONS_PATH = migrations
MIGARTIONS_FORMAT = sql
MIGRATION_NAME ?= new_migration
MIGRATION_DRIVER = postgres
MIGRATION_USER = postgres
MIGRATION_PASSWORD = postgres
MIGRATION_HOST = localhost
MIGRATION_PORT = 5435
MIGRATION_DB = core

#  make create-migrations MIGRATION_NAME=migration1 
create-migrations:
	goose -dir $(MIGARTIONS_PATH) create $(MIGRATION_NAME) $(MIGARTIONS_FORMAT)   


up-migrations:
	goose $(MIGRATION_DRIVER) $(MIGRATION_DRIVER)://$(MIGRATION_USER):$(MIGRATION_PASSWORD)@$(MIGRATION_HOST):$(MIGRATION_PORT)/$(MIGRATION_DB) up -dir $(MIGARTIONS_PATH)


down-migrations:
	goose $(MIGRATION_DRIVER) $(MIGRATION_DRIVER)://$(MIGRATION_USER):$(MIGRATION_PASSWORD)@$(MIGRATION_HOST):$(MIGRATION_PORT)/$(MIGRATION_DB) down -dir $(MIGARTIONS_PATH)


clear-migrations:
	goose $(MIGRATION_DRIVER) $(MIGRATION_DRIVER)://$(MIGRATION_USER):$(MIGRATION_PASSWORD)@$(MIGRATION_HOST):$(MIGRATION_PORT)/$(MIGRATION_DB) down-to 0 -dir $(MIGARTIONS_PATH)


version-migrations:
	goose $(MIGRATION_DRIVER) $(MIGRATION_DRIVER)://$(MIGRATION_USER):$(MIGRATION_PASSWORD)@$(MIGRATION_HOST):$(MIGRATION_PORT)/$(MIGRATION_DB) version