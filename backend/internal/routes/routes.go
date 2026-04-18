package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/config"
	"github.com/martin/sgp/internal/handlers"
	"github.com/martin/sgp/internal/middleware"
	"github.com/martin/sgp/internal/services"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(cfg *config.Config, svc *services.ServiceContainer) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(middleware.CORS(cfg))

	// Servir frontend (archivos estaticos)
	router.StaticFile("/", "../frontend/index.html")
	router.StaticFile("/index.html", "../frontend/index.html")
	router.StaticFile("/app.html", "../frontend/app.html")
	router.Static("/css", "../frontend/css")
	router.Static("/js", "../frontend/js")
	router.Static("/assets", "../frontend/assets")
	router.Static("/img", "../frontend/img")

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(svc.Auth)
	granjaHandler := handlers.NewGranjaHandler(svc.Granja)
	corralHandler := handlers.NewCorralHandler(svc.Corral)
	loteHandler := handlers.NewLoteHandler(svc.Lote)
	cerdaHandler := handlers.NewCerdaHandler(svc.Cerda)
	padrilloHandler := handlers.NewPadrilloHandler(svc.Padrillo)
	servicioHandler := handlers.NewServicioHandler(svc.Servicio)
	partoHandler := handlers.NewPartoHandler(svc.Parto)
	desteteHandler := handlers.NewDesteteHandler(svc.Destete)
	calendarioHandler := handlers.NewCalendarioHandler(svc.Calendario)
	muerteLechonHandler := handlers.NewMuerteLechonHandler(svc.MuerteLechon)
	ventaHandler := handlers.NewVentaHandler(svc.Venta)
	estadisticasHandler := handlers.NewEstadisticasHandler(svc.Estadisticas)

	// Grupo base de la API
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "API funcionando correctamente",
			})
		})

		// =============================================
		// Rutas públicas (autenticación)
		// =============================================
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Registrar)
			auth.POST("/login", authHandler.Login)
		}

		// =============================================
		// Rutas protegidas (requieren JWT)
		// =============================================
		protected := api.Group("")
		protected.Use(middleware.Auth(cfg))
		{
			// --- Granjas ---
			// Gin requiere que el wildcard tenga el mismo nombre en un nivel,
			// por eso usamos :id para todo bajo /granjas/:id/...
			granjas := protected.Group("/granjas")
			{
				granjas.POST("", granjaHandler.Crear)
				granjas.GET("", granjaHandler.Listar)
				granjas.GET("/mis-granjas", granjaHandler.ListarPorUsuario)
				granjas.GET("/:id", granjaHandler.ObtenerPorID)
				granjas.PUT("/:id", granjaHandler.Actualizar)
				granjas.DELETE("/:id", granjaHandler.Eliminar)
				granjas.GET("/:id/estadisticas", granjaHandler.GetEstadisticas)
				granjas.POST("/:id/usuarios", granjaHandler.AsignarUsuario)
				granjas.DELETE("/:id/usuarios/:usuario_id", granjaHandler.RemoverUsuario)

				// Recursos anidados: el :id aquí es el granja_id
				granjas.POST("/:id/corrales", corralHandler.CrearEnGranja)
				granjas.GET("/:id/corrales", corralHandler.ListarPorGranja)
				granjas.POST("/:id/cerdas", cerdaHandler.CrearEnGranja)
				granjas.GET("/:id/cerdas", cerdaHandler.ListarPorGranja)
				granjas.POST("/:id/padrillos", padrilloHandler.CrearEnGranja)
				granjas.GET("/:id/padrillos", padrilloHandler.ListarPorGranja)
				granjas.GET("/:id/muertes-lechones", muerteLechonHandler.ListarPorGranja)
			}

			// --- Corrales ---
			corrales := protected.Group("/corrales")
			{
				corrales.GET("/:id", corralHandler.ObtenerPorID)
				corrales.PUT("/:id", corralHandler.Actualizar)
				corrales.DELETE("/:id", corralHandler.Eliminar)
				corrales.GET("/:id/ocupacion", corralHandler.GetOcupacion)
				corrales.POST("/:id/lotes", loteHandler.CrearEnCorral)
				corrales.GET("/:id/lotes", loteHandler.ListarPorCorral)
				corrales.GET("/:id/muertes", muerteLechonHandler.ListarPorCorral)
			}

			// --- Lotes ---
			lotes := protected.Group("/lotes")
			{
				lotes.GET("", loteHandler.ListarPorEstado)
				lotes.GET("/:id", loteHandler.ObtenerPorID)
				lotes.PUT("/:id", loteHandler.Actualizar)
				lotes.POST("/:id/cerrar", loteHandler.Cerrar)
				lotes.GET("/:id/destetes", loteHandler.GetDestetes)
			}

			// --- Cerdas ---
			cerdas := protected.Group("/cerdas")
			{
				cerdas.GET("", cerdaHandler.ListarPorEstado)
				cerdas.GET("/:id", cerdaHandler.ObtenerPorID)
				cerdas.PUT("/:id", cerdaHandler.Actualizar)
				cerdas.POST("/:id/baja", cerdaHandler.DarDeBaja)
				cerdas.GET("/:id/historial", cerdaHandler.GetHistorial)
				cerdas.GET("/:id/estadisticas", cerdaHandler.GetEstadisticas)
				cerdas.GET("/:id/servicios", servicioHandler.ListarPorCerda)
				cerdas.GET("/:id/partos", partoHandler.ListarPorCerda)
				cerdas.GET("/:id/destetes", desteteHandler.ListarPorCerda)
			}

			// --- Padrillos ---
			padrillos := protected.Group("/padrillos")
			{
				padrillos.GET("/:id", padrilloHandler.ObtenerPorID)
				padrillos.PUT("/:id", padrilloHandler.Actualizar)
				padrillos.POST("/:id/baja", padrilloHandler.DarDeBaja)
				padrillos.GET("/:id/estadisticas", padrilloHandler.GetEstadisticas)
			}

			// --- Servicios ---
			servicios := protected.Group("/servicios")
			{
				servicios.POST("", servicioHandler.Crear)
				servicios.GET("", servicioHandler.ListarPorPeriodo)
				servicios.GET("/pendientes", servicioHandler.ListarPendientesConfirmacion)
				servicios.GET("/:id", servicioHandler.ObtenerPorID)
				servicios.PUT("/:id", servicioHandler.Actualizar)
				servicios.POST("/:id/confirmar-prenez", servicioHandler.ConfirmarPrenez)
				servicios.POST("/:id/cancelar-prenez", servicioHandler.CancelarPrenez)
			}

			// --- Partos ---
			partos := protected.Group("/partos")
			{
				partos.POST("", partoHandler.Crear)
				partos.GET("", partoHandler.ListarPorPeriodo)
				partos.GET("/estadisticas", partoHandler.GetEstadisticas)
				partos.GET("/:id", partoHandler.ObtenerPorID)
				partos.PUT("/:id", partoHandler.Actualizar)
				partos.GET("/:id/muertes-lechones", muerteLechonHandler.ListarPorParto)
			}

			// --- Destetes ---
			destetes := protected.Group("/destetes")
			{
				destetes.POST("", desteteHandler.Crear)
				destetes.GET("", desteteHandler.ListarPorPeriodo)
				destetes.GET("/estadisticas", desteteHandler.GetEstadisticas)
				destetes.GET("/:id", desteteHandler.ObtenerPorID)
				destetes.PUT("/:id", desteteHandler.Actualizar)
			}

			// --- Muertes de lechones ---
			muertesLechones := protected.Group("/muertes-lechones")
			{
				muertesLechones.POST("", muerteLechonHandler.Crear)
				muertesLechones.GET("", muerteLechonHandler.ListarPorPeriodo)
				muertesLechones.GET("/estadisticas", muerteLechonHandler.GetEstadisticas)
				muertesLechones.GET("/:id", muerteLechonHandler.ObtenerPorID)
				muertesLechones.PUT("/:id", muerteLechonHandler.Actualizar)
				muertesLechones.DELETE("/:id", muerteLechonHandler.Eliminar)
			}

			// --- Ventas ---
			ventas := protected.Group("/ventas")
			{
				ventas.POST("", ventaHandler.Crear)
				ventas.GET("", ventaHandler.ListarPorPeriodo)
				ventas.GET("/estadisticas", ventaHandler.GetEstadisticas)
				ventas.GET("/:id", ventaHandler.ObtenerPorID)
				ventas.PUT("/:id", ventaHandler.Actualizar)
				ventas.DELETE("/:id", ventaHandler.Eliminar)
			}

			// --- Calendario ---
			protected.GET("/calendario", calendarioHandler.GetEventosFuturos)

			// --- Estadísticas ---
			estadisticas := protected.Group("/estadisticas")
			{
				estadisticas.GET("/granja/:id", estadisticasHandler.GetResumenGranja)
				estadisticas.GET("/periodo", estadisticasHandler.GetEstadisticasPeriodo)
			}
		}
	}

	return router
}
