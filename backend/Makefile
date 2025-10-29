# Nombre del binario de la aplicación
APP_NAME=myapp
# Variables para la base de datos
DB_URL="postgres://root:root@database:5432/myapp?sslmode=disable"

##NOTAS
#	wait docker 
#	docker compose up -d
all: run

generate: 
	@echo " Generando código SQLC..."
	sqlc generate
#@echo ">= Generating Templ code..."
#@templ generate

build: generate
	@echo " Compilando aplicación..."
	go build -o $(APP_NAME) .

run: clean-db
	@echo " Iniciando servidor en http://localhost:8080"
	./$(APP_NAME)
clean:
	@rm -f $(APP_NAME)
clean-db: db-up
	@echo "Reseteando la base de datos"
#docker compose exec database psql -U root -d myapp -c "DROP TABLE IF EXISTS procesos CASCADE;"

	psql "postgres://root:root@localhost:5432/myapp?sslmode=disable" -c "DROP TABLE IF EXISTS procesos CASCADE;"
	psql "postgres://root:root@localhost:5432/myapp?sslmode=disable" -f db/schema/schema.sql

db-up: db-down
	@echo " Levantando base de datos PostgreSQL..."
	docker-compose up -d 
	docker compose up -d database

	@echo "Esperando a PostgreSQL..."
	until docker-compose exec database pg_isready -U root -d myapp; do \
		sleep 1; \
	done

db-down:
	@echo " Apagando base de datos..."
	docker-compose down
#con el "&" se consigue que el servidor quede corriendo y se ejecuten las siguientes lineas. El $$! captura el id del proceso que esta corriendo en el servidor para luego desalojarlo. 
test: build db-up
	./$(APP_NAME) & 
	SERVER_PID=$$! 
	sleep 2
	hurl request.hurl
	kill $$SERVER_PID


