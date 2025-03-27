# Makefile

# Переменные
SWAG = swag
DOCS_DIR = ./docs
AUTH_GO = internal/service/auth.go


# Генерация документации для cmd/auth
generate-swagger:
	$(SWAG) init --dir cmd/auth --output $(DOCS_DIR)
	$(SWAG) init -g $(AUTH_GO)