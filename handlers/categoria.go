package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	db "matematicos.com/servidor-go-tp1/db/sqlc"
	"matematicos.com/servidor-go-tp1/logic"
	"matematicos.com/servidor-go-tp1/views"
)

// CreateCategoria atiende POST /categorias (hx-post del formulario).
// TP6: ya no redirige con 303 (PRG del TP5). Inserta la categoría, vuelve a
// consultar la lista y responde solo el fragmento con la tabla, que HTMX
// coloca en #categorias-table-container. La respuesta también actualiza por
// OOB los <select> de categorías de ambos formularios.
func (h *Handler) CreateCategoria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// cualquier error se muestra en el párrafo de error del formulario
	fail := func(status int, msg string) {
		formError(w, r, views.CategoriaFormErrorID, status, msg)
	}

	if err := r.ParseForm(); err != nil {
		fail(http.StatusBadRequest, "Formulario invalido")
		return
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
		if _, err := h.queries.GetCategoria(ctx, newCat.IDPadre.Int32); err != nil {
			fail(http.StatusBadRequest, "La categoria padre no existe")
			return
		}
	}

	_, err := h.queries.CreateCategoria(ctx, db.CreateCategoriaParams{
		NombreCat:      newCat.NombreCat,
		DescripcionCat: newCat.DescripcionCat,
		IDPadre:        newCat.IDPadre,
	})
	if err != nil {
		fail(dbErrorMessage(err, "categoria"))
		return
	}

	categorias, err := h.queries.ListCategorias(ctx)
	if err != nil {
		log.Println("Error al listar categorias:", err)
		fail(http.StatusInternalServerError, "La categoria se creo, pero no se pudo actualizar la lista. Recargue la pagina.")
		return
	}

	render(w, r, http.StatusOK, views.CategoriaCreada(categorias))
}

// DeleteCategoria atiende DELETE /categorias/{id} (botón Borrar de cada fila).
// Si se borra bien responde 200 sin contenido principal, y HTMX reemplaza la
// fila (hx-target="closest tr") por nada. Como producto.id_cat tiene
// ON DELETE CASCADE, la respuesta trae por OOB la tabla de productos y los
// <select> actualizados.
func (h *Handler) DeleteCategoria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := pathID(r)
	if err != nil {
		noSwapError(w, http.StatusBadRequest, "ID de categoria invalido")
		return
	}

	if err := h.queries.DeleteCategoria(ctx, id); err != nil {
		// categoria.id_padre no tiene cascada: si la categoría tiene
		// subcategorías, PostgreSQL rechaza el borrado.
		if pgErrorCode(err) == pgForeignKeyViolation {
			h.categoriaRowError(w, r, id, "No se puede borrar: tiene subcategorías. Borralas primero.")
			return
		}
		log.Printf("Error al borrar la categoria %d: %v", id, err)
		noSwapError(w, http.StatusInternalServerError, "Error interno al borrar la categoria")
		return
	}

	categorias, productos, err := h.listAll(ctx)
	if err != nil {
		// la categoría ya se borró: respondemos vacío para que al menos
		// desaparezca la fila, aunque el resto de la página quede sin refrescar
		log.Println("Error al listar despues de borrar la categoria:", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	render(w, r, http.StatusOK, views.CategoriaBorrada(categorias, productos))
}

// categoriaRowError responde 409 Conflict con la misma fila de la categoría y
// el mensaje de error. En HTMX 4 las respuestas 4xx se intercambian: si
// respondiéramos un 409 vacío, la fila desaparecería igual que en un borrado
// exitoso. Devolviendo la fila, se queda en la tabla mostrando el motivo.
func (h *Handler) categoriaRowError(w http.ResponseWriter, r *http.Request, id int32, msg string) {
	ctx := r.Context()
	cat, err := h.queries.GetCategoria(ctx, id)
	if err != nil {
		noSwapError(w, http.StatusConflict, msg)
		return
	}
	// la lista completa hace falta para mostrar el nombre de la categoría padre
	categorias, err := h.queries.ListCategorias(ctx)
	if err != nil {
		noSwapError(w, http.StatusConflict, msg)
		return
	}
	render(w, r, http.StatusConflict, views.CategoriaRow(cat, categorias, msg))
}
