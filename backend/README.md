# Backend - Sistema de Gestión de Granja Porcina

Backend desarrollado en Go con Gin Framework y GORM.

## Requisitos

- Go 1.21 o superior
- MySQL 8.0 o superior
- Make (opcional, para comandos útiles)

## Instalación

1. **Clonar el repositorio** (si aún no lo has hecho)

2. **Configurar variables de entorno**
```bash
cp .env.example .env
# Editar .env con tus credenciales de MySQL
```

3. **Instalar dependencias**
```bash
make install
# o manualmente:
go mod download
go mod tidy
```

4. **Crear la base de datos**
```bash
mysql -u root -p
CREATE DATABASE granja_porcina;
```

5. **Ejecutar migraciones** (opcional, si usas AUTO_MIGRATE=true)
```bash
# Opción 1: AutoMigrate de GORM (development)
AUTO_MIGRATE=true make run

# Opción 2: SQL manual (production)
mysql -u root -p granja_porcina < ../database/schema.sql
```

## Ejecución

### Desarrollo
```bash
# Ejecutar servidor
make run

# Ejecutar con hot-reload (requiere air)
go install github.com/cosmtrek/air@latest
make dev
```

### Producción
```bash
# Compilar binario
make build

# Ejecutar binario
./bin/sgp-server
```

## Estructura del Proyecto

```
backend/
├── cmd/server/          # Punto de entrada de la aplicación
│   └── main.go         # ✅ Implementado
├── internal/            # Código privado de la aplicación
│   ├── config/         # ✅ Configuración
│   │   └── config.go
│   ├── database/       # ✅ Conexión y migraciones
│   │   ├── database.go
│   │   └── migrations.go
│   ├── handlers/       # ⏳ TODO: Controladores HTTP
│   ├── middleware/     # ✅ Middlewares (auth, cors)
│   │   ├── auth.go
│   │   └── cors.go
│   ├── models/         # ✅ Modelos GORM (10 modelos)
│   │   ├── usuario.go
│   │   ├── granja.go
│   │   ├── usuario_granja.go
│   │   ├── corral.go
│   │   ├── lote.go
│   │   ├── padrillo.go
│   │   ├── cerda.go
│   │   ├── servicio.go
│   │   ├── parto.go
│   │   ├── destete.go
│   │   ├── models.go   # Constantes y helpers
│   │   └── README.md   # Documentación de modelos
│   ├── repositories/   # ⏳ TODO: Capa de acceso a datos
│   ├── routes/         # ✅ Definición de rutas
│   │   └── routes.go
│   ├── services/       # ⏳ TODO: Lógica de negocio
│   └── utils/          # ✅ Utilidades
│       ├── jwt.go
│       └── response.go
├── migrations/         # Migraciones SQL (opcional)
├── .env.example        # ✅ Ejemplo de variables de entorno
├── go.mod             # ✅ Dependencias Go
├── Makefile           # ✅ Comandos útiles
└── README.md          # Este archivo
```

## API Endpoints

### Health Check
- `GET /api/health` - Verificar estado del servidor

### Autenticación (TODO)
- `POST /api/auth/login` - Iniciar sesión
- `POST /api/auth/register` - Registrar usuario

### Granjas (TODO)
- `GET /api/granjas` - Listar granjas
- `POST /api/granjas` - Crear granja
- `GET /api/granjas/:id` - Obtener granja
- `PUT /api/granjas/:id` - Actualizar granja
- `DELETE /api/granjas/:id` - Eliminar granja

### Cerdas (TODO)
- `GET /api/cerdas` - Listar cerdas
- `POST /api/cerdas` - Crear cerda
- `GET /api/cerdas/:id` - Obtener cerda
- `PUT /api/cerdas/:id` - Actualizar cerda
- `DELETE /api/cerdas/:id` - Dar de baja cerda

*(Más endpoints por implementar...)*

## Testing

```bash
make test
```

## Comandos Útiles

Ver `Makefile` para todos los comandos disponibles.

## Notas de Desarrollo

