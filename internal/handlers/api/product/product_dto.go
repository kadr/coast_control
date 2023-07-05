package product

import (
	"encoding/json"
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
	Name        string    `json:"name,omitempty"`
	Price       float32   `json:"price,omitempty"`
	BuyAt       time.Time `json:"buy_at,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (p *CreateProductDTO) UnmarshalJSON(bs []byte) error {
	var result map[string]interface{}
	err := json.Unmarshal(bs, &result)
	if err != nil {
		return err
	}
	if val, ok := result["buy_at"]; ok {
		p.BuyAt, err = time.ParseInLocation(dateTimeFormatJSON, val.(string), time.Local)
		if err != nil {
			return err
		}
	}
	if val, ok := result["name"]; ok {
		p.Name = val.(string)
	}
	if val, ok := result["price"]; ok {
		p.Price = float32(val.(float64))
	}
	if val, ok := result["description"]; ok {
		p.Description = val.(string)
	}
	return nil
}
