package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type Showing struct{}

func (Showing) Request(context.Context, *proto.ShowingRequest, *proto.ShowingResponse) error {
	panic("implement me")
}

func main() {
	service := micro.NewService(micro.Name("showing"))
	service.Init()
	proto.RegisterShowingHandler(service.Server(), new(Showing))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
