# Resumen de Modelos GORM - Sistema de Gestión de Granja Porcina

## ✅ Implementación Completada

Se han implementado **10 modelos GORM** con todas sus relaciones, validaciones y características.

## Diagrama de Relaciones

```
┌─────────────┐       N:M        ┌─────────────┐
│  Usuario    │ ←─────────────→  │   Granja    │
└─────────────┘   usuario_granja └──────┬──────┘
                                         │ 1:N
                                         ├──────→ Corrales
                                         ├──────→ Cerdas
                                         └──────→ Padrillos
                                         
┌─────────────┐                 ┌─────────────┐
│   Corral    │                 │  Padrillo   │
└──────┬──────┘                 └──────┬──────┘
       │ 1:N                           │ 1:N
       ↓                               ↓
┌─────────────┐                 ┌─────────────┐
│    Lote     │                 │  Servicio   │
└──────┬──────┘                 └──────┬──────┘
       │ N:1                           │ 1:N
       │                               ↓
       │                        ┌─────────────┐
       │                        │    Parto    │
       │                        └──────┬──────┘
       │                               │ 1:N
       │                               ↓
       │                        ┌─────────────┐
       └──────────────────────→ │   Destete   │
                N:1             └─────────────┘
                                       ↑
                                       │ N:1
                                ┌─────────────┐
                                │   Cerda     │
                                └─────────────┘
                                (Ciclo de Maternidad)
```

## Ciclo de Maternidad (Cerda)

```
┌─────────────┐   Registrar      ┌─────────────┐
│ Disponible  │ ────Servicio───→ │  Servicio   │
└─────────────┘                  └──────┬──────┘
       ↑                                │ Confirmar
       │                                │ Preñez
       │ Registrar                      ↓
       │ Destete              ┌─────────────┐
       │                      │  Gestación  │
┌─────────────┐               └──────┬──────┘
│    Cría     │ ←──────Registrar─────┘
└─────────────┘        Parto
```

## Modelos Implementados

### 1. Usuario (`usuario.go`)
```go
type Usuario struct {
    ID              uint
    Username        string     // Único
    Email           string     // Único
    PasswordHash    string     // No se serializa
    NombreCompleto  *string
    Rol             string     // admin, usuario, veterinario
    Activo          bool
    Granjas         []Granja   // Relación N:M
}
```

### 2. Granja (`granja.go`)
```go
type Granja struct {
    ID          uint
    Nombre      string
    Descripcion *string
    Ubicacion   *string
    Activo      bool
    Usuarios    []Usuario   // N:M
    Corrales    []Corral    // 1:N
    Cerdas      []Cerda     // 1:N
    Padrillos   []Padrillo  // 1:N
}
```

### 3. UsuarioGranja (`usuario_granja.go`)
```go
type UsuarioGranja struct {
    ID        uint
    UsuarioID uint
    GranjaID  uint
    Rol       string  // propietario, administrador, operador
}
```

### 4. Corral (`corral.go`)
```go
type Corral struct {
    ID              uint
    GranjaID        uint       // FK Granja
    Nombre          string
    Descripcion     *string
    CapacidadMaxima *int
    Activo          bool
    Lotes           []Lote     // 1:N
}
```

### 5. Lote (`lote.go`)
```go
type Lote struct {
    ID               uint
    CorralID         uint       // FK Corral (OBLIGATORIO)
    Nombre           string
    CantidadLechones int
    FechaCreacion    time.Time
    Estado           string     // activo, cerrado, vendido
    FechaCierre      *time.Time
    MotivoCierre     *string
    Destetes         []Destete  // 1:N (muchos destetes en un lote)
}
```

### 6. Padrillo (`padrillo.go`)
```go
type Padrillo struct {
    ID                    uint
    GranjaID              uint       // FK Granja
    NumeroCaravana        string     // Único por granja
    Nombre                string
    Genetica              *string
    FechaUltimaVacunacion *time.Time
    Activo                bool
    FechaBaja             *time.Time
    MotivoBaja            *string    // muerte, venta
    Servicios             []Servicio // 1:N
}
```

### 7. Cerda (`cerda.go`)
```go
type Cerda struct {
    ID             uint
    GranjaID       uint       // FK Granja
    NumeroCaravana string     // Único por granja
    DetallePelaje  *string
    Genetica       *string
    Estado         string     // disponible, servicio, gestacion, cria
    Activo         bool
    FechaBaja      *time.Time
    MotivoBaja     *string    // muerte, venta
    Servicios      []Servicio // 1:N
    Partos         []Parto    // 1:N
    Destetes       []Destete  // 1:N
}
```

