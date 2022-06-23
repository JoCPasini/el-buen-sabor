package app

import (
	"net/http"

	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor"
	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/controllers"
	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/database"
	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/services"
	"github.com/gin-gonic/gin"
)

type App struct {
	db     database.DB
	Config elbuensabor.AppConfig

	InstrumentoService    services.IInstrumentoService
	InstrumentoController controllers.IInstrumentoController

	LoginService    services.ILoginService
	LoginController controllers.ILoginController
}

func NewApp() (*App, error) {
	// scope := ("prod")
	scope := ("dev")

	config, err := NewConfig(scope)
	if err != nil {
		return &App{}, err
	}

	mysqlDB, err := database.NewMySQL(config.DB)
	if err != nil {
		return &App{}, err
	}

	container := NewContainer(config, mysqlDB)

	app := App{
		Config:                config,
		InstrumentoService:    container.InstrumentoService,
		InstrumentoController: controllers.NewInstrumentoController(container.InstrumentoService),

		LoginService:    container.LoginService,
		LoginController: controllers.NewLoginController(container.LoginService),
	}
	return &app, nil
}

func (app *App) RegisterRoutes(router *gin.Engine) {

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
	// 	AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	// }))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	instrumentoGroup := router.Group("/instrumento")
	{
		instrumentoGroup.GET("/:idInstrumento", app.InstrumentoController.GetByID)
		instrumentoGroup.POST("", app.InstrumentoController.AddInstrument)
		instrumentoGroup.GET("/getAll", app.InstrumentoController.GetAll)
		instrumentoGroup.DELETE("/:idInstrumento", app.InstrumentoController.DeleteInstrument)
		instrumentoGroup.PUT("", app.InstrumentoController.UpdateInstrument)
	}

	login := router.Group("/login")
	{
		login.GET("", app.LoginController.LoginUsuario)
		login.POST("/register", app.LoginController.AddUsuario)
	}
}

func (a *App) CerrarDB() {
	a.db.Close()
}
