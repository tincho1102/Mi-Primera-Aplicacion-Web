package logic

import (
	"errors"
	"strconv"
	db "tpe/db/sqlc"
)

func ValidateProducto(p db.Producto) error {
	if p.NombreProd == "" {
		return errors.New("El nombre del producto no puede estar vacio")
	}
	if p.DescripcionProd == "" {
		return errors.New("La descripcion del producto no puede estar vacia")
	}
	precioFloat, err := strconv.ParseFloat(p.Precio, 64)
	if err != nil {
		return errors.New("El precio del producto no es un valor valido")
	}
	if precioFloat <= 0 {
		return errors.New("El precio del producto debe ser mayor a cero")
	}
	if p.Stock < 0 {
		return errors.New("El stock del producto debe ser mayor o igual a cero")
	}
	if p.IDCat <= 0 {
		return errors.New("El id de la categoria a la que pertenece el producto debe ser mayor a cero")
	}

	return nil
}