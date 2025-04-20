# Makefile

# Переменные
SWAG = swag
DOCS_DIR = ./docs
# AUTH_GO = internal/service/auth.go
MAIN_FILE = ./cmd/auth/main.go 


# Генерация документации для cmd/auth
generate-swagger:
	$(SWAG) init --generalInfo $(MAIN_FILE)
	# $(SWAG) init --dir cmd/auth --output $(DOCS_DIR)
	# $(SWAG) init -g $(AUTH_GO)