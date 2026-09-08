# TP 2. Capa de persistencia. Mate Máticos

## Descripción

Mate Máticos es una aplicación web tipo e-commerce orientada a la venta de productos relacionados con el mate, ya sea mates, termos, bombillas y materas. Este proyecto se desarrolla de forma incremental durante la cursada de Programación Web.

Este proyecto implementa la capa de persistencia utilizando Go, sqlc y PostgreSQL usando Docker Compose.

En la carpeta db se encuantra el archivo.sql que contiene las tablas y en users.sql contiene las consultas SQL.

En el main.go se encuentran las pruebas CRUD con las operaciones a realizar sobre las tablas.



## Dominio

Cada producto tendrá asociado un identificador, su nombre, descripción, categoría (mates, termos y accesorios), precio y stock.

Ademas hay una tabla con categorias, hecha para que los productos pueden tener subcategorias.

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

## Requisitos previos

* **Go** (v1.20 o superior)
* **Docker** y **Docker Compose**
* **sqlc**

---

## Pasos para ejecutar el proyecto

Instrucciones en la terminal desde la raíz del proyecto:

### 1. Inicializar el módulo de Go (si no está creado)
```bash
go mod init [MateMaticos.com/servidor-go-tp1](https://MateMaticos.com/servidor-go-tp1)
go get [github.com/jackc/pgx/v5](https://github.com/jackc/pgx/v5)
go mod tidy
```
---
### 2. Instalar sqlc (opcional, si no se encuentra instalado)
```bash
go install [github.com/sqlc-dev/sqlc/cmd/sqlc@latest](https://github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
```
---
### 3. Generar el código Go con sqlc
```bash
sqlc generate
```
---
### 4. Levantar la base de datos en Docker
```bash
docker compose up -d
```
---
### 5. Ejecutar la aplicación
```bash
go run main.go
```
---
### OPCIONAL-> Limpiar el entorno
Para detener el contenedor y limpiar los volúmenes de la base de datos:
```bash
docker compose down -v
```
---