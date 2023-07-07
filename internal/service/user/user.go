package user

import (
	"context"
	"github.com/cost_control/internal/models"
)

type IUserRepository interface {
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type UserService struct {
	repository IUserRepository
}

func New(repository IUserRepository) *UserService {
	return &UserService{repository: repository}
}

func (u UserService) GetByEmail(ctx context.Context, email string) (UserServiceOutput, error) {
	user, err := u.repository.GetByEmail(ctx, email)
	if err != nil {
		return UserServiceOutput{}, err
	}
	return UserServiceOutput{Id: user.Id, Name: user.Name, Email: user.Email, Password: user.Password}, nil
}
