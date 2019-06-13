package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

//Wsl am besten ueber reflection field namen nutzen um offen fuer alle Anfragen zu bleiben -> bei nem Prototyp nicht unbedingt notwendig ?
type cinemaData struct {
	Rows      int32
	RowLength int32
}

type Cinema struct {
	cinemas map[string]cinemaData
}

func (me *Cinema) AddCinema(ctx context.Context, req *proto.CinemaData, rsp *proto.Response) error {
	if me.cinemas == nil {
		me.cinemas = make(map[string]cinemaData)
	}
	if _, ok := me.cinemas[req.Name]; ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Cinema %s already exists.", req.Name)
		return nil
	}

	me.cinemas[req.Name] = cinemaData{RowLength: req.RowLength, Rows: req.Rows}
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

	delete(me.cinemas, req.Name)
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Deleted cinema %s", req.Name)
	return nil
}
func (me *Cinema) GetCinemas(ctx context.Context, req *proto.CinemaRequest, rsp *proto.CinemaResponse) error {
	for k, v := range me.cinemas {
		rsp.Data = append(rsp.Data, &proto.CinemaData{Name: k, RowLength: v.RowLength, Rows: v.Rows})
	}
	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema"))
	service.Init()
	proto.RegisterCinemaHandler(service.Server(), new(Cinema))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
