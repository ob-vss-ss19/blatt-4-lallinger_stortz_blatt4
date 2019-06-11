package main

import (
	"context"
	"fmt"
	"github.com/micro/cli"
	proto "github.com/micro/examples/service/proto"
	"github.com/micro/go-micro"
)

/*

Example usage of top level service initialisation

*/

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}


func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("cinemaService"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),
	)

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	service.Init(
		// Add runtime action
		// We could actually do this above
		micro.Action(func(c *cli.Context) {

		}),
	)
	// Setup the server

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))


	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
