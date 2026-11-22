// Package handlers contiene los handlers HTTP de la aplicación, separados por
// entidad (index.go, categoria.go, producto.go). Este archivo tiene la
// estructura Handler y las funciones auxiliares que comparten todos.
package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/lib/pq"
	db "matematicos.com/servidor-go-tp1/db/sqlc"
	"matematicos.com/servidor-go-tp1/views"
)

// Handler agrupa las dependencias de los handlers (patrón de la filmina 55):
// en lugar de una variable global, las queries de sqlc se inyectan una vez
// al crear el Handler y cada método las usa con h.queries.
type Handler struct {
	queries *db.Queries
}

// New crea un Handler con las queries de sqlc ya conectadas a la base.
func New(queries *db.Queries) *Handler {
	return &Handler{queries: queries}
}

// listAll trae categorías y productos, que es lo que necesitan la página
// completa y algunas respuestas parciales.
func (h *Handler) listAll(ctx context.Context) ([]db.Categoria, []db.Producto, error) {
	categorias, err := h.queries.ListCategorias(ctx)
	if err != nil {
		return nil, nil, err
	}
	productos, err := h.queries.ListProductos(ctx)
	if err != nil {
		return nil, nil, err
	}
	return categorias, productos, nil
}

// render escribe un componente templ en la respuesta con el código indicado.
func render(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := c.Render(r.Context(), w); err != nil {
		log.Println("Error al renderizar la vista:", err)
	}
}

// formError responde un error de un formulario de alta.
// El hx-target del formulario es la tabla, así que con los headers
// HX-Retarget y HX-Reswap le indicamos a HTMX que, en este caso, ponga el
// mensaje en el párrafo de error del formulario y no toque la tabla.
// Hace falta porque en HTMX 4 las respuestas 4xx/5xx también se intercambian.
func formError(w http.ResponseWriter, r *http.Request, errorID string, status int, msg string) {
	w.Header().Set("HX-Retarget", "#"+errorID)
	w.Header().Set("HX-Reswap", "outerHTML")
	render(w, r, status, views.FormError(errorID, msg, false))
}

// noSwapError responde un error con HX-Reswap: none, para que HTMX no
// reemplace el target. Se usa en los DELETE: una respuesta de error vacía o
// con texto plano borraría o rompería la fila de la tabla.
func noSwapError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("HX-Reswap", "none")
	http.Error(w, msg, status)
}

// pathID lee el parámetro {id} de rutas como DELETE /categorias/{id}
// (r.PathValue está disponible desde Go 1.22).
func pathID(r *http.Request) (int32, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("id invalido")
	}
	return int32(id), nil
}

// pgErrorCode devuelve el código SQLSTATE de un error de PostgreSQL, o ""
// si el error no viene de la base.
func pgErrorCode(err error) string {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code)
	}
	return ""
}

// Códigos SQLSTATE de PostgreSQL que se traducen a mensajes para el usuario.
const (
	pgForeignKeyViolation = "23503"
	pgUniqueViolation     = "23505"
	pgStringTooLong       = "22001"
	pgNumericOutOfRange   = "22003"
)

// dbErrorMessage traduce los errores de restricciones de PostgreSQL más
// comunes al crear una entidad a un mensaje entendible para el usuario.
// Cualquier otro error se registra en el log y se informa como error interno.
func dbErrorMessage(err error, entidad string) (int, string) {
	switch pgErrorCode(err) {
	case pgUniqueViolation:
		return http.StatusBadRequest, "Ya existe un/a " + entidad + " con ese nombre"
	case pgStringTooLong:
		return http.StatusBadRequest, "Alguno de los textos supera el largo maximo permitido"
	case pgNumericOutOfRange:
		return http.StatusBadRequest, "El precio supera el valor maximo permitido"
	case pgForeignKeyViolation:
		return http.StatusBadRequest, "La categoria asociada ya no existe"
	}
	log.Printf("Error al crear %s: %v", entidad, err)
	return http.StatusInternalServerError, "Error interno al crear " + entidad
}
