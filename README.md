# Sistema de Gestión de Granja Porcina (SGP)

Aplicación web para gestionar granjas porcinas, enfocada en el ciclo reproductivo y el seguimiento operativo (servicios, preñez, partos, destetes, lotes, ventas y estadísticas).

## Stack actual

- **Backend**: Go + Gin + GORM
- **Frontend**: HTML/CSS/JavaScript (vanilla), servido por el backend
- **Base de datos**: MySQL 8+
- **Autenticación**: JWT

## Funcionalidades principales

- Gestión de granjas, corrales, cerdas y padrillos
- Registro de servicios, partos, destetes y muertes de lechones
- Gestión de lotes y ventas
- Calendario de eventos futuros
- Estadísticas por granja y por período
- API protegida con autenticación por token

## Requisitos

- Go 1.21+
- MySQL 8+
- Make (opcional, recomendado)

## Puesta en marcha rápida

1. **Configurar entorno backend**
   ```bash
   cd backend
   cp .env.example .env
   ```
   Edita `backend/.env` con tu configuración de base de datos y JWT.

2. **Crear la base de datos**
   ```sql
   CREATE DATABASE granja_porcina;
   ```

3. **Inicializar esquema**
   ```bash
   mysql -u <usuario> -p granja_porcina < ../database/schema.sql
   ```

4. **Instalar dependencias y ejecutar**
   ```bash
   cd backend
   make install
   make run
   ```

5. **Acceder a la app**
   - Frontend: `http://localhost:8080/`
   - Healthcheck API: `http://localhost:8080/api/health`

## Comandos útiles (backend)

```bash
cd backend
make run      # desarrollo
make dev      # hot reload (requiere air)
make test     # tests
make build    # compilar binario
```

## Estructura del proyecto

```text
backend/      # API, lógica de negocio, acceso a datos
frontend/     # interfaz web (HTML, CSS, JS)
database/     # scripts SQL y diagrama ER
ARQUITECTURA.md
DISEÑO.md
```

## Documentación adicional

- Arquitectura general: `ARQUITECTURA.md`
- Diseño funcional: `DISEÑO.md`

## Licencia

ISC
