package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/database"
	"github.com/JosePasiniMercadolibre/el-buen-sabor/internal/elbuensabor/domain"
	"github.com/jmoiron/sqlx"
)

type articuloInsumoDB struct {
	ID           int             `db:"id"`
	Denominacion sql.NullString  `db:"denominacion"`
	PrecioCompra sql.NullFloat64 `db:"precio_compra"`
	PrecioVenta  sql.NullFloat64 `db:"precio_venta"`
	StockActual  sql.NullInt32   `db:"stock_actual"`
	StockMinimo  sql.NullInt32   `db:"stock_minimo"`
	UnidadMedida sql.NullString  `db:"unidad_medida"`
	EsInsumo     sql.NullBool    `db:"es_insumo"`
	Imagen       sql.NullString  `db:"imagen"`
}

func (a *articuloInsumoDB) toArticuloInsumo() domain.ArticuloInsumo {
	return domain.ArticuloInsumo{
		ID:           a.ID,
		Denominacion: database.ToStringP(a.Denominacion),
		PrecioCompra: database.ToFloat64P(a.PrecioCompra),
		PrecioVenta:  database.ToFloat64P(a.PrecioVenta),
		StockActual:  database.ToIntP(a.StockActual),
		StockMinimo:  database.ToIntP(a.StockMinimo),
		UnidadMedida: database.ToStringP(a.UnidadMedida),
		EsInsumo:     database.ToBoolP(a.EsInsumo),
		Imagen:       database.ToStringP(a.Imagen),
	}
}

type carritoCompletoDB struct {
	ID                   int             `db:"id"`
	Denominacion         sql.NullString  `db:"denominacion"`
	PrecioCompra         sql.NullFloat64 `db:"precio_compra"`
	PrecioVenta          sql.NullFloat64 `db:"precio_venta"`
	Cantidad             sql.NullInt32   `db:"cantidad"`
	StockActual          sql.NullInt32   `db:"stock_actual"`
	StockMinimo          sql.NullInt32   `db:"stock_minimo"`
	Imagen               sql.NullString  `db:"imagen"`
	EsBebida             bool            `json:"es_bebida"`
	TiempoEstimadoCocina sql.NullInt32   `db:"tiempo_estimado_cocina"`
	IDCategoria          sql.NullInt32   `db:"id_categoria"`
}

func (a *carritoCompletoDB) toCarritoCompleto() domain.CarritoCompleto {
	return domain.CarritoCompleto{
		ID:                   a.ID,
		Denominacion:         database.ToStringP(a.Denominacion),
		PrecioCompra:         database.ToFloat64P(a.PrecioCompra),
		PrecioVenta:          database.ToFloat64P(a.PrecioVenta),
		Cantidad:             0,
		StockActual:          database.ToIntP(a.StockActual),
		StockMinimo:          database.ToIntP(a.StockMinimo),
		Imagen:               database.ToStringP(a.Imagen),
		EsBebida:             a.EsBebida,
		TiempoEstimadoCocina: database.ToIntP(a.TiempoEstimadoCocina),
		IDCategoria:          database.ToIntP(a.IDCategoria),
	}
}

type IArticuloInsumoRepository interface {
	Insert(ctx context.Context, tx *sqlx.Tx, articulo_manufacturado_detalle domain.ArticuloInsumo) error
	GetByID(ctx context.Context, tx *sqlx.Tx, id int) (*domain.ArticuloInsumo, error)
	GetAll(ctx context.Context, tx *sqlx.Tx) ([]domain.ArticuloInsumo, error)
	GetAllCarritoCompleto(ctx context.Context, tx *sqlx.Tx) ([]domain.CarritoCompleto, error)
	Update(ctx context.Context, tx *sqlx.Tx, articulo_manufacturado_detalle domain.ArticuloInsumo) error
	SumarStockInsumo(ctx context.Context, tx *sqlx.Tx, agregarStock domain.AgregarStockInsumo) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int) error
}

type MySQLArticuloInsumoRepository struct {
	qInsert     string
	qGetByID    string
	qGetAll     string
	qDeleteById string
	qUpdate     string
}

//art.Denominacion, art.PrecioCompra, art.PrecioVenta, art.StockActual, art.StockMinimo, art.UnidadMedida, art.UnidadMedida, art.EsInsumo, art.ID
func NewMySQLArticuloInsumoRepository() *MySQLArticuloInsumoRepository {
	return &MySQLArticuloInsumoRepository{
		qInsert:     "INSERT INTO articulo_insumo (denominacion, precio_compra, precio_venta, stock_actual, stock_minimo, unidad_medida, es_insumo, imagen) VALUES (?,?,?,?,?,?,?,(replace(?, ' ', '')))",
		qGetByID:    "SELECT id, denominacion, precio_compra, precio_venta, stock_actual, stock_minimo, unidad_medida, es_insumo FROM articulo_insumo WHERE id = ?",
		qGetAll:     "SELECT id, denominacion, precio_compra, precio_venta, stock_actual, stock_minimo, unidad_medida, es_insumo FROM articulo_insumo",
		qDeleteById: "DELETE FROM articulo_insumo WHERE id = ?",
		qUpdate:     "UPDATE articulo_insumo SET denominacion = COALESCE(?,denominacion), precio_compra = COALESCE(?,precio_compra), precio_venta = COALESCE(?,precio_venta), stock_actual = COALESCE(?,stock_actual), stock_minimo = COALESCE(?,stock_minimo), unidad_medida = COALESCE(?,unidad_medida), es_insumo = COALESCE(?,es_insumo) WHERE id = ?",
	}
}

