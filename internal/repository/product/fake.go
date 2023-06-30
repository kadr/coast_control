package product

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cost_control/internal/models"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
)

const filePath = "fake_db.json"

type ProductFake struct{}

func NewFakeDb() *ProductFake {
	return &ProductFake{}
}

func (pm ProductFake) Create(ctx context.Context, product models.Product) (string, error) {
	flag := os.O_RDWR
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		flag = os.O_RDWR | os.O_CREATE
	}
	file, err := os.OpenFile(filePath, flag, 0755)
	if err != nil {
		return "", err
	}
	defer file.Close()
	storage := make(map[string]models.Product)
	id := uuid.NewString()
	storage[id] = product
	jsonString, err := json.Marshal(storage)
	if err != nil {
		return "", err
	}
	_, err = file.Write(jsonString)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (pm ProductFake) Update(ctx context.Context, id string, product models.Product) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var storage map[string]models.Product
	err = json.Unmarshal(file, &storage)
	if err != nil {
		return err
	}
	updatedProduct := storage[id]
	if product.Name != "" {
		updatedProduct.Name = product.Name
	}
	if product.Price != 0 {
		updatedProduct.Price = product.Price
	}
	if product.Description != "" {
		updatedProduct.Description = product.Description
	}
	if !product.BuyAt.IsZero() {
		updatedProduct.BuyAt = product.BuyAt
	}
	storage[id] = updatedProduct
	newFileData, err := json.Marshal(storage)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, newFileData, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (pm ProductFake) GetAll(ctx context.Context, filter interface{}) ([]models.Product, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var storage map[string]models.Product
	err = json.Unmarshal(file, &storage)
	if err != nil {
		return nil, err
	}
	var products []models.Product
	for _, product := range storage {
		products = append(products, product)
	}
	return products, nil
}
func (pm ProductFake) GetById(ctx context.Context, id string) (models.Product, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return models.Product{}, err
	}
	var storage map[string]models.Product
	err = json.Unmarshal(file, &storage)
	if err != nil {
		return models.Product{}, err
	}

	return storage[id], nil
}

func (pm ProductFake) Delete(ctx context.Context, id string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var storage map[string]models.Product
	err = json.Unmarshal(file, &storage)
	if err != nil {
		return err
	}
	delete(storage, id)
	newFileData, err := json.Marshal(storage)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, newFileData, 0755)
	if err != nil {
		return err
	}
	return nil
}
