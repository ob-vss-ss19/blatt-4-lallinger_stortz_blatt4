package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type Reservation struct{}

func (Reservation) Request(context.Context, *proto.ReservationRequest, *proto.ReservationResponse) error {
	panic("implement me")
}

func main() {
	service := micro.NewService(micro.Name("reservation"))
	service.Init()
	proto.RegisterReservationHandler(service.Server(), new(Reservation))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
