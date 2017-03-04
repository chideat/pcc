package models

import (
	"context"

	"github.com/chideat/glog"
	. "github.com/chideat/pcc/action/modules/config"
	"google.golang.org/grpc"
)

var (
	articleRPCConn *grpc.ClientConn
	articleClient  ArticleRPCClient
)

func init() {
	var err error

	articleRPCConn, err = grpc.Dial(Conf.RPC.ArticleRPCAddr, grpc.WithInsecure())
	if err != nil {
		glog.Panic(err)
	}
	articleClient = NewArticleRPCClient(articleRPCConn)
}

func GetArticleById(id uint64) (*Article, error) {
	return articleClient.GetArticleById(context.Background(), &Article{Id: id})
}
