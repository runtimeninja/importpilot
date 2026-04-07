package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/runtimeninja/importpilot/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	query := `
		INSERT INTO users (client_id, email, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.ClientID,
		user.Email,
		user.PasswordHash,
		user.IsActive,
	).Scan(&id)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "users_email_key") {
			return 0, errors.New("user email already exists")
		}
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, client_id, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.ClientID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetRolesByUserID(ctx context.Context, userID int64) ([]string, error) {
	query := `
		SELECT roles.name
		FROM user_roles
		INNER JOIN roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = $1
		ORDER BY roles.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var roleName string
		if err := rows.Scan(&roleName); err != nil {
			return nil, err
		}
		roles = append(roles, roleName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
