package repository

import (
	"context"
	"project/internal/models"

	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	Conn *pgx.Conn
}

func NewUserRepo(conn *pgx.Conn) *UserRepo {
	return &UserRepo{Conn: conn}
}

func (h *UserRepo) Create(ctx context.Context, name string, email string) (models.User, error) {
	var user models.User

	sql_query := `INSERT INTO users (name,email) VALUES ($1,$2) RETURNING id,name,email`
	err := h.Conn.QueryRow(ctx, sql_query, name, email).Scan(&user.Id, &user.Name, &user.Email)

	return user, err
}

func (h *UserRepo) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	sql_query := `"SELECT id, name, email FROM users WHERE id = $1", id`
	err := h.Conn.QueryRow(
		ctx, sql_query, id).Scan(&user.Id, &user.Name, &user.Email)
	return user, err
}

func (h *UserRepo) Update(ctx context.Context, id int, name, email string) (models.User, error) {
	var user models.User
	sql_qury := `"UPDATE users SET name=$1, email=$2, updated_at=NOW() WHERE id=$3 RETURNING id, name, email"`
	err := h.Conn.QueryRow(
		ctx, sql_qury, name, email, id,
	).Scan(&user.Id, &user.Name, &user.Email)
	return user, err
}

func (h *UserRepo) Delete(ctx context.Context, id int) error {
	tag, err := h.Conn.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (h *UserRepo) Registration(name string, email string, hashpas string) (models.User, error) {
	var user models.User

	sql_quey := `INSERT INTO users (name,email,password_hash) VALUES ($1,$2,$3) RETURNING id,email,name`

	err := h.Conn.QueryRow(context.Background(), sql_quey, name, email, hashpas).Scan(&user.Id, &user.Email, &user.Name)

	return user, err
}

func (h *UserRepo) GetByEmail(email string) (models.UserWithPassword, error) {
	var passhash models.UserWithPassword
	sql_query := `SELECT id,name,email,password_hash FROM users WHERE email = $1`
	err := h.Conn.QueryRow(context.Background(), sql_query, email).Scan(&passhash.Id, &passhash.Name, &passhash.Email, &passhash.PasswordHash)
	return passhash, err
}
