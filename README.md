# TP 3. Exponiendo la Lógica como API. Mate Máticos

## Descripción

Mate Máticos es una aplicación web tipo e-commerce orientada a la venta de productos relacionados con el mate, ya sea mates, termos, bombillas y materas. Este proyecto se desarrolla de forma incremental durante la cursada de Programación Web.

Este proyecto implementa la capa de persistencia utilizando Go, sqlc y PostgreSQL usando Docker Compose.

En la carpeta `db/schema/` se encuentra `schema.sql` que contiene la definición de las tablas, y en `db/queries/queries.sql` se encuentran las consultas SQL tipadas para sqlc.

En `models/` se definen los modelos y DTOs utilizados para deserializar las peticiones HTTP (desacoplados de las estructuras internas de la base de datos).

En `logic/` se encuentra la capa de lógica de negocio, encargada de las validaciones de dominio y de orquestar las consultas a la base de datos.

La capa de API Rest se encuentra en `main.go`, incluye todos los endpoints HTTP, mapea las solicitudes a los DTOs y delega a la capa lógica utilizando el contexto de la petición (`r.Context()`).

Dentro de la carpeta `static/` hay un HTML y un CSS para servir estáticamente una interfaz amigable con la documentación de los endpoints al levantar el servidor.

Por otro lado, tenemos un archivo Makefile para automatización.

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

## Requisitos previos

* **Go** (v1.20 o superior)
* **Docker** y **Docker Compose**
* **sqlc**
* **Make**
* **Hurl**
* **Agregar el usuario al grupo docker para evitar problemas de permisos**
```bash
sudo usermod -aG docker $USER
```
---

## Pasos para ejecutar el proyecto

Instrucciones en la terminal desde la raíz del proyecto.

En el archivo Makefile se hicieron los comandos que se van a utilizar para agilizar la instalación y ejecucion de tests

### 1. Instalacion de módulo, dependencias, sqlc y Docker. Levanta la base y se asegura de que esté lista.
```bash
    make prepare
```
---
### 2. Levanta el servidor en segundo plano. Puerto 8080
```bash
    make run
```
---
### 3. Ejecución de pruebas automatizadas de los Endpoints con hurl
```bash
    make test
```
---
### 3. Limpieza de contenedores y volúmenes, borra lo que tiene las tablas. Y mata el proceso del servidor en segundo plano
```bash
    make clean
```
---
### Adicional. Ejecuta todo en orden. Hace prepare run test clean
```bash
    make all
```
---
