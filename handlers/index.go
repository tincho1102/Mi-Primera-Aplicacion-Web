package handlers

import (
	"net/http"

	"matematicos.com/servidor-go-tp1/views"
)

// Index atiende GET / : es la única respuesta con la página completa
// (layout, formularios y tablas). El resto de las interacciones se resuelven
// con HTMX intercambiando fragmentos.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	categorias, productos, err := h.listAll(r.Context())
	if err != nil {
		http.Error(w, "Error al listar categorias y productos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, r, http.StatusOK, views.Index(categorias, productos))
}
