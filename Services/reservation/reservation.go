package reservation

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type reservationData struct {
	showing int32
	seats   int32
	booked  bool
	user    string
}

type Reservation struct {
	reservations map[int32]*reservationData
	nextID       int32
}

func getFreeSeats(reservations map[int32]*reservationData, showing int32) int32 {

	var client client.Client
	show := proto.NewShowingService("showing", client)

	// Call
	rsp, err := show.GetShowings(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return -1
	}

	cinema := ""
	for _, v := range rsp.Data {
		if v.Id == showing {
			cinema = v.Cinema
			break
		}
	}

	cine := proto.NewCinemaService("cinema", client)
	resp, err := cine.GetCinemas(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return -1
	}
	var seats int32
	for _, v := range resp.Data {
		if v.Name == cinema {
			seats = v.RowLength * v.Rows
			break
		}
	}

	for _, v := range reservations {
		if v.showing == showing && v.booked {
			seats -= v.seats
		}
	}
	return seats
}

func showingExists(showing int32) bool {
	var client client.Client
	show := proto.NewShowingService("showing", client)

	// Call
	rsp, err := show.GetShowings(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, v := range rsp.Data {
		if v.Id == showing {
			return true
		}
	}
	return false
}

func userExists(name string) bool {
	var client client.Client
	user := proto.NewUserService("user", client)

	// Call
	rsp, err := user.GetUsers(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, v := range rsp.Users {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (me *Reservation) RequestReservation(ctx context.Context, req *proto.ReservationData, rsp *proto.Response) error {
	if me.reservations == nil {
		me.reservations = make(map[int32]*reservationData)
		me.nextID = 0
	}
	if !showingExists(req.Showing) {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Showing %d does not exist.", req.Showing)
		return nil
	}
	if !userExists(req.User) {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("User %s does not exist.", req.User)
		return nil
	}
	if tmp := getFreeSeats(me.reservations, req.Showing); tmp < req.Seats {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Not enough seats. Only %d remaining, trying to reserve %d.", tmp, req.Seats)
		return nil
	}

	me.reservations[me.nextID] = &reservationData{seats: req.Seats, showing: req.Showing, booked: false, user: req.User}
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Requested %d seats for showing %d. Reservation ID %d.", req.Seats, req.Showing, me.nextID)
	me.nextID++
	return nil
}
func (me *Reservation) BookReservation(ctx context.Context, req *proto.ReservationData, rsp *proto.Response) error {
	if _, ok := me.reservations[req.ReservationID]; !ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Reservation ID %d does not exist.", req.ReservationID)
		return nil
	}
	if me.reservations[req.ReservationID].booked {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Reservation %d already booked", req.ReservationID)
		return nil
	}

	if getFreeSeats(me.reservations, me.reservations[req.ReservationID].showing) < me.reservations[req.ReservationID].seats {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Could not book reservation ID %d, not enough free seats. Try again later or delete reservation.", req.ReservationID)
		return nil
	}

	me.reservations[req.ReservationID].booked = true
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Booked %d seats for showing %d.", me.reservations[req.ReservationID].seats, me.reservations[req.ReservationID].showing)

	return nil
}
func (me *Reservation) DeleteReservation(ctx context.Context, req *proto.ReservationData, rsp *proto.Response) error {
	if _, ok := me.reservations[req.ReservationID]; !ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Reservation ID %d does not exist.", req.ReservationID)
		return nil
	}

	delete(me.reservations, req.ReservationID)
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Deleted reservation %d.", req.ReservationID)

	return nil
}
func (me *Reservation) GetReservations(ctx context.Context, req *proto.Request, rsp *proto.ReservationResponse) error {
	for k, v := range me.reservations {
		rsp.Data = append(rsp.Data, &proto.ReservationData{ReservationID: k, Seats: v.seats, Showing: v.showing, Booked: v.booked, User: v.user})
	}
	return nil
}

func StartReservationService(ctx context.Context, test bool) {
	var port int64
	port = 8094
	//if test {
	//	reader := rand.Reader
	//	rsp, _ := rand.Int(reader, big.NewInt(1000))
	//	port = 1024 + 4 + rsp.Int64()
	//}

	service := micro.NewService(
		micro.Name("reservation"),
		micro.Address(fmt.Sprintf(":%v", port)),
		micro.Context(ctx),
	)

	if !test {
		service.Init()
	}
	proto.RegisterReservationHandler(service.Server(), new(Reservation))

	fmt.Println("Starting reservation service")
	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
