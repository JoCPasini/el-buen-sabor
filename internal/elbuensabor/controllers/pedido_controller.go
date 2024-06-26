package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/domain"
	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/services"
	"github.com/gin-gonic/gin"
)

var (
	StockInsuficiente = errors.New("stock insuficiente")
)

type IPedidoController interface {
	GetByID(*gin.Context)
	GetAllPedidosByIDCliente(*gin.Context)
	GetAll(*gin.Context)
	//AddPedido(*gin.Context)
	AceptarPedido(*gin.Context)
	GenerarPedido(*gin.Context)
	UpdatePedido(*gin.Context)
	UpdateEstadoPedido(*gin.Context)
	DeletePedido(*gin.Context)
	RankingComidasMasPedidas(*gin.Context)
	GetAllDetallePedidosByIDPedido(*gin.Context)
	GetRankingDePedidosPorCliente(*gin.Context)
	VerificarStock(*gin.Context)
}

type PedidoController struct {
	service services.IPedidoService
}

func NewPedidoController(service services.IPedidoService) *PedidoController {
	return &PedidoController{service}
}

func (c PedidoController) GetAllPedidosByIDCliente(ctx *gin.Context) {
	clienteID := ctx.Param("idCliente")

	ID, err := strconv.Atoi(clienteID)
	if err != nil {
		ctx.JSON(400, errors.New("invalid pedido id"))
		return
	}
	pedidos, err := c.service.GetAllPedidosByIDCliente(ctx, ID)
	if err != nil {
		if err.Error() == errInternal.Error() {
			ctx.JSON(404, gin.H{
				"message": "pedido not found",
			})
			return
		}
	}
	if err != nil {
		ctx.JSON(500, errors.New("Error internal server error: "+err.Error()))
		return
	}
	ctx.JSON(200, pedidos)
}

func (c PedidoController) GetByID(ctx *gin.Context) {
	pedidoID := ctx.Param("idPedido")

	ID, err := strconv.Atoi(pedidoID)
	if err != nil {
		ctx.JSON(400, errors.New("invalid pedido id"))
		return
	}
	pedido, err := c.service.GetByID(ctx, ID)
	if err != nil {
		if err.Error() == errInternal.Error() {
			ctx.JSON(404, gin.H{
				"message": "pedido not found",
			})
			return
		}
	}
	if err != nil {
		ctx.JSON(500, errors.New("Error internal server error: "+err.Error()))
		return
	}
	ctx.JSON(200, pedido)
}

func (c PedidoController) GetAll(ctx *gin.Context) {
	pedidos, err := c.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(500, errors.New("Error internal server error: "+err.Error()))
		return
	}
	ctx.JSON(200, pedidos)
}

func (c PedidoController) GetAllDetallePedidosByIDPedido(ctx *gin.Context) {
	idParam := ctx.Param("idPedido")
	fmt.Println("idPedido", idParam)
	idPedido, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, errors.New("invalid pedido id"))
		return
	}
	fmt.Println("idPedido", idPedido)

	pedidos, err := c.service.GetAllDetallePedidosByIDPedido(ctx, idPedido)
	if err != nil {
		ctx.JSON(500, errors.New("Error internal server error: "+err.Error()))
		return
	}
	ctx.JSON(200, pedidos)
}

func (c PedidoController) GetRankingDePedidosPorCliente(ctx *gin.Context) {
	desde := ctx.Query("desde")
	hasta := ctx.Query("hasta")
	fmt.Println("desde", desde)
	fmt.Println("hasta", hasta)

	pedidosByIdCLiente, err := c.service.GetRankingDePedidosPorCliente(ctx, desde, hasta)
	if err != nil {
		ctx.JSON(500, errors.New("Error internal server error: "+err.Error()))
		return
	}
	ctx.JSON(200, pedidosByIdCLiente)
}

func (c PedidoController) GenerarPedido(ctx *gin.Context) {
	var pedido domain.GenerarPedido
	err := ctx.BindJSON(&pedido)
	if err != nil {
		ctx.JSON(400, errors.New("generate pedido error"))
		return
	}
	fmt.Println(pedido)
	fmt.Println(":::", pedido.Pedido)
	pedido, err = c.service.GenerarPedido(ctx, pedido)
	if err != nil {
		ctx.JSON(400, errors.New("generate pedido error"))
		return
	}
	fmt.Println("3")

	ctx.JSON(200, gin.H{"status": 200, "pedido": pedido})
}

