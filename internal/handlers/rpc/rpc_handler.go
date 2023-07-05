package rpc

import (
	"github.com/cost_control/internal/handlers/rpc/product"
	pb "github.com/cost_control/internal/handlers/rpc/src"
	productRepos "github.com/cost_control/internal/repository/product"
	"github.com/cost_control/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type RpcHandler struct {
	productRpcServer product.ProductRpcServer
}

func New(db *mongo.Collection) RpcHandler {
	repos := productRepos.New(db)
	return RpcHandler{productRpcServer: product.New(service.New(repos))}
}

func (rh RpcHandler) Start() error {
	//TODO: Address and port from env variable
	address := ":5300"
	listener, err := net.Listen("tcp", address)
	log.Printf("Start rpc server on asddress: %s", address)
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	pb.RegisterProductServicesServer(grpcServer, &rh.productRpcServer)
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
