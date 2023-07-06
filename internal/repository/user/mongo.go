package product

import (
	"context"
	"errors"
	"github.com/cost_control/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongo struct {
	db *mongo.Collection
}

func New(db *mongo.Collection) *UserMongo {
	return &UserMongo{db: db}
}

func (um UserMongo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := um.db.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, errors.New("пользователя с таким email не найдено")
		}
		return models.User{}, err
	}

	return user, nil
}

