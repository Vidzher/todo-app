include .env
export

# абсолютный путь до корня проекта
export PROJECT_ROOT=${shell pwd}

env-up:
	docker compose up -d app-db

env-down:
	docker compose down app-db

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

# выключает вывод выполняемых команд
.SILENT: env-cleanup
env-cleanup:
	printf "Очистить все volume файлы окружения? Опасность утери данных. [y/N] "; \
	read answ; \
	if [ "$$answ" = "y" ]; then \
  		docker compose down app-db && \
  		rm -rf out/pgdata && \
  		echo "Файлы окружения удалены."; \
	else \
	  	echo "Очистка данных отменена."; \
	fi;

migrate-create:
	@if [ -z "$(seq)"]; then \
  		echo "Отсутствует параметр seq. Пример: migrate-create seq=migration_name"; \
		exit 1; \
	fi; \

	docker compose run --rm pg-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует параметр action. Пример: migrate-create action=up 1"; \
		exit 1; \
	fi; \

	docker compose run --rm pg-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@app-db:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/todoapp/main.go