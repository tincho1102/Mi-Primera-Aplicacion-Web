package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	db "matematicos.com/servidor-go-tp1/db/sqlc" // Reemplazá por la ruta de tu módulo
)

var queries *db.Queries

func main() {
	// Conexión con la base de datos estándar de Go
	connStr := "postgres://cliente:secret@localhost:5432/mate_db?sslmode=disable"
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer sqlDB.Close()

	// Instancia sqlc
	queries = db.New(sqlDB)

	//rutas http
	http.HandleFunc("/categorias", categoriasHandler)
	http.HandleFunc("/categorias/", categoriasHandler)

	http.HandleFunc("/productos", productosHandler)
	http.HandleFunc("/productos/", productosHandler)

	log.Println("Servidor corriendo en el puerto :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handler

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
	categorias, err := queries.ListCategorias(context.Background())
	if err != nil {
		http.Error(w, "Error al listar categorias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categorias)
}

func createCategoria(w http.ResponseWriter, r *http.Request) {
	var newCat db.Categoria

	err := json.NewDecoder(r.Body).Decode(&newCat)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdCat, err := queries.CreateCategoria(context.Background(), db.CreateCategoriaParams{
		NombreCat:      newCat.NombreCat,
		DescripcionCat: newCat.DescripcionCat,
	})
	if err != nil {
		http.Error(w, "Error al crear categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCat)
}

func getCategoria(w http.ResponseWriter, r *http.Request, id int) {
	categoria, err := queries.GetCategoria(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

func updateCategoria(w http.ResponseWriter, r *http.Request, id int) {
	var updateCat db.Categoria

	err := json.NewDecoder(r.Body).Decode(&updateCat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cat, err := queries.GetCategoria(context.Background(), int32(id))
	if err != nil {
		http.Error(w, "Categoria no encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	err = queries.UpdateCategoria(context.Background(), db.UpdateCategoriaParams{
		IDCat:          cat.IDCat,
		NombreCat:      updateCat.NombreCat,
		DescripcionCat: updateCat.DescripcionCat,
	})
	if err != nil {
		http.Error(w, "Error al actualizar categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	updatedCat, _ := queries.GetCategoria(context.Background(), int32(id))
	json.NewEncoder(w).Encode(updatedCat)
}

func deleteCategoria(w http.ResponseWriter, r *http.Request, id int) {
	cat, err := queries.GetCategoria(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = queries.DeleteCategoria(context.Background(), cat.IDCat)
	if err != nil {
		http.Error(w, "Error al eliminar categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handler para productos

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
	productos, err := queries.ListProductos(context.Background())
	if err != nil {
		http.Error(w, "Error al listar productos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

func createProducto(w http.ResponseWriter, r *http.Request) {
	var newProd db.Producto

	err := json.NewDecoder(r.Body).Decode(&newProd)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verificar si la categoria asociada existe
	_, err = queries.GetCategoria(context.Background(), newProd.IDCat)
	if err != nil {
		http.Error(w, "La categoria asociada no existe", http.StatusBadRequest)
		return
	}

	createdProd, err := queries.CreateProducto(context.Background(), db.CreateProductoParams{
		NombreProd:      newProd.NombreProd,
		DescripcionProd: newProd.DescripcionProd,
		Precio:          newProd.Precio,
		Stock:           newProd.Stock,
		IDCat:           newProd.IDCat,
	})
	if err != nil {
		http.Error(w, "Error al crear producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProd)
}

func getProducto(w http.ResponseWriter, r *http.Request, id int) {
	producto, err := queries.GetProducto(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

func updateProducto(w http.ResponseWriter, r *http.Request, id int) {
	var updateProd db.Producto

	err := json.NewDecoder(r.Body).Decode(&updateProd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prod, err := queries.GetProducto(context.Background(), int32(id))
	if err != nil {
		http.Error(w, "Producto no encontrado: "+err.Error(), http.StatusNotFound)
		return
	}

	err = queries.UpdateProducto(context.Background(), db.UpdateProductoParams{
		IDProd:          prod.IDProd,
		NombreProd:      updateProd.NombreProd,
		DescripcionProd: updateProd.DescripcionProd,
		Precio:          updateProd.Precio,
		Stock:           updateProd.Stock,
	})
	if err != nil {
		http.Error(w, "Error al actualizar producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	updatedProd, _ := queries.GetProducto(context.Background(), int32(id))
	json.NewEncoder(w).Encode(updatedProd)
}

func deleteProducto(w http.ResponseWriter, r *http.Request, id int) {
	prod, err := queries.GetProducto(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = queries.DeleteProducto(context.Background(), prod.IDProd)
	if err != nil {
		http.Error(w, "Error al eliminar producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}