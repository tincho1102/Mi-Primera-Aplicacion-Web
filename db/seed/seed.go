// seed carga datos iniciales en la base para que, lo hicimos para que al abrir la
// aplicación, ya haya categorías y productos con los que interactuar.
package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	db "matematicos.com/servidor-go-tp1/db/sqlc"
)

type categoriaSeed struct {
	nombre      string
	descripcion string
	padre       string
}

type productoSeed struct {
	nombre      string
	descripcion string
	precio      string
	stock       int32
	categoria   string // nombre de la categoría a la que pertenece
}

// Las categorías padre van antes que sus hijas, porque al insertar una hija
// se necesita el id de la padre.
var categorias = []categoriaSeed{
	{"Mates", "Mates de calabaza, madera, vidrio y acero", ""},
	{"Mates de calabaza", "Mates tradicionales de calabaza curada", "Mates"},
	{"Termos", "Termos de acero inoxidable para el agua del mate", ""},
	{"Bombillas", "Bombillas de alpaca, acero y caña", ""},
	{"Materas", "Bolsos y canastas para llevar el equipo de mate", ""},
}

var productos = []productoSeed{
	{"Mate imperial de alpaca", "Calabaza forrada en cuero con virola de alpaca", "32000.00", 5, "Mates de calabaza"},
	{"Mate de calabaza forrado", "Calabaza curada forrada en cuero vacuno", "8500.00", 15, "Mates de calabaza"},
	{"Mate de madera de algarrobo", "Mate torneado en madera de algarrobo", "6000.00", 20, "Mates"},
	{"Mate de vidrio forrado", "Vaso de vidrio con forro de cuero sintético", "5500.00", 12, "Mates"},
	{"Termo acero inoxidable 1 L", "Termo de doble pared, mantiene el agua caliente 24 h", "45000.00", 8, "Termos"},
	{"Termo con pico cebador 750ml", "Termo con pico cebador para un chorro preciso", "28000.00", 10, "Termos"},
	{"Bombilla de alpaca", "Bombilla de alpaca con filtro de resorte", "12000.00", 25, "Bombillas"},
	{"Bombilla pico de loro", "Bombilla de acero con pico curvo", "9000.00", 30, "Bombillas"},
	{"Bombilla de acero inoxidable", "Bombilla recta con filtro desmontable", "4500.00", 40, "Bombillas"},
	{"Matera de cuero", "Matera de cuero con lugar para termo, mate y yerbera", "38000.00", 6, "Materas"},
	{"Matera de lona", "Matera liviana de lona impermeable", "18000.00", 9, "Materas"},
}

// Run carga las categorías y productos por defecto, solo si la tabla de
// categorías está vacía (sin categorías tampoco puede haber productos, porque
// producto.id_cat es NOT NULL). Así no se duplican datos si la base ya tenía
// información. Todo se hace en una transacción o sea, se carga todo o nada.
func Run(ctx context.Context, sqlDB *sql.DB) error {
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // no hace nada si ya se hizo Commit

	q := db.New(sqlDB).WithTx(tx)

	existentes, err := q.ListCategorias(ctx)
	if err != nil {
		return err
	}
	if len(existentes) > 0 {
		log.Println("Seed: la base ya tiene datos, no se cargan los iniciales")
		return nil
	}

	ids := make(map[string]int32, len(categorias))
	for _, c := range categorias {
		params := db.CreateCategoriaParams{NombreCat: c.nombre, DescripcionCat: c.descripcion}
		if c.padre != "" {
			params.IDPadre = sql.NullInt32{Int32: ids[c.padre], Valid: true}
		}
		creada, err := q.CreateCategoria(ctx, params)
		if err != nil {
			return fmt.Errorf("categoria %q: %w", c.nombre, err)
		}
		ids[c.nombre] = creada.IDCat
	}

	for _, p := range productos {
		_, err := q.CreateProducto(ctx, db.CreateProductoParams{
			NombreProd:      p.nombre,
			DescripcionProd: p.descripcion,
			Precio:          p.precio,
			Stock:           p.stock,
			IDCat:           ids[p.categoria],
		})
		if err != nil {
			return fmt.Errorf("producto %q: %w", p.nombre, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("Seed: se cargaron %d categorias y %d productos iniciales", len(categorias), len(productos))
	return nil
}
