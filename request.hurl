# ===============================
#    PRUEBAS PARA CATEGORIAS Y PRODUCTOS
# ===============================

# 1. Crear categorias
# ----------------------------------------
POST http://localhost:8080/categorias
Content-Type: application/json
{
  "nombre_cat":"Mates",
  "descripcion_cat":"Mates de todo tipo"
}

HTTP/1.1 201
[Asserts]
status == 201
jsonpath "$.nombre_cat" == "Mates"
jsonpath "$.descripcion_cat" == "Mates de todo tipo"

POST http://localhost:8080/categorias
Content-Type: application/json
{
  "nombre_cat":"Termos",
  "descripcion_cat":"Termos de acero y vidrio"
}

HTTP/1.1 201
[Asserts]
status == 201
jsonpath "$.nombre_cat" == "Termos"
jsonpath "$.descripcion_cat" == "Termos de acero y vidrio"

POST http://localhost:8080/categorias
Content-Type: application/json
{
  "nombre_cat":"Accesorios",
  "descripcion_cat":"Bombillas, yerberas y accesorios"
}

HTTP/1.1 201
[Asserts]
status == 201
jsonpath "$.nombre_cat" == "Accesorios"
jsonpath "$.descripcion_cat" == "Bombillas, yerberas y accesorios"

# 2. Listar todas las categorias
# ----------------------------------------------------------------
GET http://localhost:8080/categorias

HTTP/1.1 200
[Asserts]
status == 200
jsonpath "$" count == 3

# 3. Obtener categoria por ID
# ----------------------------------------------
GET http://localhost:8080/categorias/1

HTTP/1.1 200
[Asserts]
status == 200
jsonpath "$.nombre_cat" == "Mates"

# 4. Actualizar categoria
# -------------------------------------------------
PUT http://localhost:8080/categorias/1
Content-Type: application/json
{
  "nombre_cat":"Mates Premium",
  "descripcion_cat":"Mates premium de calabaza y acero"
}

HTTP/1.1 200
[Asserts]
status == 200
jsonpath "$.nombre_cat" == "Mates Premium"
jsonpath "$.descripcion_cat" == "Mates premium de calabaza y acero"

# 5. Crear productos asociados a categorias
# -------------------------------------------------
POST http://localhost:8080/productos
Content-Type: application/json
{
  "nombre_prod":"Mate Torpedo",
  "descripcion_prod":"Mate de madera con virola",
  "precio":"30000",
  "stock":10,
  "id_cat":1
}

HTTP/1.1 201
[Asserts]
status == 201
jsonpath "$.nombre_prod" == "Mate Torpedo"
jsonpath "$.id_cat" == 1

POST http://localhost:8080/productos
Content-Type: application/json
{
  "nombre_prod":"Termo Stanley",
  "descripcion_prod":"Termo de acero de 1 litro",
  "precio":"125000",
  "stock":20,
  "id_cat":2
}

HTTP/1.1 201
[Asserts]
status == 201
jsonpath "$.nombre_prod" == "Termo Stanley"
jsonpath "$.id_cat" == 2

POST http://localhost:8080/productos
Content-Type: application/json
{
  "nombre_prod":"Bombilla de metal",
  "descripcion_prod":"Bombilla de acero inoxidable",
  "precio":"4500",
  "stock":50,
  "id_cat":3
}

HTTP/1.1 201
[Asserts]
status == 201
jsonpath "$.nombre_prod" == "Bombilla de metal"
jsonpath "$.id_cat" == 3

# 6. Listar todos los productos
# ----------------------------------------------------------------
GET http://localhost:8080/productos

HTTP/1.1 200
[Asserts]
status == 200
jsonpath "$" count == 3

# 7. Obtener producto por ID
# ----------------------------------------------------
GET http://localhost:8080/productos/1

HTTP/1.1 200
[Asserts]
status == 200
jsonpath "$.nombre_prod" == "Mate Torpedo"

# 8. Actualizar un producto
# ---------------------------------------------------
PUT http://localhost:8080/productos/1
Content-Type: application/json
{
  "nombre_prod":"Mate Imperial",
  "descripcion_prod":"Mate de madera con virola de alpaca",
  "precio":"35000.00",
  "stock":20,
  "id_cat":1
}

HTTP/1.1 200
[Asserts]
status == 200
jsonpath "$.nombre_prod" == "Mate Imperial"
jsonpath "$.stock" == 20
jsonpath "$.descripcion_prod" == "Mate de madera con virola de alpaca"

# 9. Eliminar producto
# --------------------------------------------------------
DELETE http://localhost:8080/productos/2

HTTP/1.1 204
[Asserts]
status == 204

# 10. Listar productos restantes
# ----------------------------------------------------------------
GET http://localhost:8080/productos

HTTP/1.1 200
[Asserts]
status == 200
jsonpath "$" count == 2

# 11. Eliminar categoria
# -----------------------------------------------------------------------
DELETE http://localhost:8080/categorias/3

HTTP/1.1 204
[Asserts]
status == 204

# 12. Listar productos finales
# -------------------------------------------------------------------------
GET http://localhost:8080/productos

HTTP/1.1 200
[Asserts]
status == 200
jsonpath "$" count == 1
