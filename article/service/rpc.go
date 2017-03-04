package service

import (
	"fmt"
	"net"

	"github.com/chideat/pcc/article/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RPCService struct{}

func (serv *RPCService) GetArticleById(ctx context.Context, article *models.Article) (*models.Article, error) {
	article, err := models.GetArticleById(article.Id)
	if err != nil {
		return nil, err
	} else if article == nil {
		return nil, fmt.Errorf("article not exists")
	}
	return article, nil
}

func StartRPCService(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	service := RPCService{}
	server := grpc.NewServer()
	models.RegisterArticleRPCServer(server, &service)
	reflection.Register(server)

	return server.Serve(listen)
}
