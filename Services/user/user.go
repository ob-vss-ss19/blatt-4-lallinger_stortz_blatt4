package main

import (
	"context"
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	proto "vss/blatt4/blatt-4-lallinger_stortz_blatt4/proto"
)

type User struct{}

func (User) Req(context.Context, *proto.UserRequest, *proto.UserResponse) error {
	panic("implement me")
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("user"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	service.Init(
		// Add runtime action
		// We could actually do this above
		micro.Action(func(c *cli.Context) {

		}),
	)
	// Setup the server
	// Register handler
	proto.RegisterUserHandler(service.Server(), new(User))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
