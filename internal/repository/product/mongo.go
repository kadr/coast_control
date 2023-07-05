package product

import (
	"context"
	"github.com/cost_control/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type ProductMongo struct {
	db *mongo.Collection
}

func New(db *mongo.Collection) *ProductMongo {
	return &ProductMongo{db: db}
}

func (pm ProductMongo) Create(ctx context.Context, product models.Product) (string, error) {
	result, err := pm.db.InsertOne(ctx, product)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (pm ProductMongo) Update(ctx context.Context, id string, product models.Product) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	updatedProduct := bson.M{}
	if product.Name != "" {
		updatedProduct["name"] = product.Name
	}
	if product.Price != 0 {
		updatedProduct["price"] = product.Price
	}
	if product.Description != "" {
		updatedProduct["description"] = product.Description
	}
	if !product.BuyAt.IsZero() {
		updatedProduct["buy_at"] = product.BuyAt.String()
	}
	if product.User != "" {
		updatedProduct["user"] = product.User
	}
	update := bson.D{{"$set", updatedProduct}}

	_, err = pm.db.UpdateOne(ctx, bson.D{{"_id", objectID}}, update)
	if err != nil {
		return err
	}

	return nil
}

func (pm ProductMongo) GetAll(ctx context.Context, filter interface{}) ([]models.Product, error) {
	var products []models.Product
	options := options.Find()
	// Sort by `_id` field descending
	options.SetSort(bson.D{{"buy_at", -1}})
	options.SetLimit(100)
	cursor, err := pm.db.Find(ctx, filter, options)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (pm ProductMongo) GetById(ctx context.Context, id string) (models.Product, error) {
	var product models.Product
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Product{}, err
	}
	err = pm.db.FindOne(ctx, bson.D{{"_id", objectID}}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(err)
			return models.Product{}, nil
		}
		return models.Product{}, err
	}

	return product, nil
}

func (pm ProductMongo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = pm.db.DeleteOne(ctx, bson.D{{"_id", objectID}})
	if err != nil {
		return err
	}

	return nil
}
