-- ============================================
-- MIGRACIÓN: Actualizar relación Destetes-Lotes
-- Cambios:
-- 1. Lotes pueden contener lechones de múltiples destetes (N:1)
-- 2. Lote_id es obligatorio en destetes
-- 3. Corral_id es obligatorio en lotes
-- 4. Se quita destete_id de lotes
-- ============================================

-- Paso 1: Si hay datos existentes, hacer backup primero
-- (Este script asume que no hay datos, si hay datos se debe migrar primero)

-- Paso 2: Eliminar foreign keys y restricciones existentes
ALTER TABLE destetes DROP FOREIGN KEY IF EXISTS destetes_ibfk_3;
ALTER TABLE lotes DROP FOREIGN KEY IF EXISTS lotes_ibfk_1;
ALTER TABLE lotes DROP INDEX IF EXISTS unique_destete_lote;
ALTER TABLE lotes DROP INDEX IF EXISTS idx_destete;

-- Paso 3: Eliminar columna destete_id de lotes
ALTER TABLE lotes DROP COLUMN IF EXISTS destete_id;

-- Paso 4: Hacer corral_id obligatorio en lotes
ALTER TABLE lotes MODIFY corral_id INT NOT NULL;

-- Paso 5: Hacer lote_id obligatorio en destetes
ALTER TABLE destetes MODIFY lote_id INT NOT NULL;

-- Paso 6: Agregar foreign key de destetes a lotes
ALTER TABLE destetes 
ADD FOREIGN KEY (lote_id) REFERENCES lotes(id) ON DELETE RESTRICT;

-- Nota: Si hay datos existentes, se debe:
-- 1. Crear lotes para destetes existentes sin lote
-- 2. Asignar destetes a lotes apropiados
-- 3. Actualizar cantidad_lechones en lotes (suma de destetes asociados)
