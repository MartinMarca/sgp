-- Migración: roles simplificados (admin, propietario, empleado)
-- Ejecutar sobre una BD existente antes de desplegar el nuevo backend.

-- 1. usuarios: propietario_id + nuevo enum de rol
ALTER TABLE usuarios
    ADD COLUMN propietario_id INT NULL AFTER rol,
    ADD CONSTRAINT fk_usuarios_propietario
        FOREIGN KEY (propietario_id) REFERENCES usuarios(id) ON DELETE SET NULL;

ALTER TABLE usuarios
    MODIFY COLUMN rol ENUM('admin', 'usuario', 'veterinario', 'propietario', 'empleado') DEFAULT 'usuario';

UPDATE usuarios SET rol = 'propietario' WHERE rol IN ('usuario', 'veterinario');

ALTER TABLE usuarios
    MODIFY COLUMN rol ENUM('admin', 'propietario', 'empleado') DEFAULT 'propietario';

-- 2. granjas: propietario_id
ALTER TABLE granjas
    ADD COLUMN propietario_id INT NULL AFTER activo;

-- Poblar desde usuario_granja si existe
UPDATE granjas g
SET propietario_id = (
    SELECT ug.usuario_id FROM usuario_granja ug
    WHERE ug.granja_id = g.id AND ug.rol = 'propietario' LIMIT 1
)
WHERE (g.propietario_id IS NULL OR g.propietario_id = 0);

-- Fallback: primer usuario propietario o admin
UPDATE granjas g
SET propietario_id = (SELECT MIN(id) FROM usuarios WHERE rol IN ('propietario', 'admin'))
WHERE (g.propietario_id IS NULL OR g.propietario_id = 0)
  AND EXISTS (SELECT 1 FROM usuarios WHERE rol IN ('propietario', 'admin'));

-- Último fallback: cualquier usuario
UPDATE granjas g
SET propietario_id = (SELECT MIN(id) FROM usuarios)
WHERE (g.propietario_id IS NULL OR g.propietario_id = 0)
  AND EXISTS (SELECT 1 FROM usuarios);

ALTER TABLE granjas
    MODIFY COLUMN propietario_id INT NOT NULL,
    ADD CONSTRAINT fk_granjas_propietario
        FOREIGN KEY (propietario_id) REFERENCES usuarios(id),
    ADD INDEX idx_propietario (propietario_id);

-- 3. Bootstrap admin si no hay ninguno
UPDATE usuarios SET rol = 'admin' WHERE id = (SELECT MIN(id) FROM (SELECT id FROM usuarios) AS u)
    AND NOT EXISTS (SELECT 1 FROM usuarios WHERE rol = 'admin');
