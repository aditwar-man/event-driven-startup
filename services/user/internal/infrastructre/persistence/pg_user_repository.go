package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sharedDomain "shared/pkg/domain"
	"user-service/internal/domain"

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
		INSERT INTO users (
			id, email, full_name, tier, 
			ai_description_quota_used, ai_description_quota_limit,
			ai_video_quota_used, ai_video_quota_limit,
			auto_posting_quota_used, auto_posting_quota_limit,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.FullName,
		string(user.Tier),
		user.AIDescriptionQuotaUsed,
		user.AIDescriptionQuotaLimit,
		user.AIVideoQuotaUsed,
		user.AIVideoQuotaLimit,
		user.AutoPostingQuotaUsed,
		user.AutoPostingQuotaLimit,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*sharedDomain.User, error) {
	var user sharedDomain.User
	query := `
		SELECT 
			id, email, full_name, tier,
			ai_description_quota_used, ai_description_quota_limit,
			ai_video_quota_used, ai_video_quota_limit,
			auto_posting_quota_used, auto_posting_quota_limit,
			created_at, updated_at
		FROM users WHERE id = $1
	`

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
	query := `
		SELECT 
			id, email, full_name, tier,
			ai_description_quota_used, ai_description_quota_limit,
			ai_video_quota_used, ai_video_quota_limit,
			auto_posting_quota_used, auto_posting_quota_limit,
			created_at, updated_at
		FROM users WHERE email = $1
	`

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
		SET email = $1, full_name = $2, tier = $3,
			ai_description_quota_used = $4, ai_description_quota_limit = $5,
			ai_video_quota_used = $6, ai_video_quota_limit = $7,
			auto_posting_quota_used = $8, auto_posting_quota_limit = $9,
			updated_at = $10
		WHERE id = $11
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.FullName,
		string(user.Tier),
		user.AIDescriptionQuotaUsed,
		user.AIDescriptionQuotaLimit,
		user.AIVideoQuotaUsed,
		user.AIVideoQuotaLimit,
		user.AutoPostingQuotaUsed,
		user.AutoPostingQuotaLimit,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

func (r *PostgresUserRepository) UpdateQuotas(ctx context.Context, userID string, quotas sharedDomain.QuotaInfo) error {
	query := `
		UPDATE users 
		SET ai_description_quota_used = $1,
			ai_video_quota_used = $2,
			auto_posting_quota_used = $3,
			updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		quotas.AIDescription.Used,
		quotas.AIVideo.Used,
		quotas.AutoPosting.Used,
		time.Now().UTC(),
		userID,
	)

	return err
}

func (r *PostgresUserRepository) ResetAllMonthlyQuotas(ctx context.Context) error {
	query := `
		UPDATE users 
		SET ai_description_quota_used = 0,
			ai_video_quota_used = 0,
			auto_posting_quota_used = 0,
			updated_at = $1
	`

	_, err := r.db.ExecContext(ctx, query, time.Now().UTC())
	return err
}
