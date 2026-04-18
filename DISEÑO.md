# Documento de Diseño - Sistema de Gestión de Granja Porcina

## 1. Alcance del Sistema

### 1.1 Funcionalidades Principales

El sistema permite gestionar el ciclo completo de maternidad en una granja de cría porcina, específicamente en el sector de **cría y servicio**. Las funcionalidades incluyen:

- **Registro de Cerdas y Padrillos**: Gestión completa del inventario de animales
- **Registro de Servicios**: Monta natural e inseminación artificial
- **Control de Preñez**: Confirmación y cancelación de preñeces
- **Registro de Partos**: Información completa del parto y lechones
- **Registro de Destetes**: Finalización del ciclo de cría y creación de lotes
- **Registro de Muertes de Lechones**: Carga de bajas por parto/granja/corral y ajuste de cantidades
- **Registro de Ventas**: Venta parcial o total de lotes con impacto en stock
- **Gestión de Granjas**: Administración de múltiples granjas
- **Gestión de Corrales**: Organización de espacios físicos dentro de granjas
- **Gestión de Lotes**: Administración de grupos de lechones destetados
- **Listados y Consultas**: Visualización, modificación y eliminación de registros
- **Calendario**: Visualización de partos y destetes futuros
- **Estadísticas**: Análisis de la información recopilada

### 1.2 Restricciones de Eliminación

- **Cerdas y Padrillos**: Pueden eliminarse (baja lógica) solo por venta o muerte
- **Servicios, Partos y Destetes**: Solo pueden modificarse, NO eliminarse (auditoría)
- **Lotes**: Pueden cerrarse o marcarse como vendidos, pero no eliminarse (auditoría)

### 1.3 Estructura Jerárquica del Sistema

El sistema está organizado en una jerarquía de tres niveles:

```
┌─────────────────┐
│    GRANJAS      │ ← Nivel superior
└────────┬────────┘
         │
         ├─── Cerdas (pertenecen a una granja)
         ├─── Padrillos (pertenecen a una granja)
         │
         └─── ┌─────────────────┐
              │    CORRALES     │ ← Nivel intermedio
              └────────┬────────┘
                       │
                       └─── ┌─────────────────┐
                            │     LOTES       │ ← Nivel inferior
                            └─────────────────┘
```

**Relaciones:**
- **Granja** → tiene múltiples **Corrales** (1:N)
- **Granja** → tiene múltiples **Cerdas** y **Padrillos** (1:N)
- **Granja** ↔ **Usuario** mediante **UsuarioGranja** (N:M)
- **Corral** → contiene 0, 1 o más **Lotes** (1:N opcional)
- **Corral** → tiene múltiples **Muertes de Lechones** (1:N)
- **Destete** → genera automáticamente 1 **Lote** (1:1)
- **Lote** → pertenece a un **Corral** (N:1)
- **Lote** → agrupa múltiples **Destetes** (1:N)
- **Servicio** → pertenece a una **Cerda** y puede asociar un **Padrillo** (N:1)
- **Parto** → pertenece a un **Servicio** y a una **Cerda** (N:1)
- **Muertes de Lechones** → pueden asociarse a **Parto**, **Granja** y/o **Corral**
- **Ventas** → pueden asociarse a un **Lote** (venta parcial o total)

---

## 2. Ciclo de Maternidad

### 2.1 Estados de la Cerda

El sistema gestiona el ciclo de maternidad a través de 4 estados:

```
┌─────────────┐
│ Disponible  │ ← Estado inicial y final del ciclo
└──────┬──────┘
       │ [Registrar Servicio]
       ▼
┌─────────────┐
│  Servicio   │ ← Esperando confirmación de preñez
└──────┬──────┘
       │ [Confirmar Preñez]
       ▼
┌─────────────┐
│  Gestación  │ ← Cerda preñada (114 días promedio)
└──────┬──────┘
       │ [Registrar Parto]
       ▼
┌─────────────┐
│    Cría     │ ← Cerda con lechones (30 días promedio)
└──────┬──────┘
       │ [Registrar Destete]
       ▼
┌─────────────┐
│ Disponible  │ ← Vuelve al inicio del ciclo
└─────────────┘
```

### 2.2 Transiciones de Estado

