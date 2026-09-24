-- name: CreateProducto :one
INSERT INTO producto (nombre_prod, descripcion_prod, precio, stock, id_cat)
VALUES ($1, $2, $3, $4, $5)
RETURNING id_prod, nombre_prod, descripcion_prod, precio, stock, id_cat created_at;

-- name: CreateCategoria :one
INSERT INTO categoria(nombre_cat, descripcion_cat, id_padre)
VALUES ($1, $2, $3)
RETURNING id_cat, nombre_cat, descripcion_cat, id_padre, created_at;

-- name: GetCategoria :one
SELECT id_cat, nombre_cat, descripcion_cat, id_padre, created_at
FROM categoria
WHERE id_cat = $1;

-- name: GetProducto :one
SELECT id_prod, nombre_prod, descripcion_prod, precio, stock, id_cat, created_at
FROM producto
WHERE id_prod = $1;

-- name: ListCategorias :many
SELECT id_cat, nombre_cat, descripcion_cat,id_padre, created_at
FROM categoria
ORDER BY nombre_cat ASC; 

-- name: ListProductos :many
SELECT id_prod, nombre_prod, descripcion_prod, precio, stock, id_cat, created_at
FROM producto
ORDER BY nombre_prod ASC;

-- name: ListProductosByCategoria :many
SELECT id_prod, nombre_prod, descripcion_prod, precio, stock, id_cat, created_at
FROM producto
WHERE id_cat = $1
ORDER BY nombre_prod ASC;

-- name: UpdateCategoria :exec
UPDATE categoria
SET nombre_cat = $2, descripcion_cat = $3
WHERE id_cat = $1;

-- name: UpdateProducto :exec
UPDATE producto
SET nombre_prod = $2, descripcion_prod = $3, precio = $4, stock = $5
WHERE id_prod = $1;

-- name: UpdateStockProducto :exec
UPDATE producto
SET stock = $2
WHERE id_prod = $1;

-- name: UpdatePrecioProducto :exec
UPDATE producto
SET precio = $2
WHERE id_prod = $1;

-- name: DeleteCategoria :exec
DELETE FROM categoria
WHERE id_cat = $1; 

-- name: DeleteProducto :exec
DELETE FROM producto
WHERE id_prod = $1;