package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq" // driver de PostgreSQL para database/sql
	"matematicos.com/servidor-go-tp1/db/seed"
	db "matematicos.com/servidor-go-tp1/db/sqlc"
	"matematicos.com/servidor-go-tp1/handlers"
)

func main() {
	// conexion con la base de datos de Go
	//consultamos con la IA y nos mostro que sql.open no verifica la conexion,
	// por lo que debemos hacer un ping para verificar la conexion
	connStr := "postgres://cliente:progwebunicen@localhost:5432/mate_db?sslmode=disable"
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("La base no responde al ping:", err)
	}

	// TP6: datos iniciales para no tener que cargarlos a mano
	if err := seed.Run(context.Background(), sqlDB); err != nil {
		log.Fatal("Error al cargar los datos iniciales:", err)
	}

	// los handlers reciben las queries de sqlc por inyección (handlers.Handler)
	h := handlers.New(db.New(sqlDB))

	// rutas http (patrones con método y comodines, disponibles desde Go 1.22)
	// "GET /{$}" atiende solo la raíz exacta; cualquier otra ruta devuelve 404.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("POST /categorias", h.CreateCategoria)
	mux.HandleFunc("DELETE /categorias/{id}", h.DeleteCategoria)
	mux.HandleFunc("POST /productos", h.CreateProducto)
	mux.HandleFunc("DELETE /productos/{id}", h.DeleteProducto)

	log.Println("Servidor corriendo en http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
