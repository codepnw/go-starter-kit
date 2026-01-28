package userrepository

import (
	"context"
	"database/sql"

	"github.com/codepnw/go-starter-kit/internal/features/user"
)

//go:generate mockgen -source=user_repo.go -destination=user_repo_mock.go -package=userrepository
type UserRepository interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	FindUserByEmail(ctx context.Context, email string) (*user.User, error)

	// Transaction
	InsertUserTx(ctx context.Context, tx *sql.Tx, u *user.User) error
	InsertRefreshTokenTx(ctx context.Context, tx *sql.Tx, token *user.RefreshToken) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) InsertUserTx(ctx context.Context, tx *sql.Tx, u *user.User) error {
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2) RETURNING id, created_at, updated_at
	`
	if err := tx.QueryRowContext(ctx, query, u.Email, u.Password).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var dummy bool
	query := `SELECT 1 FROM users WHERE email = $1 LIMIT 1`

	if err := r.db.QueryRowContext(ctx, query, email).Scan(&dummy); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	query := `
		SELECT id, email, password
		FROM users WHERE email = $1 LIMIT 1
	`
	if err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) InsertRefreshTokenTx(ctx context.Context, tx *sql.Tx, token *user.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, revoked)
		VALUES ($1, $2, $3, $4)
	`
	_, err := tx.ExecContext(
		ctx,
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.Revoked,
	)
	if err != nil {
		return err
	}
	return nil
}
