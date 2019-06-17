package cinemahall

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
	"math/big"
)

type cinemaData struct {
	Rows      int32
	RowLength int32
}

type Cinema struct {
	cinemas map[string]*cinemaData
}

func deleteShowings(cinema string) {

	var client client.Client
	show := proto.NewShowingService("showing", client)

	rsp, err := show.GetShowings(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range rsp.Data {
		if v.Cinema == cinema {
			resp, err := show.DeleteShowing(context.TODO(), &proto.ShowingData{Id: v.Id})
			if err != nil || !resp.Success {
				fmt.Println(err)
				return
			}
		}
	}
}

func (me *Cinema) AddCinema(ctx context.Context, req *proto.CinemaData, rsp *proto.Response) error {
	if me.cinemas == nil {
		me.cinemas = make(map[string]*cinemaData)
	}
	if _, ok := me.cinemas[req.Name]; ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Cinema %s already exists.", req.Name)
		return nil
	}

	me.cinemas[req.Name] = &cinemaData{RowLength: req.RowLength, Rows: req.Rows}
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Added %s to cinemas.", req.Name)
	return nil
}
func (me *Cinema) DeleteCinema(ctx context.Context, req *proto.CinemaData, rsp *proto.Response) error {
	if _, ok := me.cinemas[req.Name]; !ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Cinema %s does not exist.", req.Name)
		return nil
	}

	deleteShowings(req.Name)
	delete(me.cinemas, req.Name)
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Deleted cinema %s", req.Name)
	return nil
}
func (me *Cinema) GetCinemas(ctx context.Context, req *proto.Request, rsp *proto.CinemaResponse) error {
	for k, v := range me.cinemas {
		rsp.Data = append(rsp.Data, &proto.CinemaData{Name: k, RowLength: v.RowLength, Rows: v.Rows})
	}
	return nil
}

func StartCinemaService(ctx context.Context, test bool) {
	var port int64
	port = 8092
	if test {
		reader := rand.Reader
		rsp, _ := rand.Int(reader, big.NewInt(1000))
		port = 1024 + 4 + rsp.Int64()
	}

	service := micro.NewService(
		micro.Name("cinema"),
		micro.Address(fmt.Sprintf(":%v", port)),
		micro.Context(ctx),
	)

	if !test {
		service.Init()
	}
	proto.RegisterCinemaHandler(service.Server(), new(Cinema))

	fmt.Println("Starting cinema service")
	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
