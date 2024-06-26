package domain

import (
	"context"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type ILoginRepository interface {
	Insert(ctx context.Context, tx *sqlx.Tx, usuario Usuario) error
	GetHash(ctx context.Context, tx *sqlx.Tx, email string) (string, error)
	GetAllUsuarios(ctx context.Context, tx *sqlx.Tx) ([]UsuarioResponse, error)
	CantidadDeCocineros(ctx context.Context, tx *sqlx.Tx) (int, error)
	GetUsuarioByID(ctx context.Context, tx *sqlx.Tx, id int) (UsuarioResponse, error)
	GetUsuarioByEmail(ctx context.Context, tx *sqlx.Tx, email string) (UsuarioResponse, error)
	DeleteUsuarioByID(ctx context.Context, tx *sqlx.Tx, id int) (bool, error)
	UpdateUsuario(ctx context.Context, tx *sqlx.Tx, usuario Usuario) (Usuario, error)
}

const (
	USUARIO       = 100
	CAJERO        = 200
	COCINERO      = 300
	DELIVERY      = 400
	ADMINISTRADOR = 500
)

type Usuario struct {
	ID       int     `json:"id"`
	Nombre   *string `json:"nombre"`
	Apellido *string `json:"apellido"`
	Usuario  *string `json:"usuario"`
	Email    *string `json:"email"`
	Hash     *string `json:"hash"`
	Rol      int     `json:"rol"`
}

func (u Usuario) isValid() bool {
	return u.ID > 0
}

func (u Usuario) GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, err
	} else {
		return true, nil
	}
}
