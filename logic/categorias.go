package logic

import (
	"context"
	"database/sql"

	db "matematicos.com/servidor-go-tp1/db/sqlc"
	"matematicos.com/servidor-go-tp1/models"
)

type CategoriaLogic struct {
	queries *db.Queries
}

func NewCategoriaLogic(q *db.Queries) *CategoriaLogic {
	return &CategoriaLogic{queries: q}
}

func ValidateCategoria(req models.CategoriaRequest) error {
	if req.NombreCat == "" {
		return NewValidationError("El nombre de la categoria no puede estar vacio")
	}
	if req.DescripcionCat == "" {
		return NewValidationError("La descripcion de la categoria no puede estar vacia")
	}
	if req.IDPadre != nil && *req.IDPadre <= 0 {
		return NewValidationError("El id de la categoria padre debe ser mayor a cero o nulo")
	}

	return nil
}

func (l *CategoriaLogic) ListCategorias(ctx context.Context) ([]db.Categoria, error) {
	return l.queries.ListCategorias(ctx)
}

func (l *CategoriaLogic) GetCategoria(ctx context.Context, id int32) (db.Categoria, error) {
	return l.queries.GetCategoria(ctx, id)
}

func (l *CategoriaLogic) CreateCategoria(ctx context.Context, req models.CategoriaRequest) (db.Categoria, error) {
	if err := ValidateCategoria(req); err != nil {
		return db.Categoria{}, err
	}

	var idPadre sql.NullInt32
	if req.IDPadre != nil {
		idPadre = sql.NullInt32{Int32: *req.IDPadre, Valid: true}
	}

	return l.queries.CreateCategoria(ctx, db.CreateCategoriaParams{
		NombreCat:      req.NombreCat,
		DescripcionCat: req.DescripcionCat,
		IDPadre:        idPadre,
	})
}

func (l *CategoriaLogic) UpdateCategoria(ctx context.Context, id int32, req models.CategoriaRequest) (db.Categoria, error) {
	if err := ValidateCategoria(req); err != nil {
		return db.Categoria{}, err
	}

	cat, err := l.queries.GetCategoria(ctx, id)
	if err != nil {
		return db.Categoria{}, err
	}

	err = l.queries.UpdateCategoria(ctx, db.UpdateCategoriaParams{
		IDCat:          cat.IDCat,
		NombreCat:      req.NombreCat,
		DescripcionCat: req.DescripcionCat,
	})
	if err != nil {
		return db.Categoria{}, err
	}

	return l.queries.GetCategoria(ctx, id)
}

func (l *CategoriaLogic) DeleteCategoria(ctx context.Context, id int32) error {
	cat, err := l.queries.GetCategoria(ctx, id)
	if err != nil {
		return err
	}

	return l.queries.DeleteCategoria(ctx, cat.IDCat)
}
