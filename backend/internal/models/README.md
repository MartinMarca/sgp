# Modelos GORM - Sistema de Gestión de Granja Porcina

Este directorio contiene todos los modelos GORM que representan las entidades del sistema.

## Modelos Implementados

### 1. Usuario (`usuario.go`)
Representa un usuario del sistema con autenticación.

**Campos principales:**
- `Username`, `Email`, `PasswordHash` (autenticación)
- `NombreCompleto`, `Rol`, `Activo`
- Roles: `admin`, `usuario`, `veterinario`

**Relaciones:**
- Many-to-Many con `Granja` (a través de `usuario_granja`)

---

### 2. Granja (`granja.go`)
Representa una granja porcina.

**Campos principales:**
- `Nombre`, `Descripcion`, `Ubicacion`, `Activo`

**Relaciones:**
- Many-to-Many con `Usuario`
- Has-Many: `Corrales`, `Cerdas`, `Padrillos`

---

### 3. UsuarioGranja (`usuario_granja.go`)
Tabla intermedia para la relación N:M entre Usuarios y Granjas.

**Campos principales:**
- `UsuarioID`, `GranjaID`, `Rol`
- Roles en granja: `propietario`, `administrador`, `operador`

---

### 4. Corral (`corral.go`)
Representa un corral dentro de una granja.

**Campos principales:**
- `GranjaID`, `Nombre`, `Descripcion`, `CapacidadMaxima`, `Activo`

**Relaciones:**
- Belongs-To: `Granja`
- Has-Many: `Lotes`

---

### 5. Lote (`lote.go`)
Representa un lote de lechones destetados.

**Campos principales:**
- `CorralID` (obligatorio), `Nombre`, `CantidadLechones`
- `FechaCreacion`, `Estado`, `FechaCierre`, `MotivoCierre`
- Estados: `activo`, `cerrado`, `vendido`

**Relaciones:**
- Belongs-To: `Corral`
- Has-Many: `Destetes` (muchos destetes pueden estar en un lote)

---

### 6. Padrillo (`padrillo.go`)
Representa un padrillo (cerdo macho reproductor).

**Campos principales:**
- `GranjaID`, `NumeroCaravana` (único por granja), `Nombre`, `Genetica`
- `FechaUltimaVacunacion`, `Activo`, `FechaBaja`, `MotivoBaja`
- Motivos de baja: `muerte`, `venta`

**Relaciones:**
- Belongs-To: `Granja`
- Has-Many: `Servicios`

**Validación:**
- `NumeroCaravana` es único dentro de cada granja

---

### 7. Cerda (`cerda.go`)
Representa una cerda (cerda hembra reproductora).

**Campos principales:**
- `GranjaID`, `NumeroCaravana` (único por granja), `DetallePelaje`, `Genetica`
- `Estado`, `Activo`, `FechaBaja`, `MotivoBaja`
- Estados: `disponible`, `servicio`, `gestacion`, `cria`

**Relaciones:**
- Belongs-To: `Granja`
- Has-Many: `Servicios`, `Partos`, `Destetes`

**Validación:**
- `NumeroCaravana` es único dentro de cada granja

**Ciclo de Maternidad:**
```
Disponible → Servicio → Gestación → Cría → Disponible
```

---

### 8. Servicio (`servicio.go`)
Representa un servicio de monta (natural o inseminación).

**Campos principales:**
- `CerdaID`, `FechaServicio`, `TipoMonta`
- `TieneRepeticiones`, `CantidadRepeticiones`
- **Monta Natural**: `PadrilloID`, `CantidadSaltos`
- **Inseminación**: `NumeroPajuela`
- **Control de Preñez**: `PrenezConfirmada`, `FechaConfirmacionPrenez`, `PrenezCancelada`, `FechaCancelacionPrenez`, `MotivoCancelacion`

**Relaciones:**
- Belongs-To: `Cerda`, `Padrillo` (opcional)
- Has-Many: `Partos`

**Validaciones:**
- Si `TipoMonta == "natural"` → `PadrilloID` es obligatorio
- Si `TipoMonta == "inseminacion"` → `NumeroPajuela` es obligatorio

---

### 9. Parto (`parto.go`)
Representa un parto de una cerda.

**Campos principales:**
- `CerdaID`, `ServicioID`, `FechaParto`, `FechaEstimada`
- `LechonesNacidosVivos`, `LechonesNacidosTotales`
- `LechonesHembras`, `LechonesMachos`

**Relaciones:**
- Belongs-To: `Cerda`, `Servicio` (opcional)
- Has-Many: `Destetes`

**Validaciones:**
- `LechonesHembras + LechonesMachos == LechonesNacidosVivos`
- `LechonesNacidosTotales >= LechonesNacidosVivos`

---

### 10. Destete (`destete.go`)
Representa un destete de una cerda.

**Campos principales:**
- `CerdaID`, `PartoID`, `FechaDestete`, `FechaEstimada`
- `CantidadLechonesDestetados`
- `LoteID` (obligatorio: los lechones deben estar en un lote)

**Relaciones:**
- Belongs-To: `Cerda`, `Parto` (opcional), `Lote`

**Validaciones:**
- `CantidadLechonesDestetados <= parto.LechonesNacidosVivos`

---

## Características de los Modelos

### Soft Deletes
Todos los modelos principales implementan soft deletes con `gorm.DeletedAt`:
- Los registros no se eliminan físicamente de la base de datos
- Se marcan con una fecha de eliminación
- Las queries automáticamente excluyen registros eliminados

### Timestamps Automáticos
Todos los modelos tienen:
- `CreatedAt`: Fecha de creación (automático)
- `UpdatedAt`: Fecha de última actualización (automático)

### Validaciones en Hooks
Algunos modelos implementan validaciones en hooks de GORM:
- `BeforeCreate`: Validaciones antes de crear
- `BeforeUpdate`: Validaciones antes de actualizar

Modelos con validaciones:
- `Servicio`: Valida campos según tipo de monta
- `Parto`: Valida suma de lechones
- `Destete`: Valida cantidad destetados vs parto

### Tags de Binding
Los modelos usan tags de `binding` para validación con `gin`:
- `required`: Campo obligatorio
- `min`, `max`: Longitud mínima/máxima
- `gte`, `lte`: Mayor/menor que igual
- `email`: Formato de email
- `oneof`: Valor debe ser uno de los especificados

## Uso

### Importar todos los modelos
```go
import "github.com/martin/sgp/internal/models"

// Usar constantes
estado := models.EstadoCerdaDisponible

// Crear instancia
cerda := models.Cerda{
    GranjaID: 1,
    NumeroCaravana: "C-001",
    Estado: models.EstadoCerdaDisponible,
}
```

### Migraciones
```go
// Migrar todos los modelos
db.AutoMigrate(models.AllModels()...)

// O usar el helper
database.AutoMigrate(db)
```

## Notas Importantes

1. **Orden de creación**: Las tablas deben crearse en el orden correcto debido a las foreign keys. GORM maneja esto automáticamente.

2. **Lotes antes de Destetes**: La tabla `lotes` debe crearse antes que `destetes` porque hay una relación circular (`destetes` → `lotes`, `lotes` ← `destetes`).

3. **Enums en MySQL**: Los tipos ENUM se definen como strings en Go, pero GORM los mapea correctamente a ENUM de MySQL.

4. **Fechas NULL**: Se usa `sql.NullTime` para fechas que pueden ser NULL en la base de datos.

5. **Punteros para opcionales**: Los campos opcionales usan punteros (`*string`, `*int`, etc.) para distinguir entre "no proporcionado" y "valor cero".
