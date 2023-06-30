package product

import (
	"context"
	"github.com/cost_control/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type IProductService interface {
	Create(ctx context.Context, product models.Product) (string, error)
	Update(ctx context.Context, id string, product models.Product) error
	GetAll(ctx context.Context, filter interface{}) ([]models.Product, error)
	GetById(ctx context.Context, id string) (models.Product, error)
	Delete(ctx context.Context, id string) error
}

type ProductBotHandler struct {
	productService IProductService
}

func New(productService IProductService) *ProductBotHandler {
	return &ProductBotHandler{productService: productService}
}

func (pah ProductBotHandler) Create(dto CreateProductDTO) (string, error) {
	product, err := models.NewProduct(
		dto.Name,
		dto.Description,
		dto.Price,
		dto.BuyAt,
		dto.User,
	)
	if err != nil {
		return "", err
	}
	id, err := pah.productService.Create(context.Background(), product)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (pah ProductBotHandler) Update(dto UpdateProductDTO) error {
	product := models.Product{}
	if len(dto.Name) > 0 {
		product.Name = dto.Name
	}
	if dto.Price > 0 {
		product.Price = dto.Price
	}
	if len(dto.Description) > 0 {
		product.Description = dto.Description
	}
	if dto.BuyAt != nil {
		product.BuyAt = *dto.BuyAt
	}
	err := pah.productService.Update(context.Background(), dto.Id, product)
	if err != nil {
		return err
	}

	return nil
}
func (pah ProductBotHandler) GetById(id string) (models.Product, error) {
	product, err := pah.productService.GetById(context.Background(), id)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}
func (pah ProductBotHandler) Get(filter string) ([]models.Product, error) {
	preparedFilter := prepareFilter(filter)
	products, err := pah.productService.GetAll(context.Background(), preparedFilter)
	if err != nil {
		return []models.Product{}, err
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
			"$lt":  primitive.NewDateTimeFromTime(to),
		}
	case len(splitFilter) == 1:
		date, err := time.ParseInLocation(dateTimeFormatJSONWithoutTime, splitFilter[0], time.Local)
		if err != nil {
			return bson.M{}
		}
		if date.Before(time.Now()) {
			preparedFilter["buy_at"] = bson.M{"$gte": primitive.NewDateTimeFromTime(date)}
		} else {
			preparedFilter["buy_at"] = bson.M{"$lt": primitive.NewDateTimeFromTime(date)}
		}
	}

	return preparedFilter
}
