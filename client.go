package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
	"strconv"
)

var (
	help = flag.Bool("help", false, "print help")

	cine        proto.CinemaService
	movie       proto.MovieService
	reservation proto.ReservationService
	showing     proto.ShowingService
	user        proto.UserService
)

func main() {

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("client.exe SERVICE FUNCTION PARAMS")
		fmt.Println("SERVICE")
		fmt.Println(" cinema")
		fmt.Println("  FUNCTION")
		fmt.Println("  -add PARAMS: name. Example: client.exe cinema add hall1")
		fmt.Println("  -delete PARAMS: name. Example: client.exe cinema delete hall1")
		fmt.Println("  -get: Example: client.exe cinema get")
		fmt.Println(" movie")
		fmt.Println("  FUNCTION")
		fmt.Println("  -add PARAMS: title. Example: client.exe movie add shrek")
		fmt.Println("  -delete PARAMS: title. Example: client.exe movie delete shrek")
		fmt.Println("  -get: Example: client.exe movie get")
		fmt.Println(" reservation")
		fmt.Println("  FUNCTION")
		fmt.Println("  -request PARAMS: user showingID seats. Example: client.exe reservation request sepp 2 4")
		fmt.Println("   Requests a reservation.")
		fmt.Println("  -book PARAMS: reservationID. Example: client.exe reservation book 1")
		fmt.Println("   Books a reservation.")
		fmt.Println("  -delete PARAMS: reservationID. Example: client.exe reservation delete 1")
		fmt.Println("  -get: Example: client.exe reservation get")
		fmt.Println(" showing")
		fmt.Println("  FUNCTION")
		fmt.Println("  -add PARAMS: movie cinema. Example: client.exe showing add shrek hall1")
		fmt.Println("  -delete PARAMS: showingID. Example: client.exe showing delete 4")
		fmt.Println("  -get: Example: client.exe showing get")
		fmt.Println(" user")
		fmt.Println("  FUNCTION")
		fmt.Println("  -add PARAMS: name. Example: client.exe user add sepp")
		fmt.Println("  -delete PARAMS: name. Example: client.exe user delete sepp")
		fmt.Println("  -get: Example: client.exe user get")
		fmt.Println(" fill")
		fmt.Println("  -Fills services with some data. Example: client.exe fill")
		return
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	service := micro.NewService(micro.Name("client"))
	service.Init()

	switch flag.Arg(0) {
	case "cinema":
		cine = proto.NewCinemaService("cinema", service.Client())
		switch flag.Arg(1) {
		case "add":
			printResponse(cine.AddCinema(context.TODO(), &proto.CinemaData{Name: flag.Arg(2), Rows: toInt32(flag.Arg(3)), RowLength: toInt32(flag.Arg(4))}))
		case "delete":
			printResponse(cine.DeleteCinema(context.TODO(), &proto.CinemaData{Name: flag.Arg(2)}))
		case "get":
			printCinemas(cine.GetCinemas(context.TODO(), &proto.Request{}))
		}
	case "movie":
		movie = proto.NewMovieService("movie", service.Client())
		switch flag.Arg(1) {
		case "add":
			printResponse(movie.AddMovie(context.TODO(), &proto.MovieData{Title: flag.Arg(2)}))
		case "delete":
			printResponse(movie.DeleteMovie(context.TODO(), &proto.MovieData{Title: flag.Arg(2)}))
		case "get":
			printMovies(movie.GetMovies(context.TODO(), &proto.Request{}))
		}
	case "reservation":
		reservation = proto.NewReservationService("reservation", service.Client())
		switch flag.Arg(1) {
		case "request":
			printResponse(reservation.RequestReservation(context.TODO(), &proto.ReservationData{User: flag.Arg(2), Showing: toInt32(flag.Arg(3)), Seats: toInt32(flag.Arg(4))}))
		case "book":
			printResponse(reservation.BookReservation(context.TODO(), &proto.ReservationData{ReservationID: toInt32(flag.Arg(2))}))
		case "delete":
			printResponse(reservation.DeleteReservation(context.TODO(), &proto.ReservationData{ReservationID: toInt32(flag.Arg(2))}))
		case "get":
			printReservations(reservation.GetReservations(context.TODO(), &proto.Request{}))
		}
	case "showing":
		showing = proto.NewShowingService("showing", service.Client())
		switch flag.Arg(1) {
		case "add":
			printResponse(showing.AddShowing(context.TODO(), &proto.ShowingData{Movie: flag.Arg(2), Cinema: flag.Arg(3)}))
		case "delete":
			printResponse(showing.DeleteShowing(context.TODO(), &proto.ShowingData{Id: toInt32(flag.Arg(2))}))
		case "get":
			printShowings(showing.GetShowings(context.TODO(), &proto.Request{}))
		}
	case "user":
		user = proto.NewUserService("user", service.Client())
		switch flag.Arg(1) {
		case "add":
			printResponse(user.CreateUser(context.TODO(), &proto.UserData{Name: flag.Arg(2)}))
		case "delete":
			printResponse(user.DeleteUser(context.TODO(), &proto.UserData{Name: flag.Arg(2)}))
		case "get":
			printUsers(user.GetUsers(context.TODO(), &proto.Request{}))
		}
	case "fill":
		cine = proto.NewCinemaService("cinema", service.Client())
		movie = proto.NewMovieService("movie", service.Client())
		reservation = proto.NewReservationService("reservation", service.Client())
		showing = proto.NewShowingService("showing", service.Client())
		user = proto.NewUserService("user", service.Client())

		printResponse(cine.AddCinema(context.TODO(), &proto.CinemaData{Name: "saal1", Rows: 10, RowLength: 10}))
		printResponse(cine.AddCinema(context.TODO(), &proto.CinemaData{Name: "saal2", Rows: 5, RowLength: 6}))
		printResponse(movie.AddMovie(context.TODO(), &proto.MovieData{Title: "terminator"}))
		printResponse(movie.AddMovie(context.TODO(), &proto.MovieData{Title: "shrek"}))
		printResponse(movie.AddMovie(context.TODO(), &proto.MovieData{Title: "avengers"}))
		printResponse(movie.AddMovie(context.TODO(), &proto.MovieData{Title: "hulk"}))
		printResponse(showing.AddShowing(context.TODO(), &proto.ShowingData{Movie: "terminator", Cinema: "saal1"}))
		printResponse(showing.AddShowing(context.TODO(), &proto.ShowingData{Movie: "shrek", Cinema: "saal2"}))
		printResponse(showing.AddShowing(context.TODO(), &proto.ShowingData{Movie: "hulk", Cinema: "saal2"}))
		printResponse(showing.AddShowing(context.TODO(), &proto.ShowingData{Movie: "avengers", Cinema: "saal1"}))
		printResponse(user.CreateUser(context.TODO(), &proto.UserData{Name: "sepp"}))
		printResponse(user.CreateUser(context.TODO(), &proto.UserData{Name: "hans"}))
		printResponse(user.CreateUser(context.TODO(), &proto.UserData{Name: "franz"}))
		printResponse(user.CreateUser(context.TODO(), &proto.UserData{Name: "xaver"}))
		printResponse(reservation.RequestReservation(context.TODO(), &proto.ReservationData{User: "sepp", Showing: 0, Seats: 10}))
		printResponse(reservation.RequestReservation(context.TODO(), &proto.ReservationData{User: "franz", Showing: 1, Seats: 15}))
		printResponse(reservation.RequestReservation(context.TODO(), &proto.ReservationData{User: "hans", Showing: 2, Seats: 20}))
		printResponse(reservation.RequestReservation(context.TODO(), &proto.ReservationData{User: "xaver", Showing: 3, Seats: 20}))
		printResponse(reservation.BookReservation(context.TODO(), &proto.ReservationData{ReservationID: 0}))
		printResponse(reservation.BookReservation(context.TODO(), &proto.ReservationData{ReservationID: 3}))
	default:
		flag.Usage()
		return
	}
}