| Estado Origen | Acción Requerida | Estado Destino | Validaciones |
|---------------|------------------|----------------|--------------|
| **Disponible** | Registrar Servicio | **Servicio** | Cerda debe estar activa y en estado "disponible" |
| **Servicio** | Confirmar Preñez | **Gestación** | Servicio debe existir y no estar cancelado |
| **Servicio** | Cancelar Preñez | **Disponible** | Se requiere motivo de cancelación |
| **Gestación** | Registrar Parto | **Cría** | Debe existir servicio con preñez confirmada |
| **Cría** | Registrar Destete | **Disponible** | Debe existir parto asociado. **Se crea automáticamente un Lote** |

### 2.3 Duración Estimada

- **Gestación**: 114 días (desde la fecha de servicio hasta parto)
- **Cría**: 30 días (desde parto hasta destete)
- **Ciclo completo**: ~144 días (desde servicio hasta nuevo servicio)

---

## 3. Reglas de Negocio

### 3.1 Gestión de Cerdas

#### 3.1.1 Registro de Cerda
- **Número de caravana**: Obligatorio, único dentro de la granja
- **Granja**: Obligatoria
- **Detalle de pelaje**: Opcional, texto libre
- **Genética**: Opcional, texto libre (ej: raza, línea genética)
- **Estado inicial**: Puede ser **cualquiera** (disponible, servicio, gestación, cría). Las granjas existentes tienen cerdas en todos los estados al momento de cargar datos.
- **Activo**: Por defecto `true`

#### 3.1.2 Modificación de Cerda
- Se puede modificar: número de caravana, detalle de pelaje, genética
- **NO se puede modificar directamente el estado**: El estado cambia automáticamente según las acciones del ciclo
- Si la cerda está en estado "servicio" o "gestación", no se puede modificar ciertos campos críticos

#### 3.1.3 Baja de Cerda
- Solo se puede dar de baja por: **muerte** o **venta**
- Al dar de baja:
  - `activo = false`
  - `fecha_baja = fecha actual`
  - `motivo_baja = 'muerte' | 'venta'`
- Si la cerda tiene servicios activos (en estado "servicio" o "gestación"), se debe cancelar primero

### 3.2 Gestión de Padrillos

#### 3.2.1 Registro de Padrillo
- **Granja**: Obligatoria
- **Número de caravana**: Obligatorio, único dentro de la granja
- **Nombre**: Obligatorio
- **Genética**: Opcional, texto libre (ej: raza, línea genética)
- **Fecha última vacunación**: Opcional

#### 3.2.2 Baja de Padrillo
- Mismas reglas que cerdas: solo por muerte o venta
- Si el padrillo tiene servicios asociados, estos se mantienen (histórico)

### 3.3 Gestión de Servicios

#### 3.3.1 Registro de Servicio
- **Cerda**: Obligatoria, debe estar en estado "disponible"
- **Fecha de servicio**: Obligatoria
- **Tipo de monta**: Obligatorio (`natural` o `inseminacion`)

**Si es Monta Natural:**
- **Padrillo**: Obligatorio
- **Cantidad de saltos**: Opcional

**Si es Inseminación:**
- **Número de pajuela**: Obligatorio

**Repeticiones:**
- `tiene_repeticiones`: Boolean
- `cantidad_repeticiones`: Si tiene repeticiones, cantidad (opcional)

**Efecto:**
- Cambia el estado de la cerda a **"servicio"**

#### 3.3.2 Confirmación de Preñez
- Solo se puede confirmar si:
  - El servicio existe
  - La cerda está en estado "servicio"
  - La preñez no está cancelada
- Al confirmar:
  - `prenez_confirmada = true`
  - `fecha_confirmacion_prenez = fecha actual`
  - Estado de la cerda cambia a **"gestación"**
  - Se calcula `fecha_estimada_parto = fecha_servicio + 114 días`

#### 3.3.3 Cancelación de Preñez
- Solo se puede cancelar si:
  - El servicio existe
  - La preñez está confirmada
  - La preñez no está previamente cancelada
- Al cancelar:
  - `prenez_cancelada = true`
  - `fecha_cancelacion_prenez = fecha actual`
  - `motivo_cancelacion = texto obligatorio`
  - Estado de la cerda vuelve a **"disponible"**

#### 3.3.4 Modificación de Servicio
- Se puede modificar antes de confirmar preñez
- Una vez confirmada la preñez, solo se pueden modificar campos no críticos
- No se puede eliminar un servicio (auditoría)

### 3.4 Gestión de Partos

