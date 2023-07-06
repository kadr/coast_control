package product

import (
	"context"
	pb "github.com/cost_control/internal/handlers/rpc/src"
	"github.com/cost_control/internal/service/product"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IProductService interface {
	Create(ctx context.Context, product product.ProductServiceInput) (string, error)
	Update(ctx context.Context, id string, product product.ProductServiceInput) error
	GetAll(ctx context.Context, filter interface{}) ([]product.ProductServiceOutput, error)
	GetById(ctx context.Context, id string) (product.ProductServiceOutput, error)
	Delete(ctx context.Context, id string) error
}

type ProductRpcServer struct {
	pb.ProductServicesServer
	productService IProductService
}

func New(productService IProductService) ProductRpcServer {
	return ProductRpcServer{productService: productService}
}

func (s *ProductRpcServer) Add(ctx context.Context, request *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	dto := CreateProductInput(request)
	id, err := s.productService.Create(ctx, dto)
	response := pb.CreateProductResponse{}
	if err != nil {
		response.Error = true
		message := err.Error()
		response.ErrorMessage = &message
		return &response, err
	}
	response.Id = &id

	return &response, nil
}

func (s *ProductRpcServer) Update(ctx context.Context, request *pb.UpdateProductRequest) (*pb.Response, error) {
	dto := UpdateProductInput(request)
	err := s.productService.Update(ctx, request.Id, dto)
	response := pb.Response{}
	if err != nil {
		response.Error = true
		message := err.Error()
		response.ErrorMessage = &message
		return &response, err
	}

	return &response, nil
}

func (s *ProductRpcServer) Get(ctx context.Context, request *pb.ProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.productService.GetById(ctx, request.Id)
	response := pb.GetProductResponse{}
	if err != nil {
		response.Error = true
		message := err.Error()
		response.ErrorMessage = &message
		return &response, err
	}
	response.Product = &pb.Product{
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		BuyAt:       timestamppb.New(product.BuyAt),
		User:        product.User}

	return &response, nil
}

func (s *ProductRpcServer) Search(ctx context.Context, filter *pb.Filter) (*pb.SearchProductResponse, error) {
	products, err := s.productService.GetAll(ctx, filter)
	response := pb.SearchProductResponse{}
	if err != nil {
		response.Error = true
		message := err.Error()
		response.ErrorMessage = &message
		return &response, err
	}
	for _, product := range products {
		response.Products = append(response.Products, &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Price:       product.Price,
			Description: product.Description,
			BuyAt:       timestamppb.New(product.BuyAt),
			User:        product.User,
		})
	}

	return &response, nil
}

func (s *ProductRpcServer) Delete(ctx context.Context, request *pb.ProductRequest) (*pb.Response, error) {
	err := s.productService.Delete(ctx, request.Id)
	response := pb.Response{}
	if err != nil {
		response.Error = true
		message := err.Error()
		response.ErrorMessage = &message
		return &response, err
	}

	return &response, err
}

func (s *ProductRpcServer) Report(ctx context.Context, filter *pb.Filter) (*pb.ReportResponse, error) {
	products, err := s.productService.GetAll(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var sum float32
	response := pb.ReportResponse{}
	result := make(map[string]float32)
	for _, product := range products {
		sum += product.Price
		result[product.User] += product.Price
	}
	response.Sum = sum
	response.Period = &pb.Period{
		From: filter.From,
		To:   filter.To,
	}
	for user, sum := range result {
		response.GroupByUsers = append(response.GroupByUsers, &pb.GroupByUsers{
			User: user,
			Sum:  sum,
		})
	}

	return &response, nil
}

func (s *ProductRpcServer) mustEmbedUnimplementedProductServicesServer() {
	//TODO implement me
	panic("implement me")
}
