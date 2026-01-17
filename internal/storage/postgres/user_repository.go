package postgres

import (
	"context"
	"errors"
	"fmt"
	model "tic-tac-toe/internal/domain/model/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) model.UserRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (uuid, login, password)
		VALUES ($1,$2,$3)`
	_, err := r.pool.Exec(ctx, query, user.UUID, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("Ошибка создания пользователя: %w", err)
	}
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	query := `SELECT uuid, login
		FROM users
		WHERE uuid = $1`

	var userID uuid.UUID
	var userLogin string
	err := r.pool.QueryRow(ctx, query, id).Scan(&userID, &userLogin)
	if err != nil {
		return model.User{}, ErrUserNotFound
	}
	return model.User{
		UUID:  userID,
		Login: userLogin,
	}, nil
}

func (r *userRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `SELECT uuid, login, password 
		FROM users
		WHERE login = $1`

	var userID uuid.UUID
	var userLogin, userPassword string
	err := r.pool.QueryRow(ctx, query, login).Scan(&userID, &userLogin, &userPassword)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &model.User{
		UUID:     userID,
		Login:    userLogin,
		Password: userPassword,
	}, nil
}
