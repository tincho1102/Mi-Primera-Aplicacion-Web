CREATE TABLE categoria (
    id_cat SERIAL PRIMARY KEY, 
    nombre_cat VARCHAR(20) UNIQUE NOT NULL, 
    descripcion_cat VARCHAR(255) NOT NULL,
    id_padre INT NULL REFERENCES categoria(id_cat), 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE producto (
    id_prod SERIAL PRIMARY KEY, 
    nombre_prod VARCHAR(30) UNIQUE NOT NULL, 
    descripcion_prod VARCHAR(255) NOT NULL, 
    precio DECIMAL (10, 2) NOT NULL CHECK (precio >= 0),
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0), 
    id_cat INT NOT NULL REFERENCES categoria(id_cat) ON DELETE CASCADE, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)