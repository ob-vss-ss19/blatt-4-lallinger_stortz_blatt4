package user

import (
	"context"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/cinemahall"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/movie"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/reservation"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/showing"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
	"github.com/stretchr/testify/assert"
)

func TestShowing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go reservation.StartReservationService(ctx, true)
	time.Sleep(300 * time.Millisecond)
	go cinemahall.StartCinemaService(ctx, true)
	time.Sleep(300 * time.Millisecond)
	go showing.StartShowingService(ctx, true)
	time.Sleep(300 * time.Millisecond)
	go StartUserService(ctx, true)
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

	// add user 1
	requ := &proto.UserData{Name: "sepp"}
	rsp, err = use.CreateUser(context.TODO(), requ)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// add user 1 again => assert error
	rsp, err = use.CreateUser(context.TODO(), requ)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)

	// add user 2
	requ = &proto.UserData{Name: "hans"}
	rsp, err = use.CreateUser(context.TODO(), requ)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// get users
	req := &proto.Request{}
	rspr, err := use.GetUsers(context.TODO(), req)
	assert.Nil(t, err)
	assert.Len(t, rspr.Users, 2)
	i := 0
	for _, v := range rspr.Users {
		if v.Name == "sepp" {
			i++
		}
		if v.Name == "hans" {
			i++
		}
	}
	assert.Equal(t, i, 2)

	// add reservation
	reqr := &proto.ReservationData{User: "sepp", Showing: 0, Seats: 25}
	rsp, err = reserve.RequestReservation(context.TODO(), reqr)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	//delete user 0
	reqsd := &proto.UserData{Name: "sepp"}
	rsp, err = use.DeleteUser(context.TODO(), reqsd)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	//delete user 0 again => assert error
	rsp, err = use.DeleteUser(context.TODO(), reqsd)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)

	// get reservations
	rsprr, err := reserve.GetReservations(context.TODO(), req)
	assert.Nil(t, err)
	assert.Len(t, rsprr.Data, 0)

	cancel()
	time.Sleep(1 * time.Second)
}
