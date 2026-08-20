# Mate Máticos

## Descripción

Mate Máticos es una aplicación web tipo e-commerce orientada a la venta de productos relacionados con el mate, como mates, termos, bombillas y materas. Este proyecto se desarrolla de forma incremental durante la cursada de Programación Web.

## Dominio

Cada producto tendrá asociado un identificador, su nombre, descripción, categoría (mates, termos y accesorios), precio y stock.

## Tecnologías utilizadas

- Go
- HTML5

## Estructura del proyecto

```
servidor-go-tp1/
├── main.go
├── go.mod
└── static/
    └── index.html
```


## Cómo ejecutar el proyecto

1. Descargar o clonar el proyecto.
2. Abrir una terminal en la carpeta del proyecto(servidor-go-tp1).
3. Ejecutar:

   ```bash
   go run main.go
   ```

4. Abrir el navegador en:

   ```
   http://localhost:8080
   ```

## Estado del proyecto

Primera entrega: servidor web básico que sirve la página estatica de presentación (`index.html`) desde el puerto 8080.
