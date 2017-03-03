package models

import (
	"context"

	"github.com/chideat/glog"
	. "github.com/chideat/pcc/action/modules/config"
	"google.golang.org/grpc"
)

var (
	userRPCConn *grpc.ClientConn
	userClient  UserRPCClient
)

func init() {
	var err error

	userRPCConn, err = grpc.Dial(Conf.RPC.UserRPCAddr, grpc.WithInsecure())
	if err != nil {
		glog.Panic(err)
	}
	userClient = NewUserRPCClient(userRPCConn)
}

func (user *User) Info() map[string]interface{} {
	info := map[string]interface{}{}
	info["id"] = user.Id
	info["name"] = user.Name

	return info
}

func GetUserById(id uint64) (*User, error) {
	return userClient.GetUserById(context.Background(), &User{Id: id})
}
