package product

import (
	pb "github.com/cost_control/internal/handlers/rpc/src"
	"github.com/cost_control/internal/service/product"
)

func CreateProductInput(request *pb.CreateProductRequest) product.ProductServiceInput {
	product := product.ProductServiceInput{
		Name:  request.Name,
		Price: request.Price,
		User:  request.User,
	}
	if request.Description != nil {
		product.Description = *request.Description
	}
	if request.BuyAt != nil {
		product.BuyAt = request.BuyAt.AsTime()
	}

	return product
}

func UpdateProductInput(request *pb.UpdateProductRequest) product.ProductServiceInput {
	product := product.ProductServiceInput{}
	if request.Name != nil {
		product.Name = *request.Name
	}
	if request.Price != nil {
		product.Price = *request.Price
	}
	if request.Description != nil {
		product.Description = *request.Description
	}
	if request.BuyAt != nil {
		product.BuyAt = request.BuyAt.AsTime()
	}
	if request.User != nil {
		product.User = *request.User
	}

	return product
}
