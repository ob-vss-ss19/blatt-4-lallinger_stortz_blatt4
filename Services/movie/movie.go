package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type Movie struct{}

func (Movie) Request(context.Context, *proto.MovieRequest, *proto.MovieResponse) error {
	panic("implement me")
}

func main() {
	service := micro.NewService(micro.Name("movie"))
	service.Init()
	proto.RegisterMovieHandler(service.Server(), new(Movie))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
