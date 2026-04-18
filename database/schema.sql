-- ============================================
-- ESQUEMA DE BASE DE DATOS
-- Sistema de Gestión de Granja Porcina
-- ============================================

-- Tabla de Usuarios
CREATE TABLE usuarios (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nombre_completo VARCHAR(100),
    rol ENUM('admin', 'usuario', 'veterinario') DEFAULT 'usuario',
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla de Granjas
CREATE TABLE granjas (
    id INT PRIMARY KEY AUTO_INCREMENT,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    ubicacion VARCHAR(200),
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_nombre (nombre),
    INDEX idx_activo (activo)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla intermedia: Usuarios-Granjas (relación N:M)
CREATE TABLE usuario_granja (
    id INT PRIMARY KEY AUTO_INCREMENT,
    usuario_id INT NOT NULL,
    granja_id INT NOT NULL,
    rol ENUM('propietario', 'administrador', 'operador') DEFAULT 'operador',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    FOREIGN KEY (granja_id) REFERENCES granjas(id) ON DELETE CASCADE,
    UNIQUE KEY unique_usuario_granja (usuario_id, granja_id),
    INDEX idx_usuario (usuario_id),
    INDEX idx_granja (granja_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla de Corrales
CREATE TABLE corrales (
    id INT PRIMARY KEY AUTO_INCREMENT,
    granja_id INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    capacidad_maxima INT NULL,
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (granja_id) REFERENCES granjas(id) ON DELETE RESTRICT,
    INDEX idx_granja (granja_id),
    INDEX idx_nombre (nombre),
    INDEX idx_activo (activo)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla de Padrillos
CREATE TABLE padrillos (
    id INT PRIMARY KEY AUTO_INCREMENT,
    granja_id INT NOT NULL,
    numero_caravana VARCHAR(50) NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    genetica VARCHAR(100) NULL,
    fecha_ultima_vacunacion DATE,
    activo BOOLEAN DEFAULT TRUE,
    fecha_baja DATE NULL,
    motivo_baja ENUM('muerte', 'venta') NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (granja_id) REFERENCES granjas(id) ON DELETE RESTRICT,
    UNIQUE KEY unique_caravana_granja (numero_caravana, granja_id),
    INDEX idx_granja (granja_id),
    INDEX idx_caravana (numero_caravana),
    INDEX idx_nombre (nombre),
    INDEX idx_activo (activo)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla de Cerdas
CREATE TABLE cerdas (
    id INT PRIMARY KEY AUTO_INCREMENT,
    granja_id INT NOT NULL,
    numero_caravana VARCHAR(50) NOT NULL,
    detalle_pelaje TEXT,
    genetica VARCHAR(100) NULL,
    estado ENUM('disponible', 'servicio', 'gestacion', 'cria') DEFAULT 'disponible',
    activo BOOLEAN DEFAULT TRUE,
    fecha_baja DATE NULL,
    motivo_baja ENUM('muerte', 'venta') NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (granja_id) REFERENCES granjas(id) ON DELETE RESTRICT,
    UNIQUE KEY unique_caravana_granja (numero_caravana, granja_id),
    INDEX idx_granja (granja_id),
    INDEX idx_caravana (numero_caravana),
    INDEX idx_estado (estado),
    INDEX idx_activo (activo)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla de Servicios
CREATE TABLE servicios (
    id INT PRIMARY KEY AUTO_INCREMENT,
    cerda_id INT NOT NULL,
    fecha_servicio DATE NOT NULL,
    tiene_repeticiones BOOLEAN DEFAULT FALSE,
    cantidad_repeticiones INT DEFAULT 0,
    tipo_monta ENUM('natural', 'inseminacion') NOT NULL,
    -- Campos para Monta Natural
    padrillo_id INT NULL,
    cantidad_saltos INT NULL,
    -- Campos para Inseminación
    numero_pajuela VARCHAR(50) NULL,
    prenez_confirmada BOOLEAN DEFAULT FALSE,
    fecha_confirmacion_prenez DATE NULL,
    prenez_cancelada BOOLEAN DEFAULT FALSE,
    fecha_cancelacion_prenez DATE NULL,
    motivo_cancelacion TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (cerda_id) REFERENCES cerdas(id) ON DELETE RESTRICT,
    FOREIGN KEY (padrillo_id) REFERENCES padrillos(id) ON DELETE SET NULL,
    INDEX idx_cerda (cerda_id),
    INDEX idx_fecha (fecha_servicio),
    INDEX idx_prenez_confirmada (prenez_confirmada),
    INDEX idx_mes (fecha_servicio)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla de Partos
CREATE TABLE partos (
    id INT PRIMARY KEY AUTO_INCREMENT,
    cerda_id INT NOT NULL,
    servicio_id INT NULL, -- Relación con el servicio que originó este parto
    fecha_parto DATE NOT NULL,
    lechones_nacidos_vivos INT NOT NULL DEFAULT 0,
    lechones_nacidos_totales INT NOT NULL DEFAULT 0,
    lechones_hembras INT NOT NULL DEFAULT 0,
    lechones_machos INT NOT NULL DEFAULT 0,
    fecha_estimada DATE NOT NULL, -- Calculada: fecha_servicio + 114 días
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (cerda_id) REFERENCES cerdas(id) ON DELETE RESTRICT,
    FOREIGN KEY (servicio_id) REFERENCES servicios(id) ON DELETE SET NULL,
    INDEX idx_cerda (cerda_id),
    INDEX idx_fecha (fecha_parto),
    INDEX idx_fecha_estimada (fecha_estimada),
    INDEX idx_mes (fecha_parto)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla de Destetes
CREATE TABLE destetes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    cerda_id INT NOT NULL,
    parto_id INT NULL, -- Relación con el parto que originó este destete
    fecha_destete DATE NOT NULL,
    cantidad_lechones_destetados INT NOT NULL DEFAULT 0,
    fecha_estimada DATE NOT NULL, -- Calculada: fecha_parto + 30 días
    lote_id INT NULL, -- Relación con el lote generado (1:1)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (cerda_id) REFERENCES cerdas(id) ON DELETE RESTRICT,
    FOREIGN KEY (parto_id) REFERENCES partos(id) ON DELETE SET NULL,
    INDEX idx_cerda (cerda_id),
    INDEX idx_fecha (fecha_destete),
    INDEX idx_fecha_estimada (fecha_estimada),
    INDEX idx_mes (fecha_destete),
    INDEX idx_lote (lote_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabla de Lotes
CREATE TABLE lotes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    destete_id INT NOT NULL, -- Lote se crea a partir de un destete
    corral_id INT NULL, -- Lote puede estar en un corral (opcional al crear)
    nombre VARCHAR(100) NOT NULL,
    cantidad_lechones INT NOT NULL DEFAULT 0,
    fecha_creacion DATE NOT NULL,
    estado ENUM('activo', 'cerrado', 'vendido') DEFAULT 'activo',
    fecha_cierre DATE NULL,
    motivo_cierre TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (destete_id) REFERENCES destetes(id) ON DELETE RESTRICT,
    FOREIGN KEY (corral_id) REFERENCES corrales(id) ON DELETE SET NULL,
    UNIQUE KEY unique_destete_lote (destete_id), -- Un destete genera un solo lote
    INDEX idx_corral (corral_id),
    INDEX idx_destete (destete_id),
    INDEX idx_estado (estado),
    INDEX idx_fecha_creacion (fecha_creacion)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Agregar foreign key de destetes a lotes (después de crear la tabla lotes)
ALTER TABLE destetes 
ADD FOREIGN KEY (lote_id) REFERENCES lotes(id) ON DELETE SET NULL;

-- Tabla de Estadísticas (caché opcional para mejorar rendimiento)
CREATE TABLE estadisticas_cache (
    id INT PRIMARY KEY AUTO_INCREMENT,
    tipo ENUM('cerdas', 'servicios', 'partos', 'destetes') NOT NULL,
    periodo VARCHAR(20) NOT NULL, -- '2025-01', '2025', 'historico'
    datos JSON NOT NULL,
    fecha_calculo TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tipo_periodo (tipo, periodo)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- VISTAS ÚTILES
-- ============================================

-- Vista: Cerdas con su último estado y estadísticas básicas
CREATE OR REPLACE VIEW vista_cerdas_detalle AS
SELECT 
    c.id,
    c.granja_id,
    g.nombre as granja_nombre,
    c.numero_caravana,
    c.detalle_pelaje,
    c.genetica,
    c.estado,
    c.activo,
    COUNT(DISTINCT s.id) as total_servicios,
    COUNT(DISTINCT p.id) as total_partos,
    COUNT(DISTINCT CASE WHEN s.prenez_confirmada = TRUE THEN s.id END) as servicios_exitosos,
    COALESCE(AVG(p.lechones_nacidos_vivos), 0) as promedio_lechones,
    c.created_at,
    c.updated_at
FROM cerdas c
JOIN granjas g ON g.id = c.granja_id
LEFT JOIN servicios s ON s.cerda_id = c.id
LEFT JOIN partos p ON p.cerda_id = c.id
WHERE c.activo = TRUE
GROUP BY c.id;

-- Vista: Partos y destetes futuros (para calendario)
CREATE OR REPLACE VIEW vista_eventos_futuros AS
SELECT 
    'parto' as tipo_evento,
    p.id,
    p.cerda_id,
    c.numero_caravana,
    p.fecha_estimada as fecha_evento,
    p.fecha_parto as fecha_real,
    s.fecha_servicio,
    NULL as parto_id
FROM partos p
JOIN cerdas c ON c.id = p.cerda_id
LEFT JOIN servicios s ON s.id = p.servicio_id
WHERE p.fecha_estimada >= CURDATE() AND p.fecha_parto IS NULL

UNION ALL

SELECT 
    'destete' as tipo_evento,
    d.id,
    d.cerda_id,
    c.numero_caravana,
    d.fecha_estimada as fecha_evento,
    d.fecha_destete as fecha_real,
    p.fecha_parto as fecha_servicio,
    d.parto_id
FROM destetes d
JOIN cerdas c ON c.id = d.cerda_id
LEFT JOIN partos p ON p.id = d.parto_id
WHERE d.fecha_estimada >= CURDATE() AND d.fecha_destete IS NULL;

-- ============================================
-- ÍNDICES ADICIONALES PARA RENDIMIENTO
-- ============================================

-- Índices compuestos para consultas frecuentes
CREATE INDEX idx_servicios_cerda_fecha ON servicios(cerda_id, fecha_servicio);
CREATE INDEX idx_partos_cerda_fecha ON partos(cerda_id, fecha_parto);
CREATE INDEX idx_destetes_cerda_fecha ON destetes(cerda_id, fecha_destete);
CREATE INDEX idx_lotes_corral_estado ON lotes(corral_id, estado);
CREATE INDEX idx_corrales_granja_activo ON corrales(granja_id, activo);