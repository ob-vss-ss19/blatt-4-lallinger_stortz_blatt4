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
	movie:= proto.NewMovieService("movie",service.Client())
	reservation:= proto.NewReservationService("reservation",service.Client())
	showing:=proto.NewShowingService("showing",service.Client())
	user:=proto.NewUserService("user",service.Client())

	// Call
	rsp, err := cine.AddCinema(context.TODO(), &proto.CinemaData{Name: "Kino2", Rows: 5, RowLength: 5})
	if err != nil {
		fmt.Println(err)
		return
	}
	printResponse(rsp)

	rsp, err = movie.AddMovie(context.TODO(),&proto.MovieData{Title:"firstmovie"})
	if err != nil {
		fmt.Println(err)
		return
	}
	printResponse(rsp)

	rsp, err = showing.AddShowing(context.TODO(),&proto.ShowingData{Movie:"firstmovie",Cinema:"Kino2"})
	if err != nil {
		fmt.Println(err)
		return
	}
	printResponse(rsp)

	rsp, err = user.CreateUser(context.TODO(),&proto.UserData{Name:"sepp"})
	if err != nil {
		fmt.Println(err)
		return
	}
	printResponse(rsp)


	rsp, err = reservation.RequestReservation(context.TODO(),&proto.ReservationData{User:"sepp",Showing:0,Seats:26})
	if err != nil {
		fmt.Println(err)
		return
	}
	printResponse(rsp)

	rsp, err = reservation.BookReservation(context.TODO(),&proto.ReservationData{ReservationID:0})
	if err != nil {
		fmt.Println(err)
		return
	}
	printResponse(rsp)


}

func printResponse(rsp *proto.Response){
	// Print response
	if rsp.Success {
		fmt.Println("Success!")
	} else {
		fmt.Println("Error")
	}
	fmt.Printf("%s\n", rsp.Message)
}