#### 3.4.1 Registro de Parto
- **Flujo**: Se muestran las cerdas en estado "gestación"; el usuario elige una cerda para registrar el parto (más intuitivo que elegir servicio/parto primero).
- **Cerda**: Obligatoria, debe estar en estado "gestación". El sistema obtiene el servicio con preñez confirmada asociado a esa cerda (si hay más de uno, se puede elegir o tomar el más reciente).
- **Servicio**: Se asocia al parto (obtenido por la cerda seleccionada).
- **Fecha de parto**: Obligatoria
- **Lechones nacidos vivos**: Obligatorio, >= 0
- **Lechones nacidos totales**: Obligatorio, >= lechones vivos
- **Lechones hembras**: Obligatorio, >= 0
- **Lechones machos**: Obligatorio, >= 0
- **Validación**: `lechones_hembras + lechones_machos == lechones_nacidos_vivos` (la suma de hembras y machos debe ser exactamente igual a los nacidos vivos)

**Efecto:**
- Cambia el estado de la cerda a **"cría"**
- Se calcula `fecha_estimada_destete = fecha_parto + 30 días`

#### 3.4.2 Modificación de Parto
- Se puede modificar la información del parto
- No se puede eliminar (auditoría)

### 3.5 Gestión de Destetes

#### 3.5.1 Registro de Destete
- **Flujo**: Se muestran las cerdas en estado "cría"; el usuario elige una cerda para registrar el destete. No es necesario mostrar ni elegir el parto: el sistema infiere el parto sin destete asociado a esa cerda (más intuitivo).
- **Cerda**: Obligatoria, debe estar en estado "cría". El sistema obtiene automáticamente el parto sin destete de esa cerda.
- **Parto**: Se asocia automáticamente (un parto por cerda en cría sin destete).
- **Fecha de destete**: Obligatoria
- **Cantidad de lechones destetados**: Obligatorio, >= 0
- **Validación**: `cantidad_lechones_destetados <= lechones_nacidos_vivos del parto`

**Efecto:**
- Cambia el estado de la cerda a **"disponible"**
- Completa el ciclo de maternidad
- **Los lechones destetados DEBEN asignarse a un Lote** (obligatorio):
  - El lote puede existir previamente o crearse en el momento
  - Si se crea nuevo lote: requiere nombre, corral (obligatorio) y fecha de creación
  - Si se asigna a lote existente: se suma la cantidad de lechones al lote
  - Un lote puede contener lechones de múltiples destetes

#### 3.5.2 Modificación de Destete
- Se puede modificar la información del destete
- No se puede eliminar (auditoría)
- **Nota**: Si se modifica la cantidad de lechones destetados, se debe actualizar también la cantidad del lote asociado

### 3.6 Gestión de Granjas

#### 3.6.1 Registro de Granja
- **Nombre**: Obligatorio, único en el sistema
- **Descripción**: Opcional
- **Ubicación**: Opcional
- **Activo**: Por defecto `true`

#### 3.6.2 Asignación de Usuarios a Granjas
- Un usuario puede gestionar múltiples granjas (relación N:M)
- Roles en granja:
  - **Propietario**: Control total sobre la granja
  - **Administrador**: Puede gestionar todo excepto eliminar la granja
  - **Operador**: Puede registrar y modificar datos, pero no eliminar

#### 3.6.3 Modificación y Baja de Granja
- Se puede modificar nombre, descripción y ubicación
- Solo se puede dar de baja si no tiene cerdas, padrillos o corrales activos
- Al dar de baja: `activo = false`

### 3.7 Gestión de Corrales

#### 3.7.1 Registro de Corral
- **Granja**: Obligatoria
- **Nombre**: Obligatorio, único dentro de la granja
- **Descripción**: Opcional
- **Capacidad máxima**: Opcional (campo libre, no se valida automáticamente)
- **Activo**: Por defecto `true`
- **Nota**: Un corral puede crearse vacío (sin lotes)

#### 3.7.2 Modificación de Corral
- Se puede modificar todos los campos
- Se puede editar la cantidad de animales por corral (suma de lotes)
- **Nota**: No se valida automáticamente la capacidad máxima, es solo informativo

#### 3.7.3 Baja de Corral
- Solo se puede dar de baja si no tiene lotes activos
- Al dar de baja: `activo = false`
- **Importante**: Los lotes NO pueden quedar sin corral, por lo que se debe reasignar o cerrar los lotes antes de dar de baja el corral

