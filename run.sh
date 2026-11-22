#!/usr/bin/env bash
# Ciclo de vida completo de la aplicación (TP6):
#   1) Tareas previas: borrar contenedores y volúmenes, generar sqlc y templ
#      (versionados), compilar, levantar PostgreSQL y esperar a que acepte
#      conexiones.
#   2) Ejecución: correr el servidor en primer plano.
#   3) Tareas posteriores: al terminar o con Ctrl+C se borran contenedores,
#      volúmenes y el binario (trap EXIT).
set -euo pipefail

# Versiones fijas de las herramientas de generación: se ejecutan con
# "go run paquete@version", así no dependen de lo que haya instalado en el
# sistema ni en $GOPATH/bin.
TEMPL_VERSION="v0.3.1020"
SQLC_VERSION="v1.30.0"

# Trabajar siempre desde la carpeta del proyecto, aunque se llame desde otro lado
cd "$(dirname "$0")"

cleanup() {
    echo ""
    echo ">>> [Tareas posteriores] Deteniendo y borrando contenedores y volúmenes..."
    docker compose down -v || true
    rm -f server
}
# EXIT cubre la salida normal, los errores (set -e) y las señales; INT y TERM
# solo fuerzan la salida para que cleanup se ejecute una única vez.
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

echo ">>> [Tareas previas] Limpiando contenedores y volúmenes previos..."
docker compose down -v

echo ">>> [Tareas previas] Generando código con sqlc ${SQLC_VERSION}..."
go run "github.com/sqlc-dev/sqlc/cmd/sqlc@${SQLC_VERSION}" generate

echo ">>> [Tareas previas] Generando vistas con templ ${TEMPL_VERSION}..."
go run "github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}" generate

echo ">>> [Tareas previas] Compilando la aplicación..."
go build -o server .

echo ">>> [Tareas previas] Levantando PostgreSQL con Docker Compose..."
docker compose up -d

echo ">>> [Tareas previas] Esperando a que PostgreSQL esté listo..."
intentos=0
until docker compose exec -T db pg_isready -h localhost -U cliente -d mate_db -q; do
    intentos=$((intentos + 1))
    if [ "$intentos" -ge 60 ]; then
        echo ">>> ERROR: PostgreSQL no respondió después de 60 segundos." >&2
        exit 1
    fi
    sleep 1
done
echo ">>> Base de datos lista y conectada."

echo ">>> [Ejecución] Servidor en http://localhost:8080 (Ctrl+C para salir y limpiar)..."
./server
