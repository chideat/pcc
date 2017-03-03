package service

import (
	"net"

	"github.com/chideat/pcc/user/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RPCService struct{}

func (serv *RPCService) GetUserById(ctx context.Context, user *models.User) (*models.User, error) {
	return models.GetUserById(user.Id)
}

func (serv *RPCService) GetUserByName(ctx context.Context, user *models.User) (*models.User, error) {
	return models.GetUserByName(user.Name)
}

func StartRPCService(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	service := RPCService{}
	server := grpc.NewServer()
	models.RegisterUserRPCServer(server, &service)
	reflection.Register(server)

	return server.Serve(listen)
}
