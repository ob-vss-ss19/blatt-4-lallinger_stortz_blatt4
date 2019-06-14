package reservation

import (
	"context"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/cinemahall"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/movie"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/showing"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
	"github.com/stretchr/testify/assert"
)

func TestReservation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go StartReservationService(ctx, true)
	time.Sleep(300 * time.Millisecond)
	go cinemahall.StartCinemaService(ctx, true)
	time.Sleep(300 * time.Millisecond)
	go showing.StartShowingService(ctx, true)
	time.Sleep(300 * time.Millisecond)
	go user.StartUserService(ctx, true)
	time.Sleep(300 * time.Millisecond)
	go movie.StartMovieService(ctx, true)
	time.Sleep(300 * time.Millisecond)

	var client client.Client
	cinema := proto.NewCinemaService("cinema", client)
	movie := proto.NewMovieService("movie", client)
	showing := proto.NewShowingService("showing", client)
	use := proto.NewUserService("user", client)
	reserve := proto.NewReservationService("reservation", client)

	// add movie
	req1 := &proto.MovieData{Title: "mov"}
	rsp, err := movie.AddMovie(context.TODO(), req1)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// add cinema
	reqm := &proto.CinemaData{Name: "testcinema", RowLength: 5, Rows: 5}
	rsp, err = cinema.AddCinema(context.TODO(), reqm)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// add showing
	reqs := &proto.ShowingData{Movie: "mov", Cinema: "testcinema"}
	rsp, err = showing.AddShowing(context.TODO(), reqs)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// add user
	requ := &proto.UserData{Name: "sepp"}
	rsp, err = use.CreateUser(context.TODO(), requ)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// not existing user
	reqr := &proto.ReservationData{User: "hans", Showing: 0, Seats: 10}
	rsp, err = reserve.RequestReservation(context.TODO(), reqr)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "User hans does not exist.")

	// not existing showing
	reqr = &proto.ReservationData{User: "sepp", Showing: 1, Seats: 10}
	rsp, err = reserve.RequestReservation(context.TODO(), reqr)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "Showing 1 does not exist.")

	// not enough seats
	reqr = &proto.ReservationData{User: "sepp", Showing: 0, Seats: 26}
	rsp, err = reserve.RequestReservation(context.TODO(), reqr)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "Not enough seats. Only 25 remaining, trying to reserve 26.")

	// request reservation 0
	reqr = &proto.ReservationData{User: "sepp", Showing: 0, Seats: 25}
	rsp, err = reserve.RequestReservation(context.TODO(), reqr)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// request reservation 1
	rsp, err = reserve.RequestReservation(context.TODO(), reqr)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// get reservations
	req := &proto.Request{}
	rspr, err := reserve.GetReservations(context.TODO(), req)
	assert.Nil(t, err)
	assert.Len(t, rspr.Data, 2)
	i := 0
	for _, v := range rspr.Data {
		if v.User == reqr.User && v.Showing == reqr.Showing && v.Seats == reqr.Seats && v.Booked == false {
			if v.ReservationID == 0 {
				i++
			}
			if v.ReservationID == 1 {
				i++
			}
		}
	}
	assert.Equal(t, i, 2)

	// reservation not existing
	reqrb := &proto.ReservationData{ReservationID: 10}
	rsp, err = reserve.BookReservation(context.TODO(), reqrb)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "Reservation ID 10 does not exist.")

	// book reservation 0
	reqrb = &proto.ReservationData{ReservationID: 0}
	rsp, err = reserve.BookReservation(context.TODO(), reqrb)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// double book reservation 0
	reqrb = &proto.ReservationData{ReservationID: 0}
	rsp, err = reserve.BookReservation(context.TODO(), reqrb)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "Reservation 0 already booked")

	// book reservation 1 => assert error
	reqrb = &proto.ReservationData{ReservationID: 1}
	rsp, err = reserve.BookReservation(context.TODO(), reqrb)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "Could not book reservation ID 1, not enough free seats. Try again later or delete reservation.")

	// request reservation 0 again => assert error
	reqr = &proto.ReservationData{User: "sepp", Showing: 0, Seats: 1}
	rsp, err = reserve.RequestReservation(context.TODO(), reqr)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "Not enough seats. Only 0 remaining, trying to reserve 1.")

	// get reservations
	rspr, err = reserve.GetReservations(context.TODO(), req)
	assert.Nil(t, err)
	assert.Len(t, rspr.Data, 2)
	i = 0
	for _, v := range rspr.Data {
		if v.User == "sepp" && v.Showing == 0 && v.Seats == 25 {
			if v.ReservationID == 0 && v.Booked {
				i++
			}
			if v.ReservationID == 1 && !v.Booked {
				i++
			}
		}
	}
	assert.Equal(t, i, 2)

	//delete reservation 1
	rsp, err = reserve.DeleteReservation(context.TODO(), reqrb)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	//delete reservation 1 again => error
	rsp, err = reserve.DeleteReservation(context.TODO(), reqrb)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)

	// get reservations
	rspr, err = reserve.GetReservations(context.TODO(), req)
	assert.Nil(t, err)
	assert.Len(t, rspr.Data, 1)
	assert.Equal(t, rspr.Data[0].ReservationID, int32(0))

	cancel()
	time.Sleep(1 * time.Second)
}
