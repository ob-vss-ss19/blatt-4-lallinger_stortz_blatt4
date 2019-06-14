package main

import (
	"context"
	cinema "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/cinemahall"
	mov "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/movie"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/reservation"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/showing"
	"github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/Services/user"
	"time"
)

func main() {
	go reservation.StartReservationService(context.TODO(),false)
	time.Sleep(300*time.Millisecond)
	go cinema.StartCinemaService(context.TODO(),false)
	time.Sleep(300*time.Millisecond)
	go showing.StartShowingService(context.TODO(),false)
	time.Sleep(300*time.Millisecond)
	go user.StartUserService(context.TODO(),false)
	time.Sleep(300*time.Millisecond)
	mov.StartMovieService(context.TODO(),false)
}
