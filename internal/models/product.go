package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

const dateTimeFormatJSON = "02.01.2006 15:04"
const dateTimeFormatBSON = "2006-01-02 15:04:05 +0400 +04"

type Product struct {
	Id          string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Price       float32   `json:"price" bson:"price"`
	BuyAt       time.Time `json:"buy_at,omitempty" bson:"buy_at"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	User        string    `json:"user,omitempty" bson:"user,omitempty"`
}

func NewProduct(name, description string, price float32, byuAt *time.Time, user string) (Product, error) {
	var err error
	var byuAtDate time.Time
	if byuAt == nil {
		byuAtDate, err = time.ParseInLocation(dateTimeFormatBSON, time.Now().Local().String(), time.Local)
	} else {
		byuAtDate, err = time.ParseInLocation(dateTimeFormatBSON, byuAt.String(), time.Local)
	}
	if err != nil {
		log.Println(err)
		return Product{}, err
	}
	return Product{
		Name:        name,
		Price:       price,
		BuyAt:       byuAtDate,
		Description: description,
		User:        user,
	}, nil
}

func (p *Product) UnmarshalBSON(bs []byte) error {
	var result map[string]interface{}
	err := bson.Unmarshal(bs, &result)
	if err != nil {
		return err
	}
	if val, ok := result["buy_at"]; ok {
		var t time.Time
		switch val.(type) {
		case string:
			t, err = time.ParseInLocation(dateTimeFormatBSON, val.(string), time.Local)
		case primitive.DateTime:
			t = val.(primitive.DateTime).Time()
		}
		if err != nil {
			return err
		}
		p.BuyAt = t

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
	if val, ok := result["user"]; ok {
		p.User = val.(string)
	}
	return nil
}
