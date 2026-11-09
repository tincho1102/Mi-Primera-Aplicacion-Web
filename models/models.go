package models

// CategoriaRequest define la estructura para parsear el JSON recibido para crear o actualizar categorías.
type CategoriaRequest struct {
	NombreCat      string `json:"nombre_cat"`
	DescripcionCat string `json:"descripcion_cat"`
	IDPadre        *int32 `json:"id_padre,omitempty"`
}

// ProductoRequest define la estructura para parsear el JSON recibido para crear o actualizar productos.
type ProductoRequest struct {
	NombreProd      string `json:"nombre_prod"`
	DescripcionProd string `json:"descripcion_prod"`
	Precio          string `json:"precio"`
	Stock           int32  `json:"stock"`
	IDCat           int32  `json:"id_cat"`
}
