# Sistema de Gestión de Granja Porcina

Sistema web para la gestión integral del ciclo de maternidad en granjas de cría porcina.

## Características Principales

- ✅ Gestión de cerdas y padrillos
- ✅ Registro de servicios (monta natural e inseminación)
- ✅ Control de preñez (confirmación y cancelación)
- ✅ Registro de partos y destetes
- ✅ Calendario de partos y destetes futuros
- ✅ Estadísticas y reportes
- ✅ Exportación a Excel
- ✅ Sistema de autenticación y roles

## Arquitectura

Ver [ARQUITECTURA.md](./ARQUITECTURA.md) para detalles completos.

### Stack Tecnológico
- **Backend**: Node.js + Express + MySQL
- **Frontend**: HTML/CSS/JavaScript (Vanilla)
- **Base de Datos**: MySQL 8.0+
- **ORM**: Sequelize

## Instalación

### Prerrequisitos
- Node.js v18+
- MySQL 8.0+
- npm o yarn

### Pasos

1. **Instalar dependencias**
```bash
npm install
```

2. **Configurar base de datos**
   - Crear base de datos MySQL
   - Ejecutar el script `database/schema.sql`
   - Configurar variables de entorno (ver `.env.example`)

3. **Configurar variables de entorno**
```bash
cp .env.example .env
# Editar .env con tus credenciales
```

4. **Ejecutar migraciones** (si usas Sequelize CLI)
```bash
npm run migrate
```

5. **Iniciar servidor**
```bash
npm start
# o para desarrollo con auto-reload:
npm run dev
```

## Estructura del Proyecto

```
├── backend/          # API y lógica del servidor
├── database/         # Scripts SQL
├── css/             # Estilos
├── js/              # JavaScript frontend
├── index.html       # Página principal
└── package.json     # Dependencias
```

## Ciclo de Maternidad

El sistema gestiona el siguiente ciclo:

1. **Disponible** → Cerda lista para ser servida
2. **Servicio** → Se registra un servicio
3. **Gestación** → Se confirma la preñez
4. **Cría** → Se registra el parto
5. **Disponible** → Se registra el destete (vuelve al inicio)

## API Endpoints

(Se documentará cuando esté implementado)

## Licencia

ISC
