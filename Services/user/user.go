package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type User struct {
	users map[string][]int
}

func (me *User) CreateUser(ctx context.Context, req *proto.UserData, rsp *proto.Response) error {
	return nil
}
func (me *User) DeleteUser(ctx context.Context, req *proto.UserData, rsp *proto.Response) error {
	return nil
}
func (me *User) GetUsers(ctx context.Context, req *proto.UserRequest, rsp *proto.UserResponse) error {
	return nil
}

func main() {
	service := micro.NewService(micro.Name("user"))
	service.Init()
	proto.RegisterUserHandler(service.Server(), new(User))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
