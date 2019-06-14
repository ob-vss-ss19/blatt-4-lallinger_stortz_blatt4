package movie

import (
	"context"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/cinemahall"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/reservation"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/showing"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
	"github.com/stretchr/testify/assert"
)

func TestMovie(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go reservation.StartReservationService(ctx,true)
	time.Sleep(300*time.Millisecond)
	go cinemahall.StartCinemaService(ctx,true)
	time.Sleep(300*time.Millisecond)
	go showing.StartShowingService(ctx,true)
	time.Sleep(300*time.Millisecond)
	go user.StartUserService(ctx,true)
	time.Sleep(300*time.Millisecond)
	go StartMovieService(ctx,true)
	time.Sleep(300*time.Millisecond)

	var client client.Client
	cinema := proto.NewCinemaService("cinema",client)
	movie := proto.NewMovieService("movie",client)
	showing := proto.NewShowingService("showing",client)

	// add movie 1
	req1:= &proto.MovieData{Title:"mov1"}
	rsp,err:= movie.AddMovie(context.TODO(),req1)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// add movie 1 again => assert error
	rsp,err = movie.AddMovie(context.TODO(),req1)
	assert.Nil(t, err)
	assert.False(t,rsp.Success)

	// add movie 2
	req2:= &proto.MovieData{Title:"mov2"}
	rsp,err= movie.AddMovie(context.TODO(),req2)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// check if 2 movies are added
	req:=&proto.Request{}
	resp,err := movie.GetMovies(context.TODO(),req)
	assert.Nil(t, err)
	assert.Len(t,resp.Data,2)
	i:=0
	for _,v := range resp.Data {
		if v.Title==req1.Title {
			i++
		}
		if v.Title==req2.Title{
			i++
		}
	}
	assert.Equal(t,i,2)

	// add cinema
	reqm:=&proto.CinemaData{Name:"testcinema",RowLength:4,Rows:5}
	rsp,err= cinema.AddCinema(context.TODO(),reqm)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// add showing
	reqs:=&proto.ShowingData{Movie:"mov1",Cinema:"testcinema"}
	rsp,err= showing.AddShowing(context.TODO(),reqs)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// delete movie 1
	reqd:= &proto.MovieData{Title:"mov1"}
	rsp,err= movie.DeleteMovie(context.TODO(),reqd)
	assert.Nil(t, err)
	assert.True(t,rsp.Success)

	// delete movie  1 again => assert error
	rsp,err= movie.DeleteMovie(context.TODO(),reqd)
	assert.Nil(t, err)
	assert.False(t,rsp.Success)

	// get movies expect one
	resp,err = movie.GetMovies(context.TODO(),req)
	assert.Nil(t, err)
	assert.Len(t,resp.Data,1)
	assert.Equal(t,resp.Data[0].Title,req2.Title)

	// get showings expect none
	rspsg,err:=showing.GetShowings(context.TODO(),req)
	assert.Nil(t,err)
	assert.Len(t,rspsg.Data,0)

	cancel()
	time.Sleep(1 * time.Second)
}
