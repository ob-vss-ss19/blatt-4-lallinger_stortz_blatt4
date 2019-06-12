package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type User struct{}

func (User) Request(context.Context, *proto.UserRequest, *proto.UserResponse) error {
	panic("implement me")
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
