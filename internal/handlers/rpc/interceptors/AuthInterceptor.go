package interceptor

import (
	"context"
	"fmt"
	"github.com/cost_control/config"
	pb "github.com/cost_control/internal/handlers/rpc/src"
	"github.com/cost_control/pkg/jwt"
	"github.com/cost_control/pkg/password_hasher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

type AuthInterceptor struct {
	cfg           *config.Config
	jwtManager    *jwt.Token
	hasherManager *password_hasher.PasswordHasher
}

func New(cfg *config.Config) *AuthInterceptor {
	return &AuthInterceptor{cfg: cfg, jwtManager: jwt.New(), hasherManager: password_hasher.New()}
}

func (ai AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if isLoginRequest(req, info) {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("не переданы данные для авторизациия. %s", info.FullMethod)
		}
		var token []string
		if token, ok = md["authorization"]; !ok {
			return nil, fmt.Errorf("не передан токен. %s", info.FullMethod)
		}
		if !ai.jwtManager.IsValid(token[0], ai.cfg.SignedKey) {
			return nil, fmt.Errorf("токен не валиден. %s", info.FullMethod)
		}
		m, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}
		return m, err
	}
}

func (ai AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return fmt.Errorf("не переданы данные для авторизациия. %s", info.FullMethod)
		}
		var token []string
		if token, ok = md["authorization"]; !ok {
			return fmt.Errorf("не передан токен. %s", info.FullMethod)
		}
		if !ai.jwtManager.IsValid(token[0], ai.cfg.SignedKey) {
			return fmt.Errorf("токен не валиден. %s", info.FullMethod)
		}

		return handler(srv, stream)
	}
}

func isLoginRequest(request interface{}, info *grpc.UnaryServerInfo) bool {
	switch request.(type) {
	case *pb.LoginRequest:
		if request.(*pb.LoginRequest).Email != "" {
			if request.(*pb.LoginRequest).Password != "" {
				if strings.Contains(info.FullMethod, "Login") {
					return true
				}
			}
		}
	}

	return false
}
