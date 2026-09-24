package main

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	sqlc "matematicos.com/servidor-go-tp1/db/sqlc"
)

func TestCRUD(t *testing.T) {
	connStr := "postgres://cliente:progwebunicen@127.0.0.1:5432/mate_db?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	ctx := context.Background()

	//CRUD para categorias, en esta parte se crea una categoria, se listan, y se actualizan

	createdCat, err := queries.CreateCategoria(ctx, sqlc.CreateCategoriaParams{
		NombreCat:      "Mates",
		DescripcionCat: "Mates de todo tipo",
		IDPadre:        sql.NullInt32{Valid: false},
	})
	if err != nil {
		t.Fatalf("failed to create categoria: %v", err)
	}
	t.Logf("Created categoria: %+v", createdCat)

	cat, err := queries.GetCategoria(ctx, createdCat.IDCat)
	if err != nil {
		t.Fatalf("failed to get categoria: %v", err)
	}
	t.Logf("Retrieved categoria: %+v", cat)

	cats, err := queries.ListCategorias(ctx)
	if err != nil {
		t.Fatalf("failed to list categorias: %v", err)
	}
	t.Logf("All categorias: %+v", cats)

	err = queries.UpdateCategoria(ctx, sqlc.UpdateCategoriaParams{
		IDCat:          createdCat.IDCat,
		NombreCat:      "Mates",
		DescripcionCat: "Cualquier tipo de mate",
	})
	if err != nil {
		t.Fatalf("failed to update categoria: %v", err)
	}
	t.Log("Categoria updated successfully")

	// PARTE DE PRODUCTOS: crea un producto, hace un get, update y lista todos los productos que en este caso es uno solo

	prodCreado, err := queries.CreateProducto(ctx, sqlc.CreateProductoParams{
		NombreProd:      "Mate Torpedo",
		DescripcionProd: "Mate de cuero cincelado",
		Precio:          "30000",
		Stock:           10,
		IDCat:           createdCat.IDCat,
	})
	if err != nil {
		t.Fatalf("Error en CreateProducto: %v", err)
	}
	t.Logf("Producto creado con éxito (ID: %d)", prodCreado.IDProd)

	prodObtenido, err := queries.GetProducto(ctx, prodCreado.IDProd)
	if err != nil {
		t.Fatalf("Error en GetProducto: %v", err)
	}
	t.Logf("Producto obtenido: %s", prodObtenido.NombreProd)

	err = queries.UpdateProducto(ctx, sqlc.UpdateProductoParams{
		IDProd:          prodCreado.IDProd,
		NombreProd:      "Mate Imperial",
		DescripcionProd: "Con virola de alpaca",
		Precio:          "30000.00",
		Stock:           15,
	})
	if err != nil {
		t.Fatalf("Error en UpdateProducto: %v", err)
	}
	// tambien se pueden hacer update de stock, o precio por separado
	err = queries.UpdateStockProducto(ctx, sqlc.UpdateStockProductoParams{
		IDProd: prodCreado.IDProd,
		Stock:  20,
	})
	if err != nil {
		t.Fatalf("Error en UpdateStockProducto: %v", err)
	}

	err = queries.UpdatePrecioProducto(ctx, sqlc.UpdatePrecioProductoParams{
		IDProd: prodCreado.IDProd,
		Precio: "28500.00",
	})
	if err != nil {
		t.Fatalf("Error en UpdatePrecioProducto: %v", err)
	}

	prods, err := queries.ListProductos(ctx)
	if err != nil || len(prods) == 0 {
		t.Fatalf("Error o lista vacía en ListProductos: %v", err)
	}

	prodsCat, err := queries.ListProductosByCategoria(ctx, createdCat.IDCat)
	if err != nil || len(prodsCat) == 0 {
		t.Fatalf("Error en ListProductosByCategoria: %v", err)
	}

	// Parte de deletes. se borra el producto y la categoria

	err = queries.DeleteProducto(ctx, prodCreado.IDProd)
	if err != nil {
		t.Fatalf("Error en DeleteProducto: %v", err)
	}

	err = queries.DeleteCategoria(ctx, createdCat.IDCat)
	if err != nil {
		t.Fatalf("failed to delete categoria: %v", err)
	}
	t.Log("Categoria deleted successfully")

	_, err = queries.GetCategoria(ctx, createdCat.IDCat)
	if err == sql.ErrNoRows {
		t.Log("Categoria not found after deletion")
	} else if err != nil {
		t.Fatalf("failed to get categoria after deletion: %v", err)
	}
}