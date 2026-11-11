.PHONY: prepare templ generate run test clean all

MODULE_NAME = matematicos.com/servidor-go-tp1

# TP5: version fija de templ. Tiene que coincidir con la de go.mod para que el
# codigo generado (_templ.go) sea compatible con la libreria que usa el servidor.
# v0.3.1020 es la ultima version que funciona con Go 1.25.
TEMPL_VERSION = v0.3.1020
# Usamos la ruta completa del binario para no depender de que $GOPATH/bin este en el PATH.
TEMPL = $(shell go env GOPATH)/bin/templ

#	Generado con IA, nos dio un solo comando llamado test, nosotros lo dividimos en 3 comandos
#	para separar cada parte. Luego separamos uno mas para levantar el servidor para la entrega del tp3.
#	Para resolver el problema de que el contenedor de la base de datos no estaba listo,
#	agregamos un flag pg_isready para esperar a que la base de datos esté lista antes de ejecutar. Ademas para
#	que se ejecute en segundo plano, usamos el comando & y guardamos el PID en un archivo server.pid para
#	poder matarlo despues con kill.
#	En el TP5 agregamos los targets templ y generate para compilar las vistas .templ a codigo Go.
#	Ademas agregamos un target execute que hace todo menos el clean, para poder ver los logs del 
#	servidor y el resultado de las pruebas.

# 1. Instalacion de módulo, dependencias, sqlc, templ y Docker. usamos el flag pg_isready para esperar
# a que la base de datos esté lista antes de ejecutar.
#	El orden importa: primero se genera el codigo de sqlc (db/sqlc) y de templ (views/*_templ.go),
#	y recien despues go mod tidy, que necesita ver esos imports para no borrar dependencias.
prepare: templ
	@echo "Inicializando módulo de Go y dependencias..."
	go mod init $(MODULE_NAME) || true
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate
	$(TEMPL) generate
	go get github.com/a-h/templ@$(TEMPL_VERSION)
	go mod tidy
	@echo "Levantando contenedores de Docker (PostgreSQL)..."
	docker compose up -d
	@echo "Esperando a que la base de datos esté lista..."
	@until docker compose exec -T db pg_isready -h localhost -U cliente -d mate_db -q; do sleep 1; done
	@echo "Base de datos lista y conectada"

#  Descarga/verificación del CLI de templ. Si ya esta instalado en la version correcta no hace nada,
#	si falta o es otra version lo instala con go install.
templ:
	@if [ -x "$(TEMPL)" ] && [ "$$($(TEMPL) version)" = "$(TEMPL_VERSION)" ]; then \
		echo "templ $(TEMPL_VERSION) ya esta instalado"; \
	else \
		echo "Instalando templ $(TEMPL_VERSION)..."; \
		go install github.com/a-h/templ/cmd/templ@$(TEMPL_VERSION); \
	fi

# Compila los componentes .templ de views/ a archivos _templ.go
generate: templ
	@echo "Generando vistas con templ..."
	$(TEMPL) generate

#	Ejecución de main.go, lo hace en segundo plano para dejar libre la consola y poder ejecutar las pruebas.
#	Redirige la salida a server.log y guarda el PID en server.pid
#	Depende de generate: las vistas se compilan siempre antes de go build.
run: generate
	go build -o server .
	./server > server.log 2>&1 & echo $$! > server.pid
	@for i in 1 2 3 4 5 6 7 8 9 10; do curl -s -o /dev/null localhost:8080 && break; sleep 1; done
	@echo "Servidor levantado en http://localhost:8080 (logs en server.log)"


#  Limpieza de contenedores y volúmenes, borra lo que tiene las tablas.
#	Ademas mata el proceso del servidor que se ejecuta en segundo plano.
clean:
	@echo "Deteniendo servidor, contenedores y borrando volúmenes de Docker..."
	-kill $$(cat server.pid)
	rm -f server server.pid server.log
	docker compose down -v

# Todo junto y en orden, para que se pueda ejecutar con un solo comando
all: prepare run test clean

# Ejecuta todo sin hacer el clean
execute: prepare run 