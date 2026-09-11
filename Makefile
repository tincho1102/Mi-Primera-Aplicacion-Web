.PHONY: prepare run-tests clean test

MODULE_NAME = matematicos.com/servidor-go-tp1

#Generado con IA, nos dio un solo comando llamado test, nosotros lo dividimos en 3 comandos
# para separar cada parte.

# 1. Instalacion de módulo, dependencias, sqlc y Docker
prepare:
	go mod init $(MODULE_NAME) || true
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	sqlc generate
	go get github.com/jackc/pgx/v5
	go mod tidy
	docker compose up -d

# 2. Ejecución de pruebas, ejecuta el test
run-tests:
	go test -v ./...

# 3. Limpieza de contenedores y volúmenes, borra lo que tiene las tablas
clean:
	docker compose down -v

# Con este comando hacemos la preparacion, los test y la limpieza
test: prepare run-tests clean