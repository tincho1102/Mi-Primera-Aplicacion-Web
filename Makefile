.PHONY: prepare run-tests clean test

MODULE_NAME = matematicos.com/servidor-go-tp1

#	Generado con IA, nos dio un solo comando llamado test, nosotros lo dividimos en 3 comandos
#	para separar cada parte. Luego separamos uno mas para levantar el servidor para la entrega del tp3.
#	Para resolver el problema de que el contenedor de la base de datos no estaba listo,
#	agregamos un flag pg_isready para esperar a que la base de datos esté lista antes de ejecutar. Ademas para
#	que se ejecute en segundo plano, usamos el comando & y guardamos el PID en un archivo server.pid para 
#	poder matarlo despues con kill.

# 1. Instalacion de módulo, dependencias, sqlc y Docker. usamos el flag pg_isready para esperar
# a que la base de datos esté lista antes de ejecutar.
prepare:
	@echo "Inicializando módulo de Go y dependencias..."
	go mod init $(MODULE_NAME) || true
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate
	go get github.com/jackc/pgx/v5
	go mod tidy
	@echo "Levantando contenedores de Docker (PostgreSQL)..."
	docker compose up -d
	@echo "Esperando a que la base de datos esté lista..."
	@until docker compose exec -T db pg_isready -h localhost -U cliente -d mate_db -q; do sleep 1; done
	@echo "Base de datos lista y conectada"

# 2. Ejecución de main.go, lo hace en segundo plano para dejar libre la consola y poder ejecutar las pruebas.
#	Redirige la salida a server.log y guarda el PID en server.pid
run:
	go build -o server .
	./server > server.log 2>&1 & echo $$! > server.pid
	@for i in 1 2 3 4 5 6 7 8 9 10; do curl -s -o /dev/null localhost:8080 && break; sleep 1; done
	@echo "Servidor levantado (logs en server.log)"

# 3. Ejecución de pruebas automatizadas de endpoints con Hurl
test:
	@echo "Ejecutando pruebas automatizadas de endpoints con Hurl..."
	@hurl --test request.hurl

# 4. Limpieza de contenedores y volúmenes, borra lo que tiene las tablas.
#	Ademas mata el proceso del servidor que se ejecuta en segundo plano.
clean:
	@echo "Deteniendo servidor, contenedores y borrando volúmenes de Docker..."
	-kill $$(cat server.pid)
	rm -f server server.pid server.log
	docker compose down -v

# Todo junto y en orden, para que se pueda ejecutar con un solo comando
all: prepare run test clean
