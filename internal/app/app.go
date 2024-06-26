package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor"
	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/controllers"
	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/database"
	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	db     database.DB
	Config elbuensabor.AppConfig

	LoginService    services.ILoginService
	LoginController controllers.ILoginController

	DomicilioService    services.IDomicilioService
	DomicilioController controllers.IDomicilioController

	PedidoService    services.IPedidoService
	PedidoController controllers.IPedidoController

	FacturaService    services.IFacturaService
	FacturaController controllers.IFacturaController

	ArticuloManufacturadoDetalleService    services.IArticuloManufacturadoDetalleService
	ArticuloManufacturadoDetalleController controllers.IArticuloManufacturadoDetalleController

	ArticuloManufacturadoService    services.IArticuloManufacturadoService
	ArticuloManufacturadoController controllers.IArticuloManufacturadoController

	ArticuloInsumoService    services.IArticuloInsumoService
	ArticuloInsumoController controllers.IArticuloInsumoController

	CategoriaService    services.ICategoriaService
	CategoriaController controllers.ICategoriaController

	MercadoPagoController controllers.IMercadoPagoController
}

func NewApp() (*App, error) {
	scope := ("prod")
	//scope := ("dev")

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
		db:     mysqlDB,
		Config: config,

		LoginService:    container.LoginService,
		LoginController: controllers.NewLoginController(container.LoginService),

		PedidoService:    container.PedidoService,
		PedidoController: controllers.NewPedidoController(container.PedidoService),

		FacturaService:    container.FacturaService,
		FacturaController: controllers.NewFacturaController(container.FacturaService),

		ArticuloManufacturadoDetalleService:    container.ArticuloManufacturadoDetalleService,
		ArticuloManufacturadoDetalleController: controllers.NewArticuloManufacturadoDetalleController(container.ArticuloManufacturadoDetalleService),

		ArticuloManufacturadoService:    container.ArticuloManufacturadoService,
		ArticuloManufacturadoController: controllers.NewArticuloManufacturadoController(container.ArticuloManufacturadoService),

		ArticuloInsumoService:    container.ArticuloInsumoService,
		ArticuloInsumoController: controllers.NewArticuloInsumoController(container.ArticuloInsumoService),

		DomicilioService:    container.DomicilioService,
		DomicilioController: controllers.NewDomicilioController(container.DomicilioService),

		CategoriaService:    container.CategoriaService,
		CategoriaController: controllers.NewCategoriaController(container.CategoriaService),

		MercadoPagoController: controllers.NewMercadoPagoController(),
	}
	return &app, nil
}

// func controlHorarios(ctx *gin.Context) {
// 	now := time.Now()
// 	if now.Hour() < 20 || now.Hour() >= 24 {
// 		ctx.AbortWithStatusJSON(405, "comercio cerrado de 20 a 24 hs (horario Argentina). lo esperamos en ese horario. saludos")
// 		return
// 	} else {
// 		log.Println("horario aceptado")
// 		ctx.Next()
// 	}
// }

