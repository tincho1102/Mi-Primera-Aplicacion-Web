.PHONY: all prepare deps sqlc generate build run clean

MODULE_NAME = matematicos.com/servidor-go-tp1

# TP6: las herramientas de generación se ejecutan siempre versionadas con
# "go run paquete@version", sin usar binarios instalados en el sistema ni en
# $GOPATH/bin. Así el código generado no depende de la versión del host.
# templ tiene que coincidir con la versión que se agrega a go.mod (v0.3.1020)
# para que los archivos _templ.go sean compatibles con la librería del servidor.
TEMPL_VERSION = v0.3.1020
SQLC_VERSION  = v1.30.0
TEMPL = go run github.com/a-h/templ/cmd/templ@$(TEMPL_VERSION)
SQLC  = go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)

# Ciclo completo en un solo comando:
#   1) Tareas previas (prepare + build)
#   2) Ejecución del servidor en primer plano
#   3) Tareas posteriores (clean) al terminar o al cortar con Ctrl+C
#	El trap hace que el shell no se corte con Ctrl+C (o una señal TERM): espera a
#	que el servidor termine y sigue con el clean, que borra contenedores y volúmenes.
all: prepare build
	@echo ">>> Servidor en http://localhost:8080 (Ctrl+C para salir y limpiar)"
	@trap 'echo' INT TERM; ./server; $(MAKE) --no-print-directory clean

# Tareas previas: módulo y código generado, borrar contenedores y volúmenes
# viejos, levantar PostgreSQL y esperar con pg_isready a que acepte conexiones.
#	-h localhost hace que pg_isready pruebe por TCP, que recién responde cuando
#	terminó de ejecutarse el schema.sql de inicialización.
prepare: deps
	docker compose down -v
	docker compose up -d
	@echo "Esperando a que la base de datos esté lista..."
	@i=0; until docker compose exec -T db pg_isready -h localhost -U cliente -d mate_db -q; do \
		i=$$((i + 1)); \
		if [ $$i -ge 60 ]; then echo "PostgreSQL no respondio despues de 60 segundos"; exit 1; fi; \
		sleep 1; \
	done
	@echo "Base de datos lista y conectada"

# Módulo de Go y código generado. go.mod y go.sum no se suben al repositorio,
# por eso se crean acá si no existen.
#	El orden importa: primero se genera el codigo de sqlc (db/sqlc) y de templ
#	(views/*_templ.go), y recien despues go mod tidy, que necesita ver esos
#	imports para no borrar dependencias.
deps:
	@test -f go.mod || go mod init $(MODULE_NAME)
	$(SQLC) generate
	$(TEMPL) generate
	go get github.com/a-h/templ@$(TEMPL_VERSION)
	go mod tidy

# Genera solo el código Go de las queries (db/sqlc)
sqlc:
	$(SQLC) generate

# Compila solo los componentes .templ de views/ a archivos _templ.go
generate:
	$(TEMPL) generate

# Compila el binario (regenerando antes el código)
build: deps
	go build -o server .

# Compila y ejecuta el servidor sin limpiar al final, para desarrollo.
# Requiere la base levantada con "make prepare"; al terminar, "make clean".
run: build
	./server

# Tareas posteriores: detiene contenedores, borra volúmenes y el binario
clean:
	docker compose down -v
	rm -f server
