package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

//Wsl am besten ueber reflection field namen nutzen um offen fuer alle Anfragen zu bleiben -> bei nem Prototyp nicht unbedingt notwendig ?
type cinemaData struct {
	Name      string
	Rows      int32
	RowLength int32
}

var cinemaDataList []cinemaData

type Cinema struct{}

func (Cinema) Request(ctx context.Context, req *proto.CinemaRequest, rsp *proto.CinemaResponse) error {
	for _, cd := range cinemaDataList {
		if req.Value == cd.Name {
			rsp.Data = append(rsp.Data, &proto.CinemaData{Name: cd.Name, RowLength: cd.RowLength, Rows: cd.Rows})
		}
	}
	return nil
}

func main() {

	cinemaDataList = append(cinemaDataList, cinemaData{"testKino", 17, 17})
	cinemaDataList = append(cinemaDataList, cinemaData{"testKino2", 42, 42})
	cinemaDataList = append(cinemaDataList, cinemaData{"mettthaeser", 23, 23})

	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("cinema"))
	service.Init()
	proto.RegisterCinemaHandler(service.Server(), new(Cinema))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
