package logic

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	db "matematicos.com/servidor-go-tp1/db/sqlc"
	"matematicos.com/servidor-go-tp1/models"
)

type ProductoLogic struct {
	queries *db.Queries
}

func NewProductoLogic(q *db.Queries) *ProductoLogic {
	return &ProductoLogic{queries: q}
}

func ValidateProducto(p models.ProductoRequest) error {
	if p.NombreProd == "" {
		return NewValidationError("El nombre del producto no puede estar vacio")
	}
	if p.DescripcionProd == "" {
		return NewValidationError("La descripcion del producto no puede estar vacia")
	}
	precioFloat, err := strconv.ParseFloat(p.Precio, 64)
	if err != nil {
		return NewValidationError("El precio del producto no es un valor valido")
	}
	if precioFloat <= 0 {
		return NewValidationError("El precio del producto debe ser mayor a cero")
	}
	if p.Stock < 0 {
		return NewValidationError("El stock del producto debe ser mayor o igual a cero")
	}
	if p.IDCat <= 0 {
		return NewValidationError("El id de la categoria a la que pertenece el producto debe ser mayor a cero")
	}

	return nil
}

func (l *ProductoLogic) ListProductos(ctx context.Context) ([]db.Producto, error) {
	return l.queries.ListProductos(ctx)
}

func (l *ProductoLogic) GetProducto(ctx context.Context, id int32) (db.Producto, error) {
	return l.queries.GetProducto(ctx, id)
}

func (l *ProductoLogic) CreateProducto(ctx context.Context, req models.ProductoRequest) (db.Producto, error) {
	if err := ValidateProducto(req); err != nil {
		return db.Producto{}, err
	}

	// Verificar si la categoria asociada existe (regla de negocio y persistencia en logic)
	_, err := l.queries.GetCategoria(ctx, req.IDCat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Producto{}, NewValidationError("La categoria asociada no existe")
		}
		return db.Producto{}, err
	}

	return l.queries.CreateProducto(ctx, db.CreateProductoParams{
		NombreProd:      req.NombreProd,
		DescripcionProd: req.DescripcionProd,
		Precio:          req.Precio,
		Stock:           req.Stock,
		IDCat:           req.IDCat,
	})
}

func (l *ProductoLogic) UpdateProducto(ctx context.Context, id int32, req models.ProductoRequest) (db.Producto, error) {
	if err := ValidateProducto(req); err != nil {
		return db.Producto{}, err
	}

	prod, err := l.queries.GetProducto(ctx, id)
	if err != nil {
		return db.Producto{}, err
	}

	err = l.queries.UpdateProducto(ctx, db.UpdateProductoParams{
		IDProd:          prod.IDProd,
		NombreProd:      req.NombreProd,
		DescripcionProd: req.DescripcionProd,
		Precio:          req.Precio,
		Stock:           req.Stock,
	})
	if err != nil {
		return db.Producto{}, err
	}

	return l.queries.GetProducto(ctx, id)
}

func (l *ProductoLogic) DeleteProducto(ctx context.Context, id int32) error {
	prod, err := l.queries.GetProducto(ctx, id)
	if err != nil {
		return err
	}

	return l.queries.DeleteProducto(ctx, prod.IDProd)
}
