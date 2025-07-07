package repository

import (
	"context"
	"ops-monorepo/services/svc-user/internal/model"
	pg "ops-monorepo/shared-libs/storage/postgres"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type IUserSQLRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID string) error
}

type userSQLRepository struct {
	db *pg.PostgresPgx
}

func NewUserRepository(db *pg.PostgresPgx) IUserSQLRepository {
	return &userSQLRepository{
		db: db,
	}
}

func (r *userSQLRepository) CreateUser(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, password, first_name, last_name, roles, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Pool().Exec(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		pq.Array(user.Roles),
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *userSQLRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	var user model.User
	err := r.db.Pool().QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userSQLRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, roles, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	var user model.User
	err := r.db.Pool().QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		pq.Array(&user.Roles),
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userSQLRepository) UpdateUser(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users 
		SET email = $2, password = $3, first_name = $4, last_name = $5, roles = $6, is_active = $7, updated_at = $8
		WHERE id = $1
	`

	_, err := r.db.Pool().Exec(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		pq.Array(user.Roles),
		user.IsActive,
		user.UpdatedAt,
	)

	return err
}

func (r *userSQLRepository) CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error {
	if refreshToken.ID == "" {
		refreshToken.ID = uuid.New().String()
	}
	refreshToken.CreatedAt = time.Now()

	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, is_revoked)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Pool().Exec(ctx, query,
		refreshToken.ID,
		refreshToken.UserID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
		refreshToken.IsRevoked,
	)

	return err
}

func (r *userSQLRepository) GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, is_revoked
		FROM refresh_tokens
		WHERE token = $1 AND is_revoked = false AND expires_at > NOW()
	`

	var refreshToken model.RefreshToken
	err := r.db.Pool().QueryRow(ctx, query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.IsRevoked,
	)

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *userSQLRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true
		WHERE token = $1
	`

	_, err := r.db.Pool().Exec(ctx, query, token)
	return err
}

func (r *userSQLRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true
		WHERE user_id = $1
	`

	_, err := r.db.Pool().Exec(ctx, query, userID)
	return err
}
