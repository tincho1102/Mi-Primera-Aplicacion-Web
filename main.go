package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	db "matematicos.com/servidor-go-tp1/db/sqlc"
	"matematicos.com/servidor-go-tp1/logic"
	"matematicos.com/servidor-go-tp1/models"
)

var (
	catLogic  *logic.CategoriaLogic
	prodLogic *logic.ProductoLogic
)

func main() {
	// conexion con la base de datos de Go
	// consultamos con la IA y nos mostro que sql.open no verifica la conexion, 
	// por lo que debemos hacer un ping para verificar la conexion
	connStr := "postgres://cliente:progwebunicen@localhost:5432/mate_db?sslmode=disable"	
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("La base no responde al ping:", err)
	}
	defer sqlDB.Close()

	// instancia de sqlc e inyeccion a la capa de logica
	queries := db.New(sqlDB)
	catLogic = logic.NewCategoriaLogic(queries)
	prodLogic = logic.NewProductoLogic(queries)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Si entran exactamente a "/", mostramos la landing bonita
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./static/index.html")
			return
		}
		// Para cualquier otra cosa dentro de static/ (como css o js si tuvieras), lo servimos:
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
	})

	// rutas http
	http.HandleFunc("/categorias", categoriasHandler)
	http.HandleFunc("/categorias/", categoriasHandler)

	http.HandleFunc("/productos", productosHandler)
	http.HandleFunc("/productos/", productosHandler)

	log.Println("Servidor corriendo en el puerto :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handlers para categorías

func categoriasHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) == 1 && parts[0] == "categorias" {
		switch r.Method {
		case http.MethodGet:
			getCategorias(w, r)
		case http.MethodPost:
			createCategoria(w, r)
		default:
			http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 2 && parts[0] == "categorias" {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "ID de categoria invalido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			getCategoria(w, r, id)
		case http.MethodPut:
			updateCategoria(w, r, id)
		case http.MethodDelete:
			deleteCategoria(w, r, id)
		default:
			http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	http.Error(w, "Invalid URL", http.StatusBadRequest)
}

func getCategorias(w http.ResponseWriter, r *http.Request) {
	categorias, err := catLogic.ListCategorias(r.Context())
	if err != nil {
		http.Error(w, "Error al listar categorias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categorias)
}

func createCategoria(w http.ResponseWriter, r *http.Request) {
	var req models.CategoriaRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdCat, err := catLogic.CreateCategoria(r.Context(), req)
	if err != nil {
		var valErr *logic.ValidationError
		if errors.As(err, &valErr) {
			http.Error(w, valErr.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Error al crear categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCat)
}

func getCategoria(w http.ResponseWriter, r *http.Request, id int) {
	categoria, err := catLogic.GetCategoria(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

func updateCategoria(w http.ResponseWriter, r *http.Request, id int) {
	var req models.CategoriaRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedCat, err := catLogic.UpdateCategoria(r.Context(), int32(id), req)
	if err != nil {
		var valErr *logic.ValidationError
		if errors.As(err, &valErr) {
			http.Error(w, valErr.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Categoria no encontrada: "+err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Error al actualizar categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCat)
}

func deleteCategoria(w http.ResponseWriter, r *http.Request, id int) {
	err := catLogic.DeleteCategoria(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handlers para productos

func productosHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) == 1 && parts[0] == "productos" {
		switch r.Method {
		case http.MethodGet:
			getProductos(w, r)
		case http.MethodPost:
			createProducto(w, r)
		default:
			http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 2 && parts[0] == "productos" {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "ID de producto invalido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			getProducto(w, r, id)
		case http.MethodPut:
			updateProducto(w, r, id)
		case http.MethodDelete:
			deleteProducto(w, r, id)
		default:
			http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	http.Error(w, "Invalid URL", http.StatusBadRequest)
}

func getProductos(w http.ResponseWriter, r *http.Request) {
	productos, err := prodLogic.ListProductos(r.Context())
	if err != nil {
		http.Error(w, "Error al listar productos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

func createProducto(w http.ResponseWriter, r *http.Request) {
	var req models.ProductoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdProd, err := prodLogic.CreateProducto(r.Context(), req)
	if err != nil {
		var valErr *logic.ValidationError
		if errors.As(err, &valErr) {
			http.Error(w, valErr.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Error al crear producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProd)
}

func getProducto(w http.ResponseWriter, r *http.Request, id int) {
	producto, err := prodLogic.GetProducto(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

func updateProducto(w http.ResponseWriter, r *http.Request, id int) {
	var req models.ProductoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedProd, err := prodLogic.UpdateProducto(r.Context(), int32(id), req)
	if err != nil {
		var valErr *logic.ValidationError
		if errors.As(err, &valErr) {
			http.Error(w, valErr.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Producto no encontrado: "+err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Error al actualizar producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProd)
}

func deleteProducto(w http.ResponseWriter, r *http.Request, id int) {
	err := prodLogic.DeleteProducto(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
