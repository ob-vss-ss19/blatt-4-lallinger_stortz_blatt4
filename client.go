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
	rsp, err := cine.AddCinema(context.TODO(), &proto.CinemaData{Name: "Kino2", Rows: 11, RowLength: 13})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response
	if rsp.Success {
		fmt.Println("Success!")
	} else {
		fmt.Println("Error")
	}
	fmt.Printf("%s\n", rsp.Message)

	resp, err := cine.GetCinemas(context.TODO(), &proto.CinemaRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range resp.Data {
		fmt.Printf("Cinema: %s. Rows: %d. Rowlength: %d\n", v.Name, v.Rows, v.RowLength)
	}
}