### 3.8 Gestión de Lotes

#### 3.8.1 Creación de Lote
- **Manual**: Se puede crear un lote antes de registrar destetes
- **En el momento del destete**: Se puede crear un lote nuevo o asignar a uno existente
- **Datos iniciales**:
  - Nombre: Obligatorio, ingresado por usuario
  - Corral: **Obligatorio** (el lote DEBE estar asignado a un corral)
  - Cantidad de lechones: Inicialmente 0, se actualiza al asignar destetes
  - Fecha de creación: Fecha ingresada por usuario
  - Estado: "activo"

#### 3.8.2 Asignación de Destetes a Lote
- **Obligatorio**: Los lechones destetados DEBEN pertenecer a un lote
- Un lote puede contener lechones de **múltiples destetes**
- Al asignar un destete a un lote:
  - Se suma `cantidad_lechones_destetados` a `cantidad_lechones` del lote
  - El destete queda vinculado al lote (relación N:1)
- **Restricción**: Un lote NO puede moverse entre corrales una vez asignado
  - En la práctica, una vez que los cerdos están juntos en un corral, no es viable detectar qué cerdo pertenece a qué lote

#### 3.8.3 Estados del Lote
- **activo**: Lote en producción, lechones creciendo
- **cerrado**: Lote finalizado (por ejemplo, lechones movidos a otra etapa)
- **vendido**: Lote vendido

#### 3.8.4 Cierre de Lote
- Se puede cerrar un lote cambiando su estado
- Al cerrar:
  - `estado = 'cerrado' | 'vendido'`
  - `fecha_cierre = fecha actual`
  - `motivo_cierre = texto obligatorio`

#### 3.8.5 Modificación de Lote
- Se puede modificar: nombre, cantidad de lechones
- **NO se puede modificar el corral**: Una vez asignado, el lote no puede moverse entre corrales
- **Cantidad de lechones**: Se puede editar manualmente (por ejemplo, por muertes, ventas parciales, etc.)
  - La cantidad es la suma de todos los destetes asignados, pero puede ajustarse manualmente
- **Lechones individuales**: No se gestionan como entidad individual, solo cantidad
- No se puede eliminar (auditoría)

---

## 4. Diseño de Interfaces

### 4.1 Estructura de Navegación

```
┌─────────────────────────────────────┐
│  Header: Logo + Usuario + Logout    │
├─────────────────────────────────────┤
│  Menú Principal (Sidebar o Top Nav) │
│  - Dashboard                         │
│  - Granjas                           │
│  - Corrales                           │
│  - Lotes                              │
│  - Cerdas                            │
│  - Padrillos                         │
│  - Servicios                         │
│  - Partos                            │
│  - Destetes                          │
│  - Calendario                        │
│  - Estadísticas                      │
└─────────────────────────────────────┘
```

### 4.2 Páginas Principales

#### 4.2.1 Dashboard
- **Resumen general**:
  - Total de cerdas por estado
  - Total de padrillos activos
  - Servicios pendientes de confirmación
  - Partos próximos (próximos 7 días)
  - Destetes próximos (próximos 7 días)
- **Gráficos rápidos**:
  - Distribución de cerdas por estado
  - Servicios del mes actual

#### 4.2.2 Gestión de Granjas
- **Listado de granjas**:
  - Tabla con: Nombre, Ubicación, Cantidad de corrales, Cantidad de cerdas, Acciones
  - Filtros: Activas/inactivas
  - Búsqueda: Por nombre
- **Formulario de alta**:
  - Nombre (requerido, único)
  - Descripción (opcional)
  - Ubicación (opcional)
- **Formulario de edición**:
  - Mismos campos que alta
- **Gestión de usuarios**:
  - Asignar usuarios a la granja con roles
  - Ver usuarios asignados
- **Acciones**:
  - Ver detalle completo
  - Editar
  - Dar de baja (solo si no tiene datos asociados)
  - Ver corrales, cerdas, padrillos

#### 4.2.3 Gestión de Corrales
- **Listado de corrales**:
  - Tabla con: Nombre, Granja, Capacidad máxima, Lotes activos, Ocupación, Acciones
  - Filtros: Por granja, activos/inactivos
  - Búsqueda: Por nombre
- **Formulario de alta**:
  - Selector de granja (requerido)
  - Nombre (requerido, único en la granja)
  - Descripción (opcional)
  - Capacidad máxima (opcional)
