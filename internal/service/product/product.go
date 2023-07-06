package product

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

func (ps ProductService) Create(ctx context.Context, productInput ProductServiceInput) (string, error) {
	newProduct, err := CreateProductDb(productInput)
	if err != nil {
		return "", err
	}
	result, err := ps.repository.Create(ctx, newProduct)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (ps ProductService) Update(ctx context.Context, id string, productInput ProductServiceInput) error {
	updateProduct := UpdateProductDb(productInput)
	err := ps.repository.Update(ctx, id, updateProduct)
	if err != nil {
		return err
	}
	return nil
}

func (ps ProductService) GetAll(ctx context.Context, filter interface{}) ([]ProductServiceOutput, error) {
	products, err := ps.repository.GetAll(ctx, filter)
	var productsOutput []ProductServiceOutput
	for _, product := range products {
		productsOutput = append(productsOutput, ProductServiceOutput{
			Id:          product.Id,
			Name:        product.Name,
			Price:       product.Price,
			Description: product.Description,
			BuyAt:       product.BuyAt,
			User:        product.User,
		})
	}
	if err != nil {
		return nil, err
	}
	return productsOutput, nil
}

func (ps ProductService) GetById(ctx context.Context, id string) (ProductServiceOutput, error) {
	product, err := ps.repository.GetById(ctx, id)
	if err != nil {
		return ProductServiceOutput{}, err
	}
	return ProductServiceOutput{
		Id:          product.Id,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		BuyAt:       product.BuyAt,
		User:        product.User,
	}, nil
}
func (ps ProductService) Delete(ctx context.Context, id string) error {
	err := ps.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
