package persistence

import (
	"context"
	"database/sql"

	"auth-service/internal/infrastructure/auth"

	"github.com/jmoiron/sqlx"
)

type PostgresSessionRepository struct {
	db *sqlx.DB
}

func NewPostgresSessionRepository(db *sqlx.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) Create(ctx context.Context, session *auth.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.IPAddress,
		session.ExpiresAt,
		session.CreatedAt,
	)

	return err
}

func (r *PostgresSessionRepository) GetByID(ctx context.Context, sessionID string) (*auth.Session, error) {
	var session auth.Session
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at
		FROM sessions WHERE id = $1
	`

	err := r.db.GetContext(ctx, &session, query, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (r *PostgresSessionRepository) Delete(ctx context.Context, sessionID string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, sessionID)
	return err
}

func (r *PostgresSessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *PostgresSessionRepository) ListByUserID(ctx context.Context, userID string) ([]*auth.Session, error) {
	var sessions []*auth.Session
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at
		FROM sessions WHERE user_id = $1
	`

	err := r.db.SelectContext(ctx, &sessions, query, userID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
