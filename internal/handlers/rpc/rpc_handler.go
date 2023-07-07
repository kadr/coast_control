package rpc

import (
	"github.com/cost_control/config"
	"github.com/cost_control/internal/handlers/rpc/auth"
	interceptor "github.com/cost_control/internal/handlers/rpc/interceptors"
	"github.com/cost_control/internal/handlers/rpc/product"
	pb "github.com/cost_control/internal/handlers/rpc/src"
	productRepos "github.com/cost_control/internal/repository/product"
	userRepos "github.com/cost_control/internal/repository/user"
	productService "github.com/cost_control/internal/service/product"
	userService "github.com/cost_control/internal/service/user"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const productCollection = "product"
const userCollection = "user"

type RpcHandler struct {
	productRpcServer product.ProductRpcServer
	authRpcServer    auth.AuthRpcServer
	config           *config.Config
}

func New(db *mongo.Database, config *config.Config) RpcHandler {
	productRepo := productRepos.New(db.Collection(productCollection))
	userRepo := userRepos.New(db.Collection(userCollection))
	return RpcHandler{
		productRpcServer: product.New(productService.New(productRepo)),
		authRpcServer:    auth.New(userService.New(userRepo), config),
		config:           config,
	}
}

func (rh RpcHandler) Start() error {
	listener, err := net.Listen("tcp", rh.config.Rpc.Address)
	log.Printf("Start rpc server on asddress: %s", rh.config.Rpc.Address)
	if err != nil {
		return err
	}
	authInterceptor := interceptor.New(rh.config)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)
	//userGrpcServer := grpc.NewServer()

	pb.RegisterProductServicesServer(grpcServer, &rh.productRpcServer)
	pb.RegisterAuthServiceServer(grpcServer, &rh.authRpcServer)

	reflection.Register(grpcServer)
	//reflection.Register(userGrpcServer)
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
