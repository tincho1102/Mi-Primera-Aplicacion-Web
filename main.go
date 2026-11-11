package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
	db "matematicos.com/servidor-go-tp1/db/sqlc"
	"matematicos.com/servidor-go-tp1/logic"
	"matematicos.com/servidor-go-tp1/views"
)

var queries *db.Queries

func main() {
	// conexion con la base de datos de Go
	//consultamos con la IA y nos mostro que sql.open no verifica la conexion,
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

	// instancia sqlc
	queries = db.New(sqlDB)

	// rutas http (patrones con método, disponibles desde Go 1.22)
	// "GET /{$}" atiende solo la raíz exacta; cualquier otra ruta devuelve 404.
	http.HandleFunc("GET /{$}", indexHandler)
	http.HandleFunc("POST /categorias", createCategoriaHandler)
	http.HandleFunc("POST /productos", createProductoHandler)

	log.Println("Servidor corriendo en http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// renderIndex consulta categorías y productos con sqlc y renderiza la página
// principal con Templ, escribiendo el HTML directamente en el ResponseWriter.
// status permite devolver 200 en una carga normal o 400 cuando se vuelve a
// mostrar la página por un error de validación de algún formulario.
func renderIndex(w http.ResponseWriter, r *http.Request, status int, catForm, prodForm views.FormState) {
	categorias, err := queries.ListCategorias(r.Context())
	if err != nil {
		http.Error(w, "Error al listar categorias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	productos, err := queries.ListProductos(r.Context())
	if err != nil {
		http.Error(w, "Error al listar productos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	component := views.Index(categorias, productos, catForm, prodForm)
	if err := component.Render(r.Context(), w); err != nil {
		log.Println("Error al renderizar la vista:", err)
	}
}

// GET / : página principal con formularios y listados.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderIndex(w, r, http.StatusOK, views.FormState{}, views.FormState{})
}

// POST /categorias : alta de categoría desde el formulario (patrón PRG).
func createCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario invalido", http.StatusBadRequest)
		return
	}

	// si algo falla, se vuelve a mostrar la página con el mensaje de error
	// y los datos que el usuario había cargado
	fail := func(status int, msg string) {
		renderIndex(w, r, status, views.FormState{Error: msg, Values: r.PostForm}, views.FormState{})
	}

	newCat := db.Categoria{
		NombreCat:      strings.TrimSpace(r.PostForm.Get("nombre_cat")),
		DescripcionCat: strings.TrimSpace(r.PostForm.Get("descripcion_cat")),
	}

	// id_padre es opcional: vacío significa NULL (sin categoría padre)
	if padre := r.PostForm.Get("id_padre"); padre != "" {
		id, err := strconv.ParseInt(padre, 10, 32)
		if err != nil {
			fail(http.StatusBadRequest, "La categoria padre no es valida")
			return
		}
		newCat.IDPadre = sql.NullInt32{Int32: int32(id), Valid: true}
	}

	// validaciones de negocio de la capa logic
	if err := logic.ValidateCategoria(newCat); err != nil {
		fail(http.StatusBadRequest, err.Error())
		return
	}

	// si tiene padre, verificamos que exista
	if newCat.IDPadre.Valid {
		if _, err := queries.GetCategoria(r.Context(), newCat.IDPadre.Int32); err != nil {
			fail(http.StatusBadRequest, "La categoria padre no existe")
			return
		}
	}

	_, err := queries.CreateCategoria(r.Context(), db.CreateCategoriaParams{
		NombreCat:      newCat.NombreCat,
		DescripcionCat: newCat.DescripcionCat,
		IDPadre:        newCat.IDPadre,
	})
	if err != nil {
		status, msg := dbErrorMessage(err, "categoria")
		fail(status, msg)
		return
	}

	// PRG: redirigimos con 303 para que el navegador haga un GET a /
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /productos : alta de producto desde el formulario (patrón PRG).
func createProductoHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario invalido", http.StatusBadRequest)
		return
	}

	fail := func(status int, msg string) {
		renderIndex(w, r, status, views.FormState{}, views.FormState{Error: msg, Values: r.PostForm})
	}

	// los campos numéricos llegan como texto desde el formulario
	stock, err := strconv.ParseInt(r.PostForm.Get("stock"), 10, 32)
	if err != nil {
		fail(http.StatusBadRequest, "El stock debe ser un numero entero")
		return
	}
	idCat, err := strconv.ParseInt(r.PostForm.Get("id_cat"), 10, 32)
	if err != nil {
		fail(http.StatusBadRequest, "Debe elegir una categoria")
		return
	}

	newProd := db.Producto{
		NombreProd:      strings.TrimSpace(r.PostForm.Get("nombre_prod")),
		DescripcionProd: strings.TrimSpace(r.PostForm.Get("descripcion_prod")),
		Precio:          strings.TrimSpace(r.PostForm.Get("precio")), // logic se encarga de validar que sea numérico
		Stock:           int32(stock),
		IDCat:           int32(idCat),
	}

	// validaciones de negocio de la capa logic
	if err := logic.ValidateProducto(newProd); err != nil {
		fail(http.StatusBadRequest, err.Error())
		return
	}

	// Verificar si la categoria asociada existe
	if _, err := queries.GetCategoria(r.Context(), newProd.IDCat); err != nil {
		fail(http.StatusBadRequest, "La categoria asociada no existe")
		return
	}

	_, err = queries.CreateProducto(r.Context(), db.CreateProductoParams{
		NombreProd:      newProd.NombreProd,
		DescripcionProd: newProd.DescripcionProd,
		Precio:          newProd.Precio,
		Stock:           newProd.Stock,
		IDCat:           newProd.IDCat,
	})
	if err != nil {
		status, msg := dbErrorMessage(err, "producto")
		fail(status, msg)
		return
	}

	// PRG: redirigimos con 303 para que el navegador haga un GET a /
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// dbErrorMessage traduce los errores de restricciones de PostgreSQL más
// comunes a un mensaje entendible para el usuario. Cualquier otro error se
// registra en el log y se informa como error interno.
func dbErrorMessage(err error, entidad string) (int, string) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			return http.StatusBadRequest, "Ya existe un/a " + entidad + " con ese nombre"
		case "22001": // string_data_right_truncation
			return http.StatusBadRequest, "Alguno de los textos supera el largo maximo permitido"
		case "22003": // numeric_value_out_of_range
			return http.StatusBadRequest, "El precio supera el valor maximo permitido"
		}
	}
	log.Printf("Error al crear %s: %v", entidad, err)
	return http.StatusInternalServerError, "Error interno al crear " + entidad
}
