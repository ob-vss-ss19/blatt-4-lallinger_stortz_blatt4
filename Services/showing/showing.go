package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type showingData struct {
	movie  string
	cinema string
}

type Showing struct {
	showings map[int][]showingData
}

func (me *Showing) AddShowing(ctx context.Context, req *proto.ShowingData, rsp *proto.Response) error {
	return nil
}
func (me *Showing) DeleteShowing(ctx context.Context, req *proto.ShowingData, rsp *proto.Response) error {
	return nil
}
func (me *Showing) GetShowings(ctx context.Context, req *proto.ShowingRequest, rsp *proto.ShowingResponse) error {
	return nil
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
