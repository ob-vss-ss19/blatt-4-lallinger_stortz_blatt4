package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

func main() {
	service := micro.NewService(micro.Name("client"))
	service.Init()
	cine := proto.NewCinemaService("cinema", service.Client())

	// Call
	rsp, err := cine.Request(context.TODO(), &proto.CinemaRequest{Column: "Name", Value: "testKino"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response
	for _, cd := range rsp.Data {
		fmt.Printf("Name: %s, Rows: %d, RowLength:%d\n", cd.Name, cd.Rows, cd.RowLength)
	}
}
