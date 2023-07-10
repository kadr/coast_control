package auth

import (
	"context"
	"errors"
	"github.com/cost_control/config"
	pb "github.com/cost_control/internal/handlers/rpc/src"
	"github.com/cost_control/internal/service/user"
	"github.com/cost_control/pkg/jwt"
	"github.com/cost_control/pkg/password_hasher"
)

const hour = 60 * 24

type IUserService interface {
	GetByEmail(ctx context.Context, email string) (user.UserServiceOutput, error)
}

type AuthRpcServer struct {
	pb.AuthServiceServer
	userService   IUserService
	cfg           *config.Config
	jwtManager    *jwt.Token
	hasherManager *password_hasher.PasswordHasher
}

func New(userService IUserService, cfg *config.Config) AuthRpcServer {
	return AuthRpcServer{
		userService:   userService,
		cfg:           cfg,
		jwtManager:    jwt.New(),
		hasherManager: password_hasher.New(),
	}
}

func (s *AuthRpcServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	email, password := request.Email, request.Password
	findUser, err := s.userService.GetByEmail(ctx, email)
	response := pb.LoginResponse{}
	if err != nil {
		response.Error = true
		message := err.Error()
		response.ErrorMessage = &message
		return &response, err
	}
	if !s.hasherManager.Verify(findUser.Password, password) {
		return nil, errors.New("переданн не корректрный пароль")
	}
	generatedToken, err := s.jwtManager.Generate(email, hour, s.cfg.SignedKey)
	if err != nil {
		return nil, err
	}
	response.Token = &generatedToken

	return &response, nil
}
