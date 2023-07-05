package product

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateTimeFormatJSON = "02.01.2006 15:04"
const dateTimeFormatJSONWithoutTime = "02.01.2006"

type CreateProductDTO struct {
	Name        string    `json:"name"`
	Price       float32   `json:"price"`
	BuyAt       time.Time `json:"buy_at,omitempty"`
	Description string    `json:"description,omitempty"`
	User        string    `json:"user"`
}

type UpdateProductDTO struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Price       float32   `json:"price,omitempty"`
	BuyAt       time.Time `json:"buy_at,omitempty"`
	Description string    `json:"description,omitempty"`
}

type GetProductDTO struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Price       float32   `json:"price,omitempty"`
	BuyAt       time.Time `json:"buy_at,omitempty"`
	Description string    `json:"description,omitempty"`
	User        string    `json:"user"`
}

func (p *CreateProductDTO) Mapping(product string, user string) error {
	splitProduct := strings.Split(product, "\n")
	if splitProduct[0] == "" {
		splitProduct = append(splitProduct[0:], splitProduct[1:]...)
	}
	if len(splitProduct) < 3 {
		return errors.New("Не все обязательные поля были переданы.")
	}
	p.Name = splitProduct[0]
	price, err := strconv.ParseFloat(splitProduct[1], 32)
	if err != nil {
		return errors.New("Передана не корректная цена.")
	}
	p.Price = float32(price)
	p.Description = splitProduct[2]
	p.User = user
	if len(splitProduct) == 4 {
		p.BuyAt, err = time.ParseInLocation(dateTimeFormatJSON, splitProduct[3], time.Local)
		if err != nil {
			return err
		}
	}

	return nil
}
func (p *UpdateProductDTO) Mapping(product string) error {
	if len(product) > 0 {
		splitProduct := strings.Split(product, "\n")
		if splitProduct[0] == "" {
			splitProduct = append(splitProduct[:0], splitProduct[1:]...)
		}
		p.Id = splitProduct[0]
		switch len(splitProduct) {
		case 5:
			p.Name = splitProduct[1]
			price, err := strconv.ParseFloat(splitProduct[2], 32)
			if err != nil {
				return errors.New("Передана не корректная цена.")
			}
			p.Price = float32(price)
			p.Description = splitProduct[3]
			p.BuyAt, err = time.ParseInLocation(dateTimeFormatJSON, splitProduct[4], time.Local)
			if err != nil {
				return err
			}
		case 4:
			p.Name = splitProduct[1]
			price, err := strconv.ParseFloat(splitProduct[2], 32)
			if err != nil {
				return errors.New("Передана не корректная цена.")
			}
			p.Price = float32(price)
			p.Description = splitProduct[3]
		case 3:
			p.Name = splitProduct[1]
			price, err := strconv.ParseFloat(splitProduct[2], 32)
			if err != nil {
				return errors.New("Передана не корректная цена.")
			}
			p.Price = float32(price)
		case 2:
			p.Name = splitProduct[1]
		}
	}

	return nil
}