func (i *MySQLArticuloInsumoRepository) Insert(ctx context.Context, tx *sqlx.Tx, art domain.ArticuloInsumo) error {
	query := i.qInsert
	_, err := tx.ExecContext(ctx, query, art.Denominacion, art.PrecioCompra, art.PrecioVenta, art.StockActual, art.StockMinimo, art.UnidadMedida, art.EsInsumo, art.Imagen)
	return err
}

func (i *MySQLArticuloInsumoRepository) GetByID(ctx context.Context, tx *sqlx.Tx, id int) (*domain.ArticuloInsumo, error) {
	query := i.qGetByID
	var articuloInsumo articuloInsumoDB

	row := tx.QueryRowxContext(ctx, query, id)
	err := row.StructScan(&articuloInsumo)
	if err != nil {
		return nil, err
	}
	artIns := articuloInsumo.toArticuloInsumo()
	return &artIns, nil
}

func (i *MySQLArticuloInsumoRepository) GetAll(ctx context.Context, tx *sqlx.Tx) ([]domain.ArticuloInsumo, error) {
	query := i.qGetAll
	articulos := make([]domain.ArticuloInsumo, 0)

	rows, err := tx.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var articulo articuloInsumoDB
		if err := rows.StructScan(&articulo); err != nil {
			return articulos, err
		}
		articulos = append(articulos, articulo.toArticuloInsumo())
	}
	return articulos, nil
}

func (i *MySQLArticuloInsumoRepository) GetAllCarritoCompleto(ctx context.Context, tx *sqlx.Tx) ([]domain.CarritoCompleto, error) {
	queryBebidas := `SELECT id, denominacion, precio_compra, precio_venta, stock_actual, stock_minimo, imagen, id_categoria
		FROM articulo_insumo WHERE es_insumo = false AND stock_actual > 0;`
	carritoCompleto := make([]domain.CarritoCompleto, 0)
	rows, err := tx.QueryxContext(ctx, queryBebidas)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println(" 111111 :::::::::::::::::")

	var ok bool = true
	for rows.Next() {
		var carritoStruct carritoCompletoDB
		if err := rows.StructScan(&carritoStruct); err != nil {
			return carritoCompleto, err
		}
		carritoStruct.EsBebida = ok
		carritoCompleto = append(carritoCompleto, carritoStruct.toCarritoCompleto())
	}
	fmt.Println(":::::::::::::::::")
	fmt.Println("Carrito:", carritoCompleto)
	fmt.Println(":::::::::::::::::")

	// ------------------- Hasta acá la query para traer las bebidas ------------------- //
	/*
		queryPlatos := `select id, tiempo_estimado_cocina, denominacion, precio_venta, imagen from articulo_manufacturado;`
	*/
	queryPlatos := `SELECT id, tiempo_estimado_cocina, denominacion, precio_venta, imagen, id_categoria FROM articulo_manufacturado;`
	rows, err = tx.QueryxContext(ctx, queryPlatos)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var carritoStruct carritoCompletoDB
		if err := rows.StructScan(&carritoStruct); err != nil {
			return carritoCompleto, err
		}
		carritoCompleto = append(carritoCompleto, carritoStruct.toCarritoCompleto())
	}
	fmt.Println(":::::::::::::::::")
	fmt.Println("Carrito:", carritoCompleto)
	fmt.Println(":::::::::::::::::")
	return carritoCompleto, nil
}

func (i *MySQLArticuloInsumoRepository) Update(ctx context.Context, tx *sqlx.Tx, art domain.ArticuloInsumo) error {
	query := i.qUpdate
	_, err := tx.ExecContext(ctx, query, art.Denominacion, art.PrecioCompra, art.PrecioVenta, art.StockActual, art.StockMinimo, art.UnidadMedida, art.EsInsumo, art.ID)
	return err
}

func (i *MySQLArticuloInsumoRepository) SumarStockInsumo(ctx context.Context, tx *sqlx.Tx, art domain.AgregarStockInsumo) error {
	query := "UPDATE articulo_insumo SET stock_actual = (stock_actual + ? ) WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, art.Cantidad, art.IDArticuloInsumo)
	return err
}

func (i *MySQLArticuloInsumoRepository) Delete(ctx context.Context, tx *sqlx.Tx, id int) error {
	query := i.qDeleteById
	_, err := tx.ExecContext(ctx, query, id)
	return err
}
