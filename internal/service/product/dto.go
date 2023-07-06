package product

import (
	"github.com/cost_control/internal/models"
	"time"
)

type ProductServiceInput struct {
	Name        string
	Price       float32
	Description string
	BuyAt       time.Time
	User        string
}

type ProductServiceOutput struct {
	Id          string
	Name        string
	Price       float32
	Description string
	BuyAt       time.Time
	User        string
}

func CreateProductDb(product ProductServiceInput) (models.Product, error) {
	return models.NewProduct(
		product.Name,
		product.Description,
		product.Price,
		product.BuyAt,
		product.User,
	)
}

func UpdateProductDb(product ProductServiceInput) models.Product {
	return models.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		BuyAt:       product.BuyAt,
		User:        product.User,
	}
}
