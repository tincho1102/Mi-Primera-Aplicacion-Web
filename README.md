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

HTMX 4.0.0 se carga desde la CDN en el `<head>` del layout. El proyecto **no contiene archivos `.js`**.

| Acción | Petición | Respuesta del servidor |
|--------|----------|------------------------|
| Crear categoría | `POST /categorias` (`hx-post`, `hx-target="#categorias-table-container"`, `hx-swap="outerHTML"`) | Fragmento con la tabla de categorías actualizada, más los `<select>` de categorías por `hx-swap-oob` |
| Crear producto | `POST /productos` (`hx-target="#productos-table-container"`) | Fragmento con la tabla de productos actualizada |
| Borrar producto | `DELETE /productos/{id}` (`hx-target="closest tr"`, `hx-confirm`) | `200 OK` con cuerpo vacío: la fila se elimina |
| Borrar categoría | `DELETE /categorias/{id}` (`hx-target="closest tr"`, `hx-confirm`) | `200 OK`: la fila se elimina y, por `hx-swap-oob`, se actualizan la tabla de productos (se borran en cascada) y los `<select>` |

Después de un alta exitosa el formulario se limpia con `hx-on:htmx:after:swap`.

**Manejo de errores:** en HTMX 4 las respuestas 4xx/5xx también se intercambian en el DOM, por eso:
- Si falla la validación de un alta, el servidor responde `400` con los headers `HX-Retarget` y `HX-Reswap`, y el mensaje aparece arriba del formulario sin borrar lo que el usuario escribió.
- Si se intenta borrar una categoría que tiene subcategorías, el servidor responde `409 Conflict` con la misma fila y el motivo del error, así la fila no desaparece.

---

## Datos iniciales

Al arrancar, si la tabla de categorías está vacía, la aplicación carga algunas categorías (Mates, Mates de calabaza, Termos, Bombillas, Materas) y productos con precio y stock (`db/seed/seed.go`).

---

## Requisitos previos

* **Go** 1.25 o superior (lo requiere templ v0.3.1020)
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

Un solo comando realiza todo el ciclo de vida:

```bash
make        # equivale a make all
```

1. **Tareas previas:** crea el módulo de Go si no existe, genera el código de sqlc y las vistas de templ (versionados), borra contenedores y volúmenes anteriores (`docker compose down -v`), levanta PostgreSQL (`docker compose up -d`), espera con `pg_isready` a que acepte conexiones y compila el binario `server`.
2. **Ejecución:** corre el servidor en primer plano.
3. **Tareas posteriores:** al cortar con `Ctrl+C` (o si el servidor termina) borra los contenedores, los volúmenes y el binario (`make clean`).

> Abrir la aplicación en el navegador: [http://localhost:8080](http://localhost:8080)

### Targets individuales del Makefile

| Comando | Qué hace |
|---------|----------|
| `make prepare` | Tareas previas: módulo, sqlc, templ, borrar y levantar contenedores, esperar a la base |
| `make run` | Compila y ejecuta el servidor sin limpiar al final (requiere `make prepare` antes) |
| `make clean` | Tareas posteriores: detiene contenedores, borra volúmenes y el binario |
| `make generate` | Regenera solo las vistas con `templ@v0.3.1020` |
| `make sqlc` | Regenera solo el código de sqlc |