- ✅ **Modelos GORM implementados** en `internal/models/`
  - 10 modelos con todas las relaciones
  - Validaciones en hooks (BeforeCreate, BeforeUpdate)
  - Soft deletes en todos los modelos principales
  - Ver `internal/models/README.md` para documentación detallada

- ✅ **Repositories implementados** en `internal/repositories/`
  - 11 repositories (base + 9 entidades + container)
  - CRUD completo para todas las entidades
  - Queries específicas (historial, estadísticas, eventos futuros)
  - Soporte para transacciones, filtros y preload
  - Ver `internal/repositories/README.md` para documentación detallada

- ⏳ **Pendiente**: La lógica de negocio va en `internal/services/`
- ⏳ **Pendiente**: Los handlers HTTP van en `internal/handlers/`

## Modelos Implementados

1. ✅ `Usuario` - Usuarios del sistema
2. ✅ `Granja` - Granjas porcinas
3. ✅ `UsuarioGranja` - Relación N:M usuarios-granjas
4. ✅ `Corral` - Corrales dentro de granjas
5. ✅ `Lote` - Lotes de lechones destetados
6. ✅ `Padrillo` - Padrillos (machos reproductores)
7. ✅ `Cerda` - Cerdas (hembras reproductoras)
8. ✅ `Servicio` - Servicios de monta
9. ✅ `Parto` - Partos de cerdas
10. ✅ `Destete` - Destetes de lechones

### Características de los Modelos

- **Validaciones**: Hooks de GORM para validar reglas de negocio
  - `Servicio`: Valida campos según tipo de monta
  - `Parto`: Valida que hembras + machos == vivos
  - `Destete`: Valida cantidad <= lechones del parto
  
- **Relaciones**: Todas las foreign keys configuradas
- **Soft Deletes**: Implementado con `gorm.DeletedAt`
- **Timestamps**: `CreatedAt` y `UpdatedAt` automáticos
- **Índices**: Configurados para optimizar queries

Ver documentación completa en: `internal/models/README.md`

## Repositories Implementados

1. ✅ `BaseRepository` - Operaciones CRUD genéricas
2. ✅ `UsuarioRepository` - Autenticación, gestión de usuarios
3. ✅ `GranjaRepository` - Gestión de granjas y asignación de usuarios
4. ✅ `CorralRepository` - Gestión de corrales y ocupación
5. ✅ `LoteRepository` - Gestión de lotes y destetes asociados
6. ✅ `PadrilloRepository` - Gestión de padrillos y estadísticas
7. ✅ `CerdaRepository` - Gestión de cerdas, cambios de estado, historial
8. ✅ `ServicioRepository` - Gestión de servicios, confirmación de preñez
9. ✅ `PartoRepository` - Gestión de partos, partos futuros
10. ✅ `DesteteRepository` - Gestión de destetes, destetes futuros
11. ✅ `RepositoryContainer` - Container para inyección de dependencias

### Características de los Repositories

- **CRUD Completo**: Create, FindByID, Update, Delete (soft delete)
- **Preload de Relaciones**: Soporte para cargar relaciones con `Preload()`
- **Filtros**: Búsqueda por estado, granja, período (mes/año)
- **Validaciones**: Verificación de duplicados, existencia de registros
- **Queries Específicas**: Historial, estadísticas, eventos futuros
- **Transacciones**: Soporte para operaciones atómicas
- **Optimizado**: Queries con índices y joins eficientes

Ver documentación completa en: `internal/repositories/README.md`

## Próximos Pasos

1. ✅ ~~Implementar **Modelos GORM**~~ - Completado (10 modelos)
2. ✅ ~~Implementar **Repositories**~~ - Completado (11 repositories)
3. ⏳ Implementar **Services** - Lógica de negocio (ciclo de maternidad, cálculos)
4. ⏳ Implementar **Handlers** - Controladores HTTP REST
5. ⏳ Testing - Tests unitarios y de integración

---

**Estado actual**: ✅ Modelos y Repositories completamente implementados  
**Última actualización**: 2026-02-07