- **Formulario de edición**:
  - Mismos campos que alta
- **Acciones**:
  - Ver detalle completo
  - Editar
  - Dar de baja (solo si no tiene lotes activos)
  - Ver lotes asignados

#### 4.2.4 Gestión de Lotes
- **Listado de lotes**:
  - Tabla con: Nombre, Corral, Cantidad lechones, Fecha creación, Estado, Acciones
  - Filtros: Por corral, por estado, por granja
  - Búsqueda: Por nombre
- **Vista de detalle**:
  - Información del lote
  - Lista de destetes asociados (puede haber múltiples)
  - Suma total de lechones de todos los destetes
- **Formulario de alta**:
  - Nombre (obligatorio)
  - Selector de corral (obligatorio)
  - Fecha de creación
  - Cantidad inicial (opcional, se actualiza al asignar destetes)
- **Acciones**:
  - Modificar cantidad de lechones (edición manual)
  - Cerrar lote (con motivo)
  - Marcar como vendido
  - Ver destetes asociados
  - **NO**: Cambiar de corral (no permitido)

#### 4.2.5 Gestión de Cerdas
- **Listado de cerdas**:
  - Tabla con: Caravana, Granja, Genética, Estado, Último servicio, Último parto, Acciones
  - Filtros: Por estado, por granja, activas/inactivas
  - Búsqueda: Por número de caravana
- **Formulario de alta**:
  - Selector de granja (requerido)
  - Número de caravana (requerido, único en la granja)
  - Detalle de pelaje (opcional)
  - Genética (opcional, string)
  - Estado inicial (requerido: disponible, servicio, gestación o cría)
- **Formulario de edición**:
  - Mismos campos editables que alta (caravana, pelaje, genética)
  - No permite modificar estado directamente (el estado cambia por el ciclo)
- **Acciones**:
  - Ver detalle completo
  - Editar
  - Dar de baja (muerte/venta)
  - Ver historial (servicios, partos, destetes)

#### 4.2.6 Gestión de Padrillos
- Listado: Caravana, Nombre, Genética, Última vacunación, Acciones
- Formulario de alta: Granja, Caravana, Nombre, Genética (opcional), Fecha última vacunación
- Formulario de edición: Mismos campos editables

#### 4.2.7 Gestión de Servicios
- **Listado de servicios**:
  - Tabla con: Fecha, Cerda, Tipo monta, Padrillo/Pajuela, Estado preñez, Acciones
  - Filtros: Por mes/año, por cerda, por estado de preñez
- **Formulario de registro**:
  - Selector de cerda (solo disponibles)
  - Fecha de servicio
  - Tipo de monta (radio buttons)
  - Si natural: selector de padrillo, cantidad de saltos
  - Si inseminación: número de pajuela
  - Repeticiones: checkbox + cantidad
- **Acciones**:
  - Confirmar preñez (si está en estado "servicio")
  - Cancelar preñez (si está confirmada)
  - Editar (con restricciones)

#### 4.2.8 Gestión de Partos
- **Listado de partos**:
  - Tabla con: Fecha, Cerda, Lechones vivos, Lechones totales, Acciones
  - Filtros: Por mes/año, por cerda
- **Formulario de registro** (flujo intuitivo):
  - Paso 1: Se listan las **cerdas en gestación**; el usuario elige una cerda
  - Paso 2: Con la cerda elegida, se cargan fecha de parto y datos de lechones
  - Fecha de parto
  - Lechones nacidos vivos
  - Lechones nacidos totales
  - Lechones hembras, Lechones machos (validación: hembras + machos == vivos)
  - El servicio se asocia automáticamente (prenez confirmada de esa cerda)
- **Acciones**: Editar

#### 4.2.9 Gestión de Destetes
- **Listado de destetes**:
  - Tabla con: Fecha, Cerda, Lechones destetados, Lote, Acciones
  - Filtros: Por mes/año, por cerda, por granja
- **Formulario de registro** (flujo intuitivo):
  - Paso 1: Se listan las **cerdas en cría**; el usuario elige una cerda (no se muestran partos)
  - Paso 2: Con la cerda elegida, el sistema usa el parto sin destete de esa cerda
  - Fecha de destete
  - Cantidad de lechones destetados
  - **Asignación a lote** (obligatorio):
    - Opción 1: Seleccionar lote existente
    - Opción 2: Crear nuevo lote (nombre y corral obligatorio)
