# TP 5. Renderizado del Lado del Servidor. Mate Máticos

## Descripción

Mate Máticos es una aplicación web tipo e-commerce orientada a la venta de productos relacionados con el mate, ya sea mates, termos, bombillas y materas. Este proyecto se desarrolla de forma incremental durante la cursada de Programación Web.

Este proyecto implementa la capa de persistencia utilizando Go, sqlc y PostgreSQL usando Docker Compose.

En la carpeta db se encuantra el archivo.sql que contiene las tablas y en users.sql contiene las consultas SQL.

En logic/ se encuentra la capa lógico con las validaciones del negocio

La capa web se encuentra en main.go. Desde el TP5 la aplicación usa renderizado del lado del servidor (SSR): las vistas están en views/ (componentes Templ) y la carga de datos se hace con formularios HTML.

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
* **sqlc** y **templ** (los instala el Makefile)
* **Make**
* **Hurl**
* **Agregar el usuario al grupo docker para evitar problemas de permisos**
```bash
sudo usermod -aG docker $USER
```
---

## Pasos para ejecutar el proyecto

Instrucciones para ejecutar el proyecto desde la terminal en la raíz del repositorio. El `Makefile` automatiza mediante los siguientes comandos:

### 1. Preparación del entorno y dependencias

Instala dependencias de Go, genera el código de `sqlc` y de las vistas con `templ`, levanta la base de datos PostgreSQL en Docker y aguarda con `pg_isready` a que esté lista para recibir conexiones.

```bash
make prepare

```

---

### 2. Compilación y ejecución del servidor

Compila las vistas `.templ`, genera el binario `server` y levanta la aplicación en segundo plano en el puerto `8080`, redirigiendo la salida a `server.log`.

```bash
make run

```

> Abrir la aplicación en el navegador: [http://localhost:8080](http://localhost:8080)

---

### 3. Limpieza del entorno

Detiene el servidor en segundo plano matando su PID, borra los archivos temporales de ejecución (`server`, `server.pid`, `server.log`) y elimina los contenedores y volúmenes de datos de Docker.

```bash
make clean

```

---

### Otros comandos

* **Ejecución completa del ciclo de vida** (`prepare` $\rightarrow$ `run` $\rightarrow$ `test` $\rightarrow$ `clean`):
```bash
make all

```


* **Ejecución de desarrollo sin limpieza** (`prepare` $\rightarrow$ `run`, útil para mantener el servidor corriendo y revisar `server.log`):
```bash
make execute

```


* **Generación manual de vistas** (recompila únicamente las plantillas de `views/` a código Go):
```bash
make generate

```