func (c PedidoController) AceptarPedido(ctx *gin.Context) {
	fmt.Println(" -- principio --")
	idParam := ctx.Param("idPedido")

	fmt.Println(idParam)
	idPedido, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, errors.New("add pedido error"))
		return
	}
	ok, err := c.service.AceptarPedido(ctx, idPedido)
	if err != nil || !ok {
		fmt.Println("ERRORRR:", err)
		if err.Error() == StockInsuficiente.Error() {
			_ = c.service.CancelarPedido(ctx, idPedido)
		}
		ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
		return
	}
	fmt.Println(" -- fin -- ")
	ctx.JSON(200, gin.H{"status": 200, "ok": ok})
}

func (c PedidoController) VerificarStock(ctx *gin.Context) {
	idParam := ctx.Param("idArticulo")
	amountParam := ctx.Param("amount")
	esBebidaParam := ctx.Param("esBebida")

	idPedido, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, errors.New("add pedido error"))
		return
	}

	amount, err := strconv.Atoi(amountParam)
	if err != nil {
		ctx.JSON(400, errors.New("add pedido error"))
		return
	}

	esBebida, err := strconv.ParseBool(esBebidaParam)
	if err != nil {
		ctx.JSON(400, errors.New("add pedido error"))
		return
	}

	fmt.Println(idPedido)
	fmt.Println(amount)
	fmt.Println(esBebida)

	ok, err := c.service.VerificarStock(ctx, idPedido, amount, esBebida)
	if err != nil || !ok {
		ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
		return
	}
	fmt.Println(" -- fin -- ")
	ctx.JSON(200, gin.H{"status": 200, "ok": ok})
}

// func (c PedidoController) AddPedido(ctx *gin.Context) {
// 	var pedido domain.Pedido
// 	err := ctx.BindJSON(&pedido)
// 	if err != nil {
// 		ctx.JSON(400, errors.New("error bind"))
// 		return
// 	}

// 	err = c.service.AddPedido(ctx, pedido)
// 	if err != nil {
// 		ctx.JSON(400, errors.New("add pedido error"))
// 		return
// 	}
// 	ctx.JSON(200, gin.H{"status": 200, "pedido": pedido})
// }

func (c PedidoController) UpdatePedido(ctx *gin.Context) {
	var pedido domain.Pedido

	err := ctx.BindJSON(&pedido)

	if err != nil {
		ctx.JSON(400, errors.New("invalid pedido to be updated"))
		return
	}

	err = c.service.UpdatePedido(ctx, pedido)
	if err != nil {
		ctx.JSON(500, errors.New("internal error server"))
		return
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "pedido updated successfully",
	})
}

func (c PedidoController) UpdateEstadoPedido(ctx *gin.Context) {
	var pedido domain.PedidoEstado

	err := ctx.BindJSON(&pedido)

	if err != nil {
		ctx.JSON(400, errors.New("invalid pedido to be updated"))
		return
	}

	err = c.service.UpdateEstadoPedido(ctx, pedido.Estado, pedido.IDPedido)
	if err != nil {
		ctx.JSON(500, errors.New("internal error server"))
		return
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "pedido updated successfully",
	})
}

func (c PedidoController) DeletePedido(ctx *gin.Context) {
	idPedido := ctx.Param("idPedido")
	if idPedido == "" {
		ctx.JSON(400, errors.New("invalid pedido"))
		return
	}

	ID, err := strconv.Atoi(idPedido)
	if err != nil {
		ctx.JSON(400, errors.New("id pedido must be a number"))
		return
	}

	err = c.service.DeletePedido(ctx, ID)
	if err != nil {
		if err.Error() == errNotFound.Error() {
			ctx.JSON(404, gin.H{
				"message": "pedido not found",
			})
			return
		}
	}

	if err != nil {
		ctx.JSON(500, errors.New("Error internal server error: "+err.Error()))
		return
	}

	ctx.JSON(200, gin.H{
		"message": "pedido deleted successfully",
	})
}

func (c PedidoController) RankingComidasMasPedidas(ctx *gin.Context) {
	var err error
	var comidasMasPedidas []domain.RankingComidasMasPedidas
	desde := ctx.Query("desde")
	hasta := ctx.Query("hasta")
	fmt.Println("desde", desde)
	fmt.Println("hasta", hasta)
	comidasMasPedidas, err = c.service.RankingComidasMasPedidas(ctx, desde, hasta)
	if err != nil {
		ctx.JSON(400, errors.New("generate pedido error"))
		return
	}
	ctx.JSON(200, comidasMasPedidas)
}