- **Acciones**: 
  - Editar
  - Ver lote asociado
  - Ver otros destetes del mismo lote

#### 4.2.10 Calendario
- **Vista mensual**:
  - Muestra partos y destetes futuros
  - Color diferente para cada tipo de evento
  - Al hacer clic: muestra detalles
- **Filtros**: Por tipo de evento, por rango de fechas

#### 4.2.11 Estadísticas
- **Métricas principales**:
  - Total de servicios por mes
  - Tasa de preñez confirmada
  - Total de partos por mes
  - Promedio de lechones por parto
  - Total de destetes por mes
  - Promedio de lechones destetados
- **Gráficos**:
  - Evolución mensual de servicios/partos/destetes
  - Distribución de lechones por parto
  - Tasa de éxito por padrillo (si aplica)

---

## 5. Flujos de Usuario

### 5.1 Flujo: Registro de Servicio Completo

```
1. Usuario accede a "Servicios" → "Nuevo Servicio"
2. Selecciona cerda (solo disponibles)
3. Ingresa fecha de servicio
4. Selecciona tipo de monta:
   - Si Natural: selecciona padrillo, ingresa saltos
   - Si Inseminación: ingresa número de pajuela
5. Opcionalmente marca "tiene repeticiones" y cantidad
6. Guarda → Sistema cambia estado cerda a "servicio"
7. Usuario ve servicio en listado con estado "Pendiente confirmación"
```

### 5.2 Flujo: Confirmación de Preñez

```
1. Usuario ve servicio en estado "Pendiente confirmación"
2. Hace clic en "Confirmar Preñez"
3. Sistema muestra confirmación
4. Usuario confirma → Sistema:
   - Marca preñez como confirmada
   - Cambia estado cerda a "gestación"
   - Calcula fecha estimada de parto
5. Usuario ve servicio con estado "Preñez confirmada"
6. En calendario aparece parto estimado
```

### 5.3 Flujo: Registro de Parto

```
1. Usuario accede a "Partos" → "Nuevo Parto"
2. Sistema muestra listado de **cerdas en gestación**
3. Usuario elige una cerda
4. Sistema asocia el servicio con preñez confirmada de esa cerda (o permite elegir si hay más de uno)
5. Usuario ingresa fecha de parto
6. Usuario ingresa datos de lechones:
   - Nacidos vivos
   - Nacidos totales
   - Hembras
   - Machos
7. Sistema valida: hembras + machos == vivos (exactamente igual)
8. Guarda → Sistema:
   - Cambia estado cerda a "cría"
   - Calcula fecha estimada de destete
9. En calendario aparece destete estimado
```

### 5.4 Flujo: Registro de Destete y Asignación a Lote

```
1. Usuario accede a "Destetes" → "Nuevo Destete"
2. Sistema muestra listado de **cerdas en cría** (no se muestran partos)
3. Usuario elige una cerda
4. Sistema obtiene automáticamente el parto sin destete de esa cerda
5. Usuario ingresa fecha de destete
6. Usuario ingresa cantidad de lechones destetados
7. Sistema valida: destetados <= nacidos vivos del parto
8. Sistema muestra opciones de asignación a lote:
   Opción A - Lote existente: selector de lote (filtrado por granja/corral)
   Opción B - Crear nuevo lote: nombre y corral (obligatorios), fecha de creación
9. Usuario selecciona opción y completa datos
10. Guarda → Sistema:
    - Cambia estado cerda a "disponible"
    - Completa el ciclo
    - Si lote nuevo: crea lote con datos ingresados
    - Si lote existente: suma cantidad al lote
    - Asocia destete al lote
    - Actualiza cantidad_lechones del lote
11. Cerda vuelve a estar disponible para nuevo servicio
12. Lote queda disponible para recibir más destetes
```

### 5.5 Flujo: Creación de Lote Anticipado

```
1. Usuario accede a "Lotes" → "Nuevo Lote"
2. Ingresa nombre del lote (obligatorio)
3. Selecciona corral (obligatorio, de la granja actual)
4. Ingresa fecha de creación
5. Cantidad inicial: 0 (se actualizará al asignar destetes)
6. Guarda → Lote creado y disponible para asignar destetes
7. El lote aparece en listado, listo para recibir lechones destetados
```

### 5.6 Flujo: Gestión de Corrales

