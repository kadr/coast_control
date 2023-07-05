package product

import (
	"context"
	"github.com/cost_control/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type IProductService interface {
	Create(ctx context.Context, product service.ProductServiceInput) (string, error)
	GetAll(ctx context.Context, filter interface{}) ([]service.ProductServiceOutput, error)
	GetById(ctx context.Context, id string) (service.ProductServiceOutput, error)
	Delete(ctx context.Context, id string) error
}

type ProductBotHandler struct {
	productService IProductService
}

func New(productService IProductService) *ProductBotHandler {
	return &ProductBotHandler{productService: productService}
}

func (pah ProductBotHandler) Create(dto CreateProductDTO) (string, error) {
	product := service.ProductServiceInput{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		BuyAt:       dto.BuyAt,
		User:        dto.User,
	}

	id, err := pah.productService.Create(context.Background(), product)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (pah ProductBotHandler) GetById(id string) (service.ProductServiceOutput, error) {
	product, err := pah.productService.GetById(context.Background(), id)
	if err != nil {
		return service.ProductServiceOutput{}, err
	}
	return product, nil
}
func (pah ProductBotHandler) Get(filter string) ([]service.ProductServiceOutput, error) {
	preparedFilter := prepareFilter(filter)
	products, err := pah.productService.GetAll(context.Background(), preparedFilter)
	if err != nil {
		return []service.ProductServiceOutput{}, err
	}

	return products, nil
}
func (pah ProductBotHandler) Delete(id string) error {
	err := pah.productService.Delete(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}

func (pah ProductBotHandler) Report(filter string) (map[string]float32, error) {
	preparedFilter := prepareFilter(filter)
	products, err := pah.productService.GetAll(context.Background(), preparedFilter)
	if err != nil {
		return nil, err
	}
	var sum float32
	result := make(map[string]float32)
	for _, product := range products {
		sum += product.Price
		result[product.User] += product.Price
	}
	result["sum"] = sum

	return result, nil
}

func prepareFilter(filter string) bson.M {
	preparedFilter := bson.M{}
	splitFilter := strings.Split(filter, " ")
	switch {
	case len(splitFilter) == 2:
		var from, to time.Time
		from, err := time.ParseInLocation(dateTimeFormatJSONWithoutTime, splitFilter[0], time.Local)
		if err != nil {
			return bson.M{}
		}
		to, err = time.ParseInLocation(dateTimeFormatJSONWithoutTime, splitFilter[1], time.Local)
		if err != nil {
			return bson.M{}
		}
		preparedFilter["buy_at"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(from),
			"$lte": primitive.NewDateTimeFromTime(to),
		}
	case len(splitFilter) == 1:
		date, err := time.ParseInLocation(dateTimeFormatJSONWithoutTime, splitFilter[0], time.Local)
		if err != nil {
			return bson.M{}
		}
		if date.Before(time.Now()) {
			preparedFilter["buy_at"] = bson.M{"$gte": primitive.NewDateTimeFromTime(date)}
		} else {
			preparedFilter["buy_at"] = bson.M{"$lte": primitive.NewDateTimeFromTime(date)}
		}
	}

	return preparedFilter
}
