package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type Movie struct {
	// use map as set
	movies map[string]bool
}

func deleteShowings(movie string) {
	service := micro.NewService(micro.Name("movieRequest"))
	service.Init()
	show := proto.NewShowingService("showing", service.Client())

	rsp, err := show.GetShowings(context.TODO(), &proto.ShowingRequest{})
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
func (me *Movie) GetMovies(ctx context.Context, req *proto.MovieRequest, rsp *proto.MovieResponse) error {
	for k := range me.movies {
		rsp.Data = append(rsp.Data, &proto.MovieData{Title: k})
	}
	return nil
}

func main() {
	service := micro.NewService(micro.Name("movie"))
	service.Init()
	proto.RegisterMovieHandler(service.Server(), new(Movie))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
