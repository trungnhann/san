MIGRATION_DIR := ./internal/db/migration

migration:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migration name=create_users_table"; \
		exit 1; \
	fi
	@mkdir -p $(MIGRATION_DIR)
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	touch $(MIGRATION_DIR)/$${TIMESTAMP}_$(name).up.sql; \
	touch $(MIGRATION_DIR)/$${TIMESTAMP}_$(name).down.sql; \
	echo "Created:"; \
	echo "  $(MIGRATION_DIR)/$${TIMESTAMP}_$(name).up.sql"; \
	echo "  $(MIGRATION_DIR)/$${TIMESTAMP}_$(name).down.sql"

sqlc:
	@sqlc generate

dev:
	air

swagger:
	swag init -g cmd/server/main.go