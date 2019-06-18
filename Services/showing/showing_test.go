package showing

import (
	"context"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/cinemahall"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/movie"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/reservation"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/user"
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
	go StartShowingService(ctx, true)
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

	// add second showing
	rsp, err = showing.AddShowing(context.TODO(), reqs)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// not existing movie
	reqs = &proto.ShowingData{Movie: "nomov", Cinema: "testcinema"}
	rsp, err = showing.AddShowing(context.TODO(), reqs)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "Movie nomov does not exist.")

	// not existing movie
	reqs = &proto.ShowingData{Movie: "mov", Cinema: "notestcinema"}
	rsp, err = showing.AddShowing(context.TODO(), reqs)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)
	assert.Equal(t, rsp.Message, "Cinema notestcinema does not exist.")

	// get showings
	req := &proto.Request{}
	rspr, err := showing.GetShowings(context.TODO(), req)
	assert.Nil(t, err)
	assert.Len(t, rspr.Data, 2)
	i := 0
	for _, v := range rspr.Data {
		if v.Movie == "mov" && v.Cinema == "testcinema" {
			if v.Id == 0 {
				i++
			}
			if v.Id == 1 {
				i++
			}
		}
	}
	assert.Equal(t, i, 2)

	// add user
	requ := &proto.UserData{Name: "sepp"}
	rsp, err = use.CreateUser(context.TODO(), requ)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	// add reservation
	reqr := &proto.ReservationData{User: "sepp", Showing: 0, Seats: 25}
	rsp, err = reserve.RequestReservation(context.TODO(), reqr)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	//delete showing 0
	reqsd := &proto.ShowingData{Id: 0}
	rsp, err = showing.DeleteShowing(context.TODO(), reqsd)
	assert.Nil(t, err)
	assert.True(t, rsp.Success)

	//delete showing 0 again => assert error
	rsp, err = showing.DeleteShowing(context.TODO(), reqsd)
	assert.Nil(t, err)
	assert.False(t, rsp.Success)

	// get reservations
	rsprr, err := reserve.GetReservations(context.TODO(), req)
	assert.Nil(t, err)
	assert.Len(t, rsprr.Data, 0)

	// get showings
	rspr, err = showing.GetShowings(context.TODO(), req)
	assert.Nil(t, err)
	assert.Len(t, rspr.Data, 1)
	assert.Equal(t, rspr.Data[0].Id, int32(1))

	cancel()
	time.Sleep(1 * time.Second)
}
