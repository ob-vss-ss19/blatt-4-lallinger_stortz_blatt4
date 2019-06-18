package showing

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type showingData struct {
	movie  string
	cinema string
}

type Showing struct {
	showings map[int32]*showingData
	nextID   int32
}

func deleteReservations(showing int32) {
	var client client.Client
	reservation := proto.NewReservationService("reservation", client)

	rsp, err := reservation.GetReservations(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range rsp.Data {
		if v.Showing == showing {
			resp, err := reservation.DeleteReservation(context.TODO(), &proto.ReservationData{ReservationID: v.ReservationID})
			if err != nil || !resp.Success {
				fmt.Println(err)
				return
			}
		}
	}
}

func movieExists(title string) bool {
	var client client.Client
	mov := proto.NewMovieService("movie", client)

	// Call
	rsp, err := mov.GetMovies(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, v := range rsp.Data {
		if v.Title == title {
			return true
		}
	}
	return false
}

func cinemaExists(name string) bool {
	var client client.Client
	cine := proto.NewCinemaService("cinema", client)

	// Call
	rsp, err := cine.GetCinemas(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, v := range rsp.Data {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (me *Showing) AddShowing(ctx context.Context, req *proto.ShowingData, rsp *proto.Response) error {
	if me.showings == nil {
		me.showings = make(map[int32]*showingData)
		me.nextID = 0
	}
	if !movieExists(req.Movie) {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Movie %s does not exist.", req.Movie)
		return nil
	}
	if !cinemaExists(req.Cinema) {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Cinema %s does not exist.", req.Cinema)
		return nil
	}

	me.showings[me.nextID] = &showingData{cinema: req.Cinema, movie: req.Movie}
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Added showing %d for %s in %s.", me.nextID, req.Movie, req.Cinema)
	me.nextID++
	return nil
}
func (me *Showing) DeleteShowing(ctx context.Context, req *proto.ShowingData, rsp *proto.Response) error {
	if _, ok := me.showings[req.Id]; !ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Showing ID %d does not exist.", req.Id)
		return nil
	}

	deleteReservations(req.Id)
	delete(me.showings, req.Id)
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Deleted showing %d.", req.Id)
	return nil
}
func (me *Showing) GetShowings(ctx context.Context, req *proto.Request, rsp *proto.ShowingResponse) error {
	for k, v := range me.showings {
		rsp.Data = append(rsp.Data, &proto.ShowingData{Id: k, Cinema: v.cinema, Movie: v.movie})
	}
	return nil
}

func StartShowingService(ctx context.Context, test bool) {
	var port int64
	port = 8095
	//if test {
	//	reader := rand.Reader
	//	rsp, _ := rand.Int(reader, big.NewInt(1000))
	//	port = 1024 + 4 + rsp.Int64()
	//}

	service := micro.NewService(
		micro.Name("showing"),
		micro.Address(fmt.Sprintf(":%v", port)),
		micro.Context(ctx),
	)

	if !test {
		service.Init()
	}
	proto.RegisterShowingHandler(service.Server(), new(Showing))

	fmt.Println("Starting showing service")
	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
