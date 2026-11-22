package views

import "github.com/a-h/templ"

// IDs de los párrafos de error de cada formulario. Los handlers los usan en el
// header HX-Retarget para que HTMX muestre ahí el mensaje de validación, en
// lugar de reemplazar la tabla (que es el hx-target normal del formulario).
const (
	CategoriaFormErrorID = "categoria-form-error"
	ProductoFormErrorID  = "producto-form-error"
)

// oobAttrs agrega hx-swap-oob="true" cuando un componente viaja como fragmento
// "Out of Band" (filminas 48-49): HTMX busca en la página el elemento con el
// mismo id y lo reemplaza, aparte del swap principal de la respuesta.
// En una carga normal de la página devuelve nil y no agrega nada.
func oobAttrs(oob bool) templ.Attributes {
	if oob {
		return templ.Attributes{"hx-swap-oob": "true"}
	}
	return nil
}
