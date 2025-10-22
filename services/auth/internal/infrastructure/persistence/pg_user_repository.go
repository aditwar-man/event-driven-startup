package persistence

import (
	"auth-service/internal/domain"
	"context"
	"database/sql"
	"errors"

	sharedDomain "shared/pkg/domain"

	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *sharedDomain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, tier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		string(user.Tier),
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*sharedDomain.User, error) {
	var user sharedDomain.User
	query := `SELECT id, email, password_hash, full_name, tier, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*sharedDomain.User, error) {
	var user sharedDomain.User
	query := `SELECT id, email, password_hash, full_name, tier, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *sharedDomain.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, full_name = $3, tier = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		string(user.Tier),
		user.UpdatedAt,
		user.ID,
	)

	return err
}