func uploader(ctx *gin.Context) {
	r := ctx.Request
	r.ParseMultipartForm(2000)
	file, fileInfo, err := r.FormFile("archivo")

	f, err := os.OpenFile("./files/"+fileInfo.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		ctx.JSON(500, "error cargando el archivo")
		log.Fatal("error al cargar el archivo")
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Println("File:", file)
	fmt.Println("fileInfo:", fileInfo.Filename)
	ctx.JSON(200, fileInfo.Filename)
}

func (app *App) RegisterRoutes(router *gin.Engine) {

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE", "GET"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers, access-control-allow-credentials"},
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/files", uploader)

	register := router.Group("/register")
	{
		register.POST("", app.LoginController.AddUsuario)
	}

	login := router.Group("/login")
	{
		login.POST("", app.LoginController.LoginUsuario)
	}

	// mercado pago API - test
	mercadopago := router.Group("/mercado-pago")
	mercadopago.POST("/pagar", app.MercadoPagoController.Pagar)
	mercadopago.GET("/metodos-de-pago", app.MercadoPagoController.MetodosDePago)

	pedido := router.Group("")
	pedido.PUT("/generar-pedido", app.PedidoController.GenerarPedido)
	pedido.PUT("/aceptar-pedido/:idPedido", app.PedidoController.AceptarPedido)
	pedido.GET("/detalle-pedido/:idPedido", app.PedidoController.GetAllDetallePedidosByIDPedido)
	pedido.PUT("/pedido/update-estado", app.PedidoController.UpdateEstadoPedido)
	pedido.GET("/verificar-stock/:idArticulo/:amount/:esBebida", app.PedidoController.VerificarStock)
	router.GET("/carrito-completo-getAll", app.ArticuloInsumoController.GetAllCarritoCompleto)

	// Rankings :: Excels
	// Agregar esta validación al repo cuando esté :: (AND p.estado = 'confirmado y pagado')
	router.GET("/ranking-comidas", app.PedidoController.RankingComidasMasPedidas)
	router.GET("/pedidos-por-cliente", app.PedidoController.GetRankingDePedidosPorCliente)
	router.GET("/recaudaciones-diarias", app.FacturaController.RecaudacionesDiarias)
	router.GET("/recaudaciones-mensuales", app.FacturaController.RecaudacionesMensuales)
	router.GET("/recaudaciones-periodo-tiempo", app.FacturaController.RecaudacionesPeriodoTiempo)
	router.GET("/ganancias", app.FacturaController.ObtenerGanancias)

	usuarios := router.Group("/usuarios")
	{
		usuarios.GET("", app.LoginController.GetAllUsuarios)
		usuarios.GET("/:id", app.LoginController.GetUsuarioByID)
		usuarios.DELETE("/:id", app.LoginController.DeleteUsuarioByID)
		usuarios.PUT("", app.LoginController.UpdateUsuario)
	}

	domicilio := router.Group("/domicilio")
	{
		domicilio.GET("/:idUsuario", app.DomicilioController.GetAllDomicilioByUsuario)
		domicilio.POST("", app.DomicilioController.AddDomicilio)
		domicilio.PUT("", app.DomicilioController.UpdateDomicilio)
		// domicilio.DELETE("/:id", app.LoginController.DeleteUsuarioByID)
	}

	categoria := router.Group("/categoria")
	{
		categoria.GET("/getAll", app.CategoriaController.GetAllCategoria)
		categoria.POST("", app.CategoriaController.AddCategoria)
	}

	instrumentoGroup := router.Group("/factura")
	{
		instrumentoGroup.GET("/:idFactura", app.FacturaController.GetByID)
		instrumentoGroup.POST("", app.FacturaController.AddFactura)
		instrumentoGroup.GET("/getAll", app.FacturaController.GetAll)
		instrumentoGroup.DELETE("/:idFactura", app.FacturaController.DeleteFactura)
		instrumentoGroup.PUT("", app.FacturaController.UpdateFactura)
		instrumentoGroup.GET("/getByPedido/:idPedido", app.FacturaController.GetByIDPedido)
		instrumentoGroup.GET("/getAllByCliente/:idCliente", app.FacturaController.GetAllByCliente)
	}

	productoGroup := router.Group("/pedido")
	{
		productoGroup.GET("/:idPedido", app.PedidoController.GetByID)
		productoGroup.GET("/byCliente/:idCliente", app.PedidoController.GetAllPedidosByIDCliente)
		//productoGroup.POST("", app.PedidoController.AddPedido)
		productoGroup.GET("/getAll", app.PedidoController.GetAll)
		productoGroup.DELETE("/:idPedido", app.PedidoController.DeletePedido)
		productoGroup.PUT("", app.PedidoController.UpdatePedido)
	}

	articuloInsumo := router.Group("/articulo-insumo")
	{
		articuloInsumo.GET("/:id", app.ArticuloInsumoController.GetByID)
		articuloInsumo.POST("", app.ArticuloInsumoController.AddArticuloInsumo)
		articuloInsumo.GET("/getAll", app.ArticuloInsumoController.GetAll)
		articuloInsumo.DELETE("/:id", app.ArticuloInsumoController.DeleteArticuloInsumo)
		articuloInsumo.PUT("", app.ArticuloInsumoController.UpdateArticuloInsumo)
		articuloInsumo.PUT("/agregar-stock-insumo", app.ArticuloInsumoController.AgregarStockInsumo)
	}

	articuloManufacturadoDetalle := router.Group("/articulo-manufacturado-detalle")
	{
		articuloManufacturadoDetalle.GET("/:id", app.ArticuloManufacturadoDetalleController.GetByID)
		articuloManufacturadoDetalle.POST("", uploader, app.ArticuloManufacturadoDetalleController.AddArticuloManufacturadoDetalle)
		articuloManufacturadoDetalle.GET("/getAll", uploader, app.ArticuloManufacturadoDetalleController.GetAll)
		articuloManufacturadoDetalle.DELETE("/:id", app.ArticuloManufacturadoDetalleController.DeleteArticuloManufacturadoDetalle)
		articuloManufacturadoDetalle.PUT("", app.ArticuloManufacturadoDetalleController.UpdateArticuloManufacturadoDetalle)
	}

	articuloManufacturado := router.Group("/articulo-manufacturado")
	{
		articuloManufacturado.GET("/:id", app.ArticuloManufacturadoController.GetByID)
		articuloManufacturado.POST("", app.ArticuloManufacturadoController.AddArticuloManufacturado)
		articuloManufacturado.GET("/getAll", app.ArticuloManufacturadoController.GetAll)
		articuloManufacturado.DELETE("/:id", app.ArticuloManufacturadoController.DeleteArticuloManufacturado)
		articuloManufacturado.PUT("", app.ArticuloManufacturadoController.UpdateArticuloManufacturado)
		articuloManufacturado.GET("/getAllAvailable", app.ArticuloManufacturadoController.GetAllAvailable)
	}

}

func (a *App) CerrarDB() {
	a.db.Close()
}
