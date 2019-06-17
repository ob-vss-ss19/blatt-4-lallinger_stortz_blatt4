package movie

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
	"math/big"
)

type Movie struct {
	// use map as set
	movies map[string]bool
}

func deleteShowings(movie string) {
	var client client.Client
	show := proto.NewShowingService("showing", client)

	rsp, err := show.GetShowings(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range rsp.Data {
		if v.Movie == movie {
			resp, err := show.DeleteShowing(context.TODO(), &proto.ShowingData{Id: v.Id})
			if err != nil || !resp.Success {
				fmt.Println(err)
				return
			}
		}
	}
}

func (me *Movie) AddMovie(ctx context.Context, req *proto.MovieData, rsp *proto.Response) error {
	if me.movies == nil {
		me.movies = make(map[string]bool)
	}
	if _, ok := me.movies[req.Title]; ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Movie %s already exists.", req.Title)
		return nil
	}

	me.movies[req.Title] = true
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Added %s to movies.", req.Title)
	return nil
}
func (me *Movie) DeleteMovie(ctx context.Context, req *proto.MovieData, rsp *proto.Response) error {
	if _, ok := me.movies[req.Title]; !ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Movie %s does not exist.", req.Title)
		return nil
	}

	deleteShowings(req.Title)
	delete(me.movies, req.Title)
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Deleted %s from movies.", req.Title)
	return nil
}
func (me *Movie) GetMovies(ctx context.Context, req *proto.Request, rsp *proto.MovieResponse) error {
	for k := range me.movies {
		rsp.Data = append(rsp.Data, &proto.MovieData{Title: k})
	}
	return nil
}

func StartMovieService(ctx context.Context, test bool) {
	var port int64
	port = 8093
	if test {
		reader := rand.Reader
		rsp, _ := rand.Int(reader, big.NewInt(1000))
		port = 1024 + 4 + rsp.Int64()
	}

	service := micro.NewService(
		micro.Name("movie"),
		micro.Address(fmt.Sprintf(":%v", port)),
		micro.Context(ctx),
	)

	if !test {
		service.Init()
	}
	proto.RegisterMovieHandler(service.Server(), new(Movie))

	fmt.Println("Starting movie service")
	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