### 8. Servicio (`servicio.go`)
```go
type Servicio struct {
    ID                      uint
    CerdaID                 uint       // FK Cerda
    FechaServicio           time.Time
    TieneRepeticiones       bool
    CantidadRepeticiones    int
    TipoMonta               string     // natural, inseminacion
    // Monta Natural:
    PadrilloID              *uint      // FK Padrillo (si natural)
    CantidadSaltos          *int
    // Inseminación:
    NumeroPajuela           *string    // (si inseminacion)
    // Control de Preñez:
    PrenezConfirmada        bool
    FechaConfirmacionPrenez *time.Time
    PrenezCancelada         bool
    FechaCancelacionPrenez  *time.Time
    MotivoCancelacion       *string
    Partos                  []Parto    // 1:N
}

// Validación: Si natural → requiere PadrilloID
// Validación: Si inseminacion → requiere NumeroPajuela
```

### 9. Parto (`parto.go`)
```go
type Parto struct {
    ID                     uint
    CerdaID                uint       // FK Cerda
    ServicioID             *uint      // FK Servicio
    FechaParto             time.Time
    LechonesNacidosVivos   int
    LechonesNacidosTotales int
    LechonesHembras        int
    LechonesMachos         int
    FechaEstimada          time.Time  // fecha_servicio + 114 días
    Destetes               []Destete  // 1:N
}

// Validación: LechonesHembras + LechonesMachos == LechonesNacidosVivos
// Validación: LechonesNacidosTotales >= LechonesNacidosVivos
```

### 10. Destete (`destete.go`)
```go
type Destete struct {
    ID                         uint
    CerdaID                    uint       // FK Cerda
    PartoID                    *uint      // FK Parto
    FechaDestete               time.Time
    CantidadLechonesDestetados int
    FechaEstimada              time.Time  // fecha_parto + 30 días
    LoteID                     uint       // FK Lote (OBLIGATORIO)
}

// Validación: CantidadLechonesDestetados <= parto.LechonesNacidosVivos
```

## Características Implementadas

### ✅ Validaciones en Hooks
- **Servicio**: Valida campos según tipo de monta (natural/inseminación)
- **Parto**: Valida que suma de hembras y machos sea igual a vivos
- **Destete**: Valida que cantidad destetados no supere nacidos vivos

### ✅ Soft Deletes
Todos los modelos principales usan `gorm.DeletedAt`:
- Los registros no se eliminan físicamente
- Las queries automáticamente excluyen registros eliminados

### ✅ Timestamps Automáticos
- `CreatedAt`: Fecha de creación
- `UpdatedAt`: Fecha de última actualización

### ✅ Índices de Base de Datos
- Índices únicos: `username`, `email`, `numero_caravana` (por granja)
- Índices compuestos: `(cerda_id, fecha_servicio)`, etc.
- Índices para foreign keys y campos de búsqueda frecuente

### ✅ Tags de Validación
- `binding:"required"` - Campo obligatorio
- `binding:"email"` - Formato de email
- `binding:"min=X,max=Y"` - Longitud mínima/máxima
- `binding:"gte=0"` - Mayor o igual que
- `binding:"oneof=A B C"` - Valor debe ser uno de los especificados

### ✅ JSON Serialization
- Campos sensibles (`PasswordHash`) no se serializan: `json:"-"`
- Relaciones opcionales: `json:"...,omitempty"`

### ✅ Constantes para Enums
Definidas en `models.go`:
```go
// Estados de Cerda
EstadoCerdaDisponible = "disponible"
EstadoCerdaServicio   = "servicio"
EstadoCerdaGestacion  = "gestacion"
EstadoCerdaCria       = "cria"

// Estados de Lote
EstadoLoteActivo  = "activo"
EstadoLoteCerrado = "cerrado"
EstadoLoteVendido = "vendido"

// Tipos de Monta
TipoMontaNatural      = "natural"
TipoMontaInseminacion = "inseminacion"

// Y más...
```

## Uso

### Importar modelos
```go
import "github.com/martin/sgp/internal/models"

// Usar constantes
cerda := models.Cerda{
    Estado: models.EstadoCerdaDisponible,
}
```

### Migrar todos los modelos
```go
db.AutoMigrate(models.AllModels()...)
```

## Próximos Pasos

1. ⏳ Implementar **Repositories** - Capa de acceso a datos
2. ⏳ Implementar **Services** - Lógica de negocio (ciclo de maternidad, cálculos de fechas)
3. ⏳ Implementar **Handlers** - Controladores HTTP para cada entidad
4. ⏳ Testing - Unit tests para validaciones y reglas de negocio

---

**Estado**: ✅ Modelos GORM completamente implementados  
**Fecha**: 2025-02-07
