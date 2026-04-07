package service

import (
	"context"
	"errors"
	"strings"

	"github.com/runtimeninja/importpilot/internal/domain"
	"github.com/runtimeninja/importpilot/internal/repository"
)

type ClientService struct {
	clientRepo *repository.ClientRepository
}

type CreateClientInput struct {
	Name    string
	Email   string
	ShopURL string
	Plan    string
}

func NewClientService(clientRepo *repository.ClientRepository) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
	}
}

func (s *ClientService) CreateClient(ctx context.Context, input CreateClientInput) (int64, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	shopURL := strings.TrimSpace(input.ShopURL)
	plan := strings.TrimSpace(strings.ToLower(input.Plan))

	if name == "" {
		return 0, errors.New("name is required")
	}

	if email == "" {
		return 0, errors.New("email is required")
	}

	if shopURL == "" {
		return 0, errors.New("shop_url is required")
	}

	if plan == "" {
		plan = "free"
	}

	client := domain.Client{
		Name:    name,
		Email:   email,
		ShopURL: shopURL,
		Status:  "active",
		Plan:    plan,
	}

	id, err := s.clientRepo.Create(ctx, client)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *ClientService) ListClients(ctx context.Context, limit, offset int) ([]domain.Client, error) {
	return s.clientRepo.List(ctx, limit, offset)
}

func (s *ClientService) GetClientByID(ctx context.Context, id int64) (*domain.Client, error) {
	return s.clientRepo.GetByID(ctx, id)
}

func (s *ClientService) UpdateClientStatus(ctx context.Context, id int64, status string) error {
	status = strings.TrimSpace(strings.ToLower(status))

	if id <= 0 {
		return domain.ErrInvalidClientID
	}

	if status != "active" && status != "inactive" {
		return domain.ErrInvalidClientStatus
	}

	return s.clientRepo.UpdateStatus(ctx, id, status)
}