func printResponse(rsp *proto.Response, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print response
	if rsp.Success {
		fmt.Print("Success: ")
	} else {
		fmt.Print("Error: ")
	}
	fmt.Printf("%s\n", rsp.Message)
}

func printCinemas(rsp *proto.CinemaResponse, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%-10s %-5s\n", "Name", "Seats")
	for _, v := range rsp.Data {
		fmt.Printf("%-10s %-5d\n", v.Name, v.Rows*v.RowLength)
	}
}

func printMovies(rsp *proto.MovieResponse, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Title\n")
	for _, v := range rsp.Data {
		fmt.Printf("%s\n", v.Title)
	}
}

func printReservations(rsp *proto.ReservationResponse, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%-3s %-4s %-10s %-4s %-5s\n", "Id", "Show", "User", "Seats", "Booked")
	for _, v := range rsp.Data {
		fmt.Printf("%-3d %-4d %-10s %-4d %-5t\n", v.ReservationID, v.Showing, v.User, v.Seats, v.Booked)
	}
}

func printShowings(rsp *proto.ShowingResponse, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%-3s %-10s %-10s\n", "Id", "Cinema", "Movie")
	for _, v := range rsp.Data {
		fmt.Printf("%-3d %-10s %-10s\n", v.Id, v.Cinema, v.Movie)
	}
}

func printUsers(rsp *proto.UserResponse, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Name\n")
	for _, v := range rsp.Users {
		fmt.Printf("%s\n", v.Name)
	}
}

func toInt32(text string) int32 {
	tmp, _ := strconv.Atoi(text)
	return int32(tmp)
}
