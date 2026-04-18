# Arquitectura del Sistema de Gestión de Granja Porcina

## Stack Tecnológico

### Backend
- **Lenguaje**: Go 1.21+
- **Framework**: Gin (Web Framework)
- **Base de Datos**: MySQL 8.0+
- **ORM**: GORM (Go Object-Relational Mapping)
- **Autenticación**: golang-jwt/jwt
- **Validación**: go-playground/validator/v10
- **Variables de entorno**: godotenv
- **Exportación**: excelize (para Excel)
- **Migrations**: golang-migrate o GORM AutoMigrate

### Frontend
- **HTML5/CSS3/JavaScript** (Vanilla JS)
- **Fetch API** para comunicación con backend
- **Librerías opcionales**: 
  - Chart.js (para gráficos de estadísticas)
  - FullCalendar (para calendario de partos/destetes)

## Estructura de Carpetas

```
proyecto/
├── backend/                      # Backend en Go
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # Punto de entrada de la aplicación
│   ├── internal/                # Código privado de la aplicación
│   │   ├── config/
│   │   │   └── config.go        # Configuración y variables de entorno
│   │   ├── models/              # Modelos GORM
│   │   │   ├── usuario.go
│   │   │   ├── granja.go
│   │   │   ├── corral.go
│   │   │   ├── lote.go
│   │   │   ├── cerda.go
│   │   │   ├── padrillo.go
│   │   │   ├── servicio.go
│   │   │   ├── parto.go
│   │   │   └── destete.go
│   │   ├── handlers/            # Controladores HTTP (Gin handlers)
│   │   │   ├── auth_handler.go
│   │   │   ├── granja_handler.go
│   │   │   ├── corral_handler.go
│   │   │   ├── lote_handler.go
│   │   │   ├── cerda_handler.go
│   │   │   ├── padrillo_handler.go
│   │   │   ├── servicio_handler.go
│   │   │   ├── parto_handler.go
│   │   │   └── destete_handler.go
│   │   ├── services/            # Lógica de negocio
│   │   │   ├── cerda_service.go
│   │   │   ├── parto_service.go
│   │   │   ├── destete_service.go
│   │   │   ├── estadisticas_service.go
│   │   │   ├── calendario_service.go
│   │   │   └── reporte_service.go
│   │   ├── repositories/        # Capa de acceso a datos
│   │   │   ├── cerda_repository.go
│   │   │   ├── parto_repository.go
│   │   │   └── ...
│   │   ├── middleware/          # Middlewares
│   │   │   ├── auth.go          # Verificación JWT
│   │   │   ├── cors.go          # CORS
│   │   │   └── roles.go         # Control de roles
│   │   ├── routes/              # Definición de rutas
│   │   │   └── routes.go
│   │   ├── database/            # Conexión y migraciones
│   │   │   ├── database.go
│   │   │   └── migrations.go
│   │   └── utils/               # Utilidades
│   │       ├── jwt.go           # Manejo de JWT
│   │       ├── response.go      # Respuestas estándar
│   │       └── validators.go    # Validadores personalizados
│   ├── migrations/              # Archivos de migración SQL (opcional)
│   ├── .env.example             # Ejemplo de variables de entorno
│   ├── go.mod                   # Dependencias Go
│   ├── go.sum                   # Checksums de dependencias
│   └── Makefile                 # Comandos útiles (build, run, test)
│
├── frontend/                     # (archivos en raíz del proyecto)
│   ├── index.html
│   ├── cerdas.html
│   ├── css/
│   ├── js/
│   │   ├── app.js
│   │   ├── api.js               # Cliente API
│   │   ├── auth.js              # Manejo de autenticación
│   │   └── modules/             # Módulos por funcionalidad
│   │       ├── cerdas.js
│   │       ├── servicios.js
│   │       └── ...
│   └── assets/
│
├── database/
│   ├── schema.sql               # Esquema SQL inicial
│   ├── diagrama-er.puml         # Diagrama ER
│   └── migration_lotes_destetes.sql
│
└── README.md
```

## Arquitectura por Capas

### 1. Capa de Presentación (Frontend)
- **Responsabilidad**: Interfaz de usuario, validación de formularios, visualización
- **Tecnología**: HTML, CSS, JavaScript vanilla

### 2. Capa de API (Backend - Routes/Handlers)
- **Responsabilidad**: Endpoints REST, validación de entrada, respuestas HTTP
- **Tecnología**: Gin Framework (Go)

### 3. Capa de Servicios (Backend - Services)
- **Responsabilidad**: Lógica de negocio compleja (estadísticas, cálculos, reglas de ciclo de maternidad)
- **Tecnología**: Go

### 4. Capa de Repositorios (Backend - Repositories)
- **Responsabilidad**: Abstracción de acceso a datos, queries complejas
- **Tecnología**: Go + GORM

### 5. Capa de Datos (Backend - Models)
- **Responsabilidad**: Definición de estructuras, relaciones, validaciones de modelo
- **Tecnología**: GORM (Go ORM)

## Flujo de Datos

```
Frontend (JS) 
  → API Client (fetch)
    → Gin Router
      → Middleware (auth/cors/roles)
        → Handler (controller)
          → Service (lógica de negocio)
            → Repository (queries)
              → GORM Models
                → MySQL Database
```

## Patrón de Diseño

El backend sigue una **arquitectura hexagonal simplificada** (Clean Architecture):

- **Handlers**: Capa de presentación HTTP, maneja requests/responses
- **Services**: Capa de aplicación, contiene lógica de negocio y reglas del dominio
- **Repositories**: Capa de infraestructura, abstrae el acceso a datos
- **Models**: Entidades del dominio, estructuras de datos

Beneficios:
- Separación de responsabilidades
- Fácil testing (mocking de repositories)
- Independencia de frameworks
- Escalabilidad y mantenibilidad

## Seguridad

- **Autenticación**: JWT tokens almacenados en localStorage
- **Autorización**: Middleware de roles (Propietario, Administrador, Operador)
- **Validación**: Backend (validator) y frontend
- **SQL Injection**: Prevenido por ORM (GORM con prepared statements)
- **XSS**: Sanitización de inputs
- **CORS**: Configuración restrictiva
- **Rate Limiting**: Middleware para prevenir abuso (opcional)

## Escalabilidad

- **Concurrencia**: Goroutines para operaciones concurrentes
- **100+ usuarios simultáneos**: Go maneja muy bien concurrencia nativa
- **Conexión a BD**: Pool de conexiones MySQL (configurable en GORM)
- **Performance**: Compilado a binario nativo, muy rápido
- **Caché**: Considerar Redis para estadísticas frecuentes (opcional)
- **Deployment**: Binario único, fácil de desplegar (Docker, systemd, etc.)

## Dependencias Go Principales

```go
require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.5
    gorm.io/driver/mysql v1.5.2
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/go-playground/validator/v10 v10.16.0
    github.com/joho/godotenv v1.5.1
    golang.org/x/crypto v0.17.0 // bcrypt
)
```

## Comandos Útiles

```bash
# Inicializar módulo Go
go mod init github.com/usuario/sgp

# Descargar dependencias
go mod tidy

# Ejecutar servidor (desarrollo)
go run cmd/server/main.go

# Compilar binario
go build -o bin/server cmd/server/main.go

# Ejecutar tests
go test ./...

# Ejecutar con hot-reload (air)
air
```
