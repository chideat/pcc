package service

import (
	"fmt"
	"net"

	"github.com/chideat/pcc/user/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RPCService struct{}

func (serv *RPCService) GetUserById(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := models.GetUserById(user.Id)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, fmt.Errorf("user not exists")
	}
	return user, nil
}

func (serv *RPCService) GetUserByName(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := models.GetUserByName(user.Name)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, fmt.Errorf("user not exists")
	}
	return user, nil
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
