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

type IUserService interface {
	GetByEmail(ctx context.Context, email string) (user.UserServiceOutput, error)
}

type AuthRpcServer struct {
	pb.AuthServiceServer
	userService IUserService
	cfg         *config.Config
}

func New(userService IUserService, cfg *config.Config) AuthRpcServer {
	return AuthRpcServer{userService: userService, cfg: cfg}
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
	if !password_hasher.Verify(findUser.Password, password) {
		return nil, errors.New("переданн не корректрный пароль")
	}
	token := jwt.New(email, 60*60*24)
	generatedToken, err := token.Generate(s.cfg.SignedKey)
	if err != nil {
		return nil, err
	}
	response.Token = &generatedToken

	return &response, nil
}
