package repository

import (
	"context"
	"database/sql"

	"github.com/runtimeninja/importpilot/internal/domain"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Create(ctx context.Context, client domain.Client) (int64, error) {
	query := `
		INSERT INTO clients (name, email, shop_url, status, plan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		client.Name,
		client.Email,
		client.ShopURL,
		client.Status,
		client.Plan,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	query := `
		SELECT id, name, email, shop_url, status, plan, created_at, updated_at
		FROM clients
		WHERE id = $1
	`

	var client domain.Client
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.ShopURL,
		&client.Status,
		&client.Plan,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) List(ctx context.Context, limit, offset int) ([]domain.Client, error) {
	query := `
		SELECT id, name, email, shop_url, status, plan, created_at, updated_at
		FROM clients
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []domain.Client

	for rows.Next() {
		var client domain.Client
		if err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.Email,
			&client.ShopURL,
			&client.Status,
			&client.Plan,
			&client.CreatedAt,
			&client.UpdatedAt,
		); err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *ClientRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE clients
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
