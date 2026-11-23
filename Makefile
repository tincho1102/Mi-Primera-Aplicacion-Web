.PHONY: all prepare generate run clean

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

# Ciclo completo: prepara el entorno y levanta el servidor en segundo plano
all: prepare run

# Tareas previas: módulo y código generado (generate), borrar contenedores y volúmenes
# viejos, levantar PostgreSQL, esperar con pg_isready a que acepte conexiones y compilar.
prepare: generate
	docker compose down -v
	docker compose up -d
	@echo "Esperando a que la base de datos esté lista..."
	@i=0; until docker compose exec -T db pg_isready -h localhost -U cliente -d mate_db -q; do \
		i=$$((i + 1)); \
		if [ $$i -ge 60 ]; then echo "PostgreSQL no respondio despues de 60 segundos"; exit 1; fi; \
		sleep 1; \
	done
	@echo "Base de datos lista y conectada"
	go build -o server .

# Módulo de Go y código generado por sqlc y templ.
# go.mod y go.sum no se suben al repositorio, por eso se crean acá si no existen.
generate:
	@test -f go.mod || go mod init $(MODULE_NAME)
	$(SQLC) generate
	$(TEMPL) generate
	go get github.com/a-h/templ@$(TEMPL_VERSION)
	go mod tidy

# Levanta el servidor en segundo plano para no bloquear la consola
run:
	@test -f server || go build -o server .
	@if [ -f server.pid ] && kill -0 $$(cat server.pid) 2>/dev/null; then \
		echo "El servidor ya se encuentra corriendo en segundo plano (PID: $$(cat server.pid))"; \
	else \
		./server > server.log 2>&1 & echo $$! > server.pid; \
		echo ">>> Servidor iniciado en segundo plano en http://localhost:8080 (PID: $$(cat server.pid))"; \
		echo ">>> Logs en server.log (tail -f server.log). Para detener y limpiar: make clean"; \
	fi

# Tareas posteriores: detiene el servidor, contenedores, volúmenes y archivos generados
clean:
	@if [ -f server.pid ]; then \
		echo "Deteniendo servidor (PID: $$(cat server.pid))..."; \
		kill $$(cat server.pid) 2>/dev/null || true; \
		rm -f server.pid; \
	fi
	@pkill -f "./server" 2>/dev/null || true
	docker compose down -v
	rm -f server server.log
