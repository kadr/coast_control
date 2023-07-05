package product

import (
	"context"
	"github.com/cost_control/internal/handlers/utils"
	"github.com/cost_control/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type IProductService interface {
	Create(ctx context.Context, product service.ProductServiceInput) (string, error)
	Update(ctx context.Context, id string, product service.ProductServiceInput) error
	GetAll(ctx context.Context, filter interface{}) ([]service.ProductServiceOutput, error)
	GetById(ctx context.Context, id string) (service.ProductServiceOutput, error)
	Delete(ctx context.Context, id string) error
}

type ProductApiHandler struct {
	productService IProductService
	Response       utils.Response
}

func New(productService IProductService) *ProductApiHandler {
	return &ProductApiHandler{productService: productService}
}

func (pah ProductApiHandler) Create(c *gin.Context) {
	var dto CreateProductDTO
	if err := c.BindJSON(&dto); err != nil {
		pah.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	productInput := service.ProductServiceInput{
		Name:        dto.Name,
		Price:       dto.Price,
		Description: dto.Description,
		BuyAt:       dto.BuyAt,
		User:        dto.User,
	}
	id, err := pah.productService.Create(context.Background(), productInput)
	if err != nil {
		pah.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pah.Response.Success(c, http.StatusOK, map[string]string{"id": id})
}

func (pah ProductApiHandler) Update(c *gin.Context) {
	productInput := service.ProductServiceInput{}
	if err := c.BindJSON(&productInput); err != nil {
		pah.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	err := pah.productService.Update(context.Background(), c.Param("id"), productInput)
	if err != nil {
		pah.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pah.Response.Success(c, http.StatusOK, nil)
}
func (pah ProductApiHandler) GetById(c *gin.Context) {
	product, err := pah.productService.GetById(context.Background(), c.Param("id"))
	if err != nil {
		pah.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pah.Response.Success(c, http.StatusOK, product)
}
func (pah ProductApiHandler) Get(c *gin.Context) {
	products, err := pah.productService.GetAll(context.Background(), bson.D{})
	if err != nil {
		pah.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pah.Response.Success(c, http.StatusOK, products)
}
func (pah ProductApiHandler) Delete(c *gin.Context) {
	err := pah.productService.Delete(context.Background(), c.Param("id"))
	if err != nil {
		pah.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pah.Response.Success(c, http.StatusOK, nil)
}

func (pah ProductApiHandler) Report(c *gin.Context) {
	filter, err := prepareFilter(c)
	if err != nil {
		return
	}
	products, err := pah.productService.GetAll(context.Background(), filter)
	if err != nil {
		pah.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var sum float32
	result := make(map[string]float32)
	for _, product := range products {
		sum += product.Price
		result[product.User] += product.Price
	}
	result["sum"] = sum

	pah.Response.Success(c, http.StatusOK, result)

}

func prepareFilter(ctx *gin.Context) (filter primitive.M, err error) {
	var from, to time.Time
	if _from, ok := ctx.GetQuery("from"); ok {
		from, err = time.ParseInLocation(dateTimeFormatJSONWithoutTime, _from, time.Local)
		if err != nil {
			return
		}
	}
	if _to, ok := ctx.GetQuery("to"); ok {
		to, err = time.ParseInLocation(dateTimeFormatJSONWithoutTime, _to, time.Local)
		if err != nil {
			return
		}
	}
	filter = bson.M{}
	switch {
	case !from.IsZero() && !to.IsZero():
		filter["buy_at"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(from),
			"$lt":  primitive.NewDateTimeFromTime(to),
		}
	case !from.IsZero():
		filter["buy_at"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(from),
		}
	case !to.IsZero():
		filter["buy_at"] = bson.M{
			"$lt": primitive.NewDateTimeFromTime(to),
		}
	}

	return
}
