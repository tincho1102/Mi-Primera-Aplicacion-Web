package main

import (
	"fmt"
)

func main(){

	fmt.Println("Servidor Go - TP Persistencia")
	fmt.Println("Para ejecutar las pruebas del CRUD y el pipeline automatizado, utiliza: make test")
      //connStr := "user=cliente password=progwebunicen dbname=mate_db"
      //no funcionaba y lo consulté con gemini
	  /*
	  connStr := "postgres://cliente:progwebunicen@127.0.0.1:5432/mate_db?sslmode=disable"
      db, err := sql.Open("pgx", connStr)
      if err != nil {
        log.Fatalf("failed to connect to DB: %v", err)
      }
      defer db.Close()
      queries := sqlc.New(db)
      ctx := context.Background()
      

      // crud

	//Crear una categoría
	createdCat, err := queries.CreateCategoria(ctx, sqlc.CreateCategoriaParams{
		NombreCat:      "Mates",
		DescripcionCat: "Mates de todo tipo",
		IDPadre: sql.NullInt32{Valid: false},
	})
	if err != nil {
		log.Fatalf("failed to create categoria: %v", err)
	}
	fmt.Printf("Created categoria: %+v\n", createdCat)

	//  Obtener categoría por ID -> lee una sola
	cat, err := queries.GetCategoria(ctx, createdCat.IDCat)
	if err != nil {
		log.Fatalf("failed to get categoria: %v", err)
	}
	fmt.Printf("Retrieved categoria: %+v\n", cat)

	// Lista de todas las categorias
	cats, err := queries.ListCategorias(ctx)
	if err != nil {
		log.Fatalf("failed to list categorias: %v", err)
	}
	fmt.Printf("All categorias: %+v\n", cats)

	// Actualizar categoría
	err = queries.UpdateCategoria(ctx, sqlc.UpdateCategoriaParams{
		IDCat:          createdCat.IDCat,
		NombreCat:      "Mates",
		DescripcionCat: "Cualquier tipo de mate", // Pasa un string directo
	})
	if err != nil {
		log.Fatalf("failed to update categoria: %v", err)
	}
	fmt.Println("Categoria updated successfully")

	updatedCat, err := queries.GetCategoria(ctx, createdCat.IDCat)
	if err != nil {
		log.Fatalf("failed to get updated categoria: %v", err)
	}
	fmt.Printf("Updated categoria: %+v\n", updatedCat)

	// Eliminar  1 categoría
	err = queries.DeleteCategoria(ctx, createdCat.IDCat)
	if err != nil {
		log.Fatalf("failed to delete categoria: %v", err)
	}
	fmt.Println("Categoria deleted successfully")

	_, err = queries.GetCategoria(ctx, createdCat.IDCat)
	if err == sql.ErrNoRows {
		fmt.Println("Categoria not found after deletion")
	} else if err != nil {
		log.Fatalf("failed to get categoria after deletion: %v", err)
	}

	*/
}