package service

import (
	"context"
	"github.com/cost_control/internal/models"
)

type IProductRepository interface {
	Create(ctx context.Context, product models.Product) (string, error)
	Update(ctx context.Context, id string, product models.Product) error
	GetAll(ctx context.Context, filter interface{}) ([]models.Product, error)
	GetById(ctx context.Context, id string) (models.Product, error)
	Delete(ctx context.Context, id string) error
}

type ProductService struct {
	repository IProductRepository
}

func New(repository IProductRepository) *ProductService {
	return &ProductService{repository: repository}
}

func (ps ProductService) Create(ctx context.Context, product models.Product) (string, error) {
	result, err := ps.repository.Create(ctx, product)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (ps ProductService) Update(ctx context.Context, id string, product models.Product) error {
	err := ps.repository.Update(ctx, id, product)
	if err != nil {
		return err
	}
	return nil
}

func (ps ProductService) GetAll(ctx context.Context, filter interface{}) ([]models.Product, error) {
	products, err := ps.repository.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps ProductService) GetById(ctx context.Context, id string) (models.Product, error) {
	product, err := ps.repository.GetById(ctx, id)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}
func (ps ProductService) Delete(ctx context.Context, id string) error {
	err := ps.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
