# TP 6. Interfaces Reactivas con HTMX. Mate Máticos

## Descripción

Mate Máticos es una aplicación web tipo e-commerce orientada a la venta de productos relacionados con el mate, ya sea mates, termos, bombillas y materas. Este proyecto se desarrolla de forma incremental durante la cursada de Programación Web.

Este proyecto implementa la capa de persistencia utilizando Go, sqlc y PostgreSQL usando Docker Compose. Las vistas se renderizan en el servidor con Templ y, desde el TP6, la interfaz es reactiva con **HTMX**: las altas y bajas de categorías y productos se hacen sin recargar la página.

---

## Estructura del proyecto

| Ruta | Contenido |
|------|-----------|
| `main.go` | Conexión a la base, carga de datos iniciales y registro de rutas |
| `handlers/` | Handlers HTTP separados por entidad (`index.go`, `categoria.go`, `producto.go`) y la estructura `Handler` con las dependencias (`common.go`) |
| `db/schema/` y `db/queries/` | Tablas y consultas SQL que usa sqlc |
| `db/seed/` | Datos iniciales (categorías y productos) que se cargan si la base está vacía |
| `logic/` | Validaciones de negocio |
| `views/` | Componentes Templ (layout, formularios, tablas y fragmentos de respuesta HTMX) |
| `Makefile` | Automatización del ciclo de vida completo |

---

## Dominio

Cada producto tendrá asociado un identificador, su nombre, descripción, categoría (mates, termos y accesorios), precio y stock.

Ademas hay una tabla con categorias, hecha para que los productos puedan tener subcategorias.

---

# Modelo de organización
El modelo de organización de Mate Máticos se puede ver en 2 tablas:

### producto
    - id_prod SERIAL PRIMARY KEY
    - nombre_prod VARCHAR(30) UNIQUE NOT NULL
    - descripcion_prod VARCHAR(255) NOT NULL
    - precio DECIMAL (10, 2) NOT NULL CHECK (precio >= 0)
    - stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0)
    - id_cat INT NOT NULL REFERENCES categoria(id_cat) ON DELETE CASCADE
    - created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

### categoria
    - id_cat SERIAL PRIMARY KEY
    - nombre_cat VARCHAR(20) UNIQUE NOT NULL
    - descripcion_cat VARCHAR(255) NOT NULL
    - id_padre INT NULL REFERENCES categoria(id_cat)
    - created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

---

## Interacción con HTMX

HTMX se carga desde la CDN en el `<head>` del layout. El proyecto **no contiene archivos `.js`**.

| Acción | Petición | Respuesta del servidor |
|--------|----------|------------------------|
| Crear categoría | `POST /categorias` (`hx-post`, `hx-target="#categorias-table-container"`, `hx-swap="outerHTML"`) | Fragmento con la tabla de categorías actualizada, más los `<select>` de categorías por `hx-swap-oob` |
| Crear producto | `POST /productos` (`hx-target="#productos-table-container"`) | Fragmento con la tabla de productos actualizada |
| Borrar producto | `DELETE /productos/{id}` (`hx-target="closest tr"`, `hx-confirm`) | `200 OK` con cuerpo vacío: la fila se elimina |
| Borrar categoría | `DELETE /categorias/{id}` (`hx-target="closest tr"`, `hx-confirm`) | `200 OK`: la fila se elimina y, por `hx-swap-oob`, se actualizan la tabla de productos (se borran en cascada) y los `<select>` |

Después de un alta exitosa el formulario se limpia con `hx-on:htmx:after:swap`.

Usamos HTMX porque permite que una aplicacion web se sienta mas fluida, sin js.

Las ventajas de usar HTMX son: 
 * Evita recargar la pagina completa
 * Mejora el HTML estandar
 * Mantiene la logica en el servidor

---

## Datos iniciales

Al arrancar, si la tabla de categorías está vacía, la aplicación carga algunas categorías (Mates, Mates de calabaza, Termos, Bombillas, Materas) y productos con precio y stock (`db/seed/seed.go`).

---

## Requisitos previos

* **Go** 1.22 o superior
* **Docker** y **Docker Compose**
* **Make**
* **Agregar el usuario al grupo docker para evitar problemas de permisos**
```bash
sudo usermod -aG docker $USER
```

No hace falta instalar **sqlc** ni **templ**: se ejecutan con versiones fijas mediante `go run` (`sqlc@v1.30.0` y `templ@v0.3.1020`), sin depender de lo que esté instalado en el sistema. Tampoco hace falta crear el módulo: `go.mod` y `go.sum` no se suben al repositorio y el Makefile los genera (`go mod init`, `go get github.com/a-h/templ@v0.3.1020` y `go mod tidy`).

---

## Pasos para ejecutar el proyecto
Para clonar el proyecto se usa git clone -b tp6 https://github.com/tincho1102/Mi-Primera-Aplicacion-Web.git
Esto se hace para descargar la rama correcta que en este caso es la del tp6.

Luego hay que entrar a la carpeta usando:
    cd Mi-Primera-Aplicacion-Web

### Ejecución completa (recomendada)

Un solo comando realiza la preparación y levanta el servidor en segundo plano:

```bash
make        # equivale a make all
```

1. **Tareas previas:** crea el módulo de Go si no existe, genera el código de sqlc y las vistas de templ (versionados), borra contenedores y volúmenes anteriores (`docker compose down -v`), levanta PostgreSQL (`docker compose up -d`), espera con `pg_isready` a que acepte conexiones y compila el binario `server`.
2. **Ejecución:** inicia el servidor en segundo plano (`server.pid`) y redirige su salida a `server.log`, dejando la terminal libre para seguir introduciendo comandos.
3. **Tareas posteriores:** cuando desees detener la aplicación y limpiar el entorno, ejecutás:
```bash
make clean
```
Esto detiene el proceso del servidor, elimina contenedores y volúmenes de Docker, y borra los archivos temporales.

> Abrir la aplicación en el navegador: [http://localhost:8080](http://localhost:8080)  
> Para ver los logs del servidor en tiempo real: `tail -f server.log`

### Comandos del Makefile

| Comando | Qué hace |
|---------|----------|
| `make` o `make all` | **Ciclo completo:** ejecuta tareas previas y levanta el servidor en segundo plano, liberando la consola |
| `make prepare` | **Tareas previas:** genera código (`sqlc` y `templ`), levanta PostgreSQL con Docker, espera conexión y compila el servidor |
| `make run` | **Ejecución:** levanta el servidor en segundo plano (requiere la base levantada con `make prepare`) |
| `make generate` | **Generación de código:** crea el módulo si no existe y regenera el código de `sqlc` y las vistas de `templ` |
| `make clean` | **Tareas posteriores:** detiene el proceso del servidor en segundo plano, borra volúmenes/contenedores de Docker y los archivos temporales (`server`, `server.pid`, `server.log`) |
