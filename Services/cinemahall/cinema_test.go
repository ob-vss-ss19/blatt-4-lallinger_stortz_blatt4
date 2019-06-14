package cinemahall

import (
	"context"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/reservation"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/showing"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/movie"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
	"github.com/stretchr/testify/assert"
)

func TestCinema(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go reservation.StartReservationService(ctx,true)
	time.Sleep(300*time.Millisecond)
	go StartCinemaService(ctx,true)
	time.Sleep(300*time.Millisecond)
	go showing.StartShowingService(ctx,true)
	time.Sleep(300*time.Millisecond)
	go user.StartUserService(ctx,true)
	time.Sleep(300*time.Millisecond)
	go movie.StartMovieService(ctx,true)
	time.Sleep(300*time.Millisecond)

	var client client.Client
	cinema := proto.NewCinemaService("cinema",client)
	movie := proto.NewMovieService("movie",client)
	showing := proto.NewShowingService("showing",client)

	// add cinema 1
	req1:= &proto.CinemaData{Name:"testcinemaone",RowLength:4,Rows:10}
	rsp,err:= cinema.AddCinema(context.TODO(),req1)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// add cinema 1 again => assert error
	rsp,err= cinema.AddCinema(context.TODO(),req1)
	assert.Nil(t, err)
	assert.False(t,rsp.Success)

	// add cinema 2
	req2 := &proto.CinemaData{Name:"testcinematwo",RowLength:6,Rows:6}
	rsp,err = cinema.AddCinema(context.TODO(),req2)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// check if 2 cinemas are added
	req:=&proto.Request{}
	resp,err := cinema.GetCinemas(context.TODO(),req)
	assert.Nil(t, err)
	assert.Len(t,resp.Data,2)
	i:=0
	for _,v := range resp.Data {
		if v.Name==req1.Name&&v.Rows==req1.Rows&&v.RowLength==req1.RowLength {
			i++
		}
		if v.Name==req2.Name&&v.Rows==req2.Rows&&v.RowLength==req2.RowLength {
			i++
		}
	}
	assert.Equal(t,i,2)

	// add movie
	reqm:=&proto.MovieData{Title:"testmovie"}
	rsp,err= movie.AddMovie(context.TODO(),reqm)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// add showing
	reqs:=&proto.ShowingData{Movie:"testmovie",Cinema:"testcinemaone"}
	rsp,err= showing.AddShowing(context.TODO(),reqs)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// delete cinema 1
	reqd:= &proto.CinemaData{Name:"testcinemaone"}
	rsp,err= cinema.DeleteCinema(context.TODO(),reqd)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// delete cinema 1 again => assert error
	rsp,err= cinema.DeleteCinema(context.TODO(),reqd)
	assert.Nil(t, err)
	assert.False(t,rsp.Success)

	// get cinemas expect one
	resp,err = cinema.GetCinemas(context.TODO(),req)
	assert.Nil(t, err)
	assert.Len(t,resp.Data,1)

	// get showings expect none
	rspsg,err:=showing.GetShowings(context.TODO(),req)
	assert.Nil(t,err)
	assert.Len(t,rspsg.Data,0)

	cancel()
	time.Sleep(1 * time.Second)
}
