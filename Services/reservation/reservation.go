package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type reservationData struct {
	user    string
	showing int
	seats   int
	booked  bool
}

type Reservation struct {
	reservations map[int][]reservationData
}

func (me *Reservation) RequestReservation(ctx context.Context, req *proto.ReservationData, rsp *proto.ReservationData) error {
	return nil
}
func (me *Reservation) BookReservation(ctx context.Context, req *proto.ReservationData, rsp *proto.Response) error {
	return nil
}
func (me *Reservation) DeleteReservation(ctx context.Context, req *proto.ReservationData, rsp *proto.Response) error {
	return nil
}
func (me *Reservation) GetReservations(ctx context.Context, req *proto.ReservationRequest, rsp *proto.ReservationResponse) error {
	return nil
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
