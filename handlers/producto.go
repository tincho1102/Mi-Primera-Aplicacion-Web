package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	db "matematicos.com/servidor-go-tp1/db/sqlc"
	"matematicos.com/servidor-go-tp1/logic"
	"matematicos.com/servidor-go-tp1/views"
)

// CreateProducto atiende POST /productos (hx-post del formulario).
// TP6: ya no redirige con 303. Inserta el producto, vuelve a consultar
// productos y categorías y responde solo el fragmento con la tabla, que HTMX
// coloca en #productos-table-container.
func (h *Handler) CreateProducto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// cualquier error se muestra en el párrafo de error del formulario
	fail := func(status int, msg string) {
		formError(w, r, views.ProductoFormErrorID, status, msg)
	}

	if err := r.ParseForm(); err != nil {
		fail(http.StatusBadRequest, "Formulario invalido")
		return
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
	if _, err := h.queries.GetCategoria(ctx, newProd.IDCat); err != nil {
		fail(http.StatusBadRequest, "La categoria asociada no existe")
		return
	}

	_, err = h.queries.CreateProducto(ctx, db.CreateProductoParams{
		NombreProd:      newProd.NombreProd,
		DescripcionProd: newProd.DescripcionProd,
		Precio:          newProd.Precio,
		Stock:           newProd.Stock,
		IDCat:           newProd.IDCat,
	})
	if err != nil {
		fail(dbErrorMessage(err, "producto"))
		return
	}

	categorias, productos, err := h.listAll(ctx)
	if err != nil {
		log.Println("Error al listar productos:", err)
		fail(http.StatusInternalServerError, "El producto se creo, pero no se pudo actualizar la lista. Recargue la pagina.")
		return
	}

	render(w, r, http.StatusOK, views.ProductoCreado(productos, categorias))
}

// DeleteProducto atiende DELETE /productos/{id} (botón Borrar de cada fila).
// Responde 200 OK con cuerpo vacío: con hx-swap="outerHTML", HTMX reemplaza
// la fila (hx-target="closest tr") por nada y el producto desaparece de la
// vista sin recargar la página.
func (h *Handler) DeleteProducto(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		noSwapError(w, http.StatusBadRequest, "ID de producto invalido")
		return
	}

	if err := h.queries.DeleteProducto(r.Context(), id); err != nil {
		log.Printf("Error al borrar el producto %d: %v", id, err)
		noSwapError(w, http.StatusInternalServerError, "Error interno al borrar el producto")
		return
	}

	w.WriteHeader(http.StatusOK)
}
