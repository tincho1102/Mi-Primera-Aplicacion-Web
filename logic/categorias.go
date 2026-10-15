package logic

import (
	"errors"
	db "matematicos.com/servidor-go-tp1/db/sqlc"
)

func ValidateCategoria(p db.Categoria) error {
	if p.NombreCat == "" {
		return errors.New("El nombre de la categoria no puede estar vacio")
	}
	if p.DescripcionCat == "" {
		return errors.New("La descripcion de la categoria no puede estar vacia")
	}
	if p.IDPadre.Valid && p.IDPadre.Int32 <= 0 {
		return errors.New("El id de la categoria padre debe ser mayor a cero o nulo")
	}//preguntar

	return nil
}