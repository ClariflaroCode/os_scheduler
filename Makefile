# Nombre del binario de la aplicación
APP_NAME=my-app
# Variables para la base de datos
DB_URL="postgres://user:password@localhost:5432/mydb?sslmode=disable"

##NOTAS
#	wait docker 
#	docker compose up -d
all: run

generate: 
	@echo " Generando código SQLC..."
	sqlc generate

build: generate
	@echo " Compilando aplicación..."
	go build -o $(APP_NAME) .

run: db-up
	@echo " Iniciando servidor en http://localhost:8080"
	./$(APP_NAME)

clean:
	@echo " Limpiando binarios..."
	rm -f $(APP_NAME)
clean-db: 
	@echo "Reseteando la base de datos"
	psql "postgres://root:root@localhost:5432/myapp?sslmode=disable" -c "DROP TABLE IF EXISTS procesos CASCADE;"
	psql "postgres://root:root@localhost:5432/myapp?sslmode=disable" -f db/schema/schema.sql

db-up: db-down
	@echo " Levantando base de datos PostgreSQL..."
	docker-compose up -d 
	@echo "Esperando a PostgreSQL..."
	until docker-compose exec database pg_isready -U root -d myapp; do \
		sleep 1; \
	done

db-down: clean-db
	@echo " Apagando base de datos..."
	docker-compose down

test: all
	@echo " Ejecutando script de pruebas requests.hurl..."
	hurl --test request.hurl