```
1. Usuario accede a "Corrales" → "Nuevo Corral"
2. Selecciona granja
3. Ingresa nombre (único en la granja)
4. Opcionalmente ingresa descripción y capacidad máxima
5. Guarda → Corral creado
6. Corral aparece en listado, disponible para asignar lotes
```

### 5.7 Flujo: Cancelación de Preñez

```
1. Usuario ve servicio con preñez confirmada
2. Hace clic en "Cancelar Preñez"
3. Sistema muestra formulario con campo "Motivo"
4. Usuario ingresa motivo (obligatorio)
5. Confirma → Sistema:
   - Marca preñez como cancelada
   - Cambia estado cerda a "disponible"
6. Cerda vuelve a estar disponible para nuevo servicio
```

---

## 6. Validaciones y Restricciones

### 6.1 Validaciones de Frontend
- Campos requeridos marcados con *
- Validación de formato de fechas
- Validación de números (>= 0)
- Validación de relaciones lógicas (hembras + machos == vivos en partos)
- Mensajes de error claros y específicos

### 6.2 Validaciones de Backend
- Todas las validaciones de frontend se repiten en backend
- Validación de existencia de registros relacionados
- Validación de estados antes de transiciones
- Validación de unicidad (caravanas)
- Validación de integridad referencial

### 6.3 Restricciones de Negocio
- No se puede registrar servicio si cerda no está disponible
- No se puede confirmar preñez si servicio no existe o está cancelado
- No se puede registrar parto si cerda no está en gestación
- No se puede registrar destete si cerda no está en cría
- No se pueden eliminar servicios, partos o destetes (solo modificar)

---

## 7. Plan de Desarrollo

### Fase 1: Infraestructura Base (Prioridad Alta)
- [x] Diagrama ER
- [x] Schema SQL
- [ ] Modelos GORM
- [ ] Middleware de autenticación
- [ ] Middleware de roles
- [ ] Configuración de base de datos

### Fase 2: Backend - CRUD Básico (Prioridad Alta)
- [ ] Controladores y rutas de Cerdas
- [ ] Controladores y rutas de Padrillos
- [ ] Controladores y rutas de Servicios
- [ ] Controladores y rutas de Partos
- [ ] Controladores y rutas de Destetes
- [ ] Validaciones de request en handlers/services

### Fase 3: Backend - Lógica de Negocio (Prioridad Alta)
- [ ] Transiciones de estado automáticas
- [ ] Confirmación y cancelación de preñez
- [ ] Cálculo de fechas estimadas
- [ ] Validaciones de ciclo de maternidad

### Fase 4: Frontend - Interfaces Básicas (Prioridad Media)
- [ ] Página de login
- [ ] Dashboard
- [ ] Gestión de Cerdas (listado, alta, edición)
- [ ] Gestión de Padrillos (listado, alta, edición)

### Fase 5: Frontend - Ciclo de Maternidad (Prioridad Media)
- [ ] Gestión de Servicios
- [ ] Confirmación/Cancelación de Preñez
- [ ] Gestión de Partos
- [ ] Gestión de Destetes

### Fase 6: Funcionalidades Avanzadas (Prioridad Baja)
- [ ] Calendario de eventos futuros
- [ ] Estadísticas y gráficos
- [ ] Exportación a Excel
- [ ] Reportes

### Fase 7: Mejoras y Optimización (Prioridad Baja)
- [ ] Caché de estadísticas
- [ ] Optimización de consultas
- [ ] Mejoras de UI/UX
- [ ] Testing

---

## 8. Consideraciones Técnicas

### 8.1 Manejo de Estados
- Los cambios de estado deben ser atómicos (transacciones)
- Registrar en logs todos los cambios de estado
- Permitir rollback en caso de error

### 8.2 Fechas Estimadas
- Se calculan automáticamente al crear/confirmar
- Se pueden ajustar manualmente si es necesario
- Se muestran en calendario con indicador de "estimada"

### 8.3 Auditoría
- Todos los registros tienen `created_at` y `updated_at`
- No se eliminan físicamente servicios, partos o destetes
- Se mantiene historial completo de cada cerda

### 8.4 Rendimiento
- Índices en campos de búsqueda frecuente
- Caché de estadísticas para consultas pesadas
- Paginación en listados grandes

*Documento creado: 2025-01-XX*
*Última actualización: 2026-04-18*
