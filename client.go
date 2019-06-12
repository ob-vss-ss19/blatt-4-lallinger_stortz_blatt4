package main

import (
	"fmt"
	"os"

	"context"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

/*

Example usage of top level service initialisation

*/

//type Greeter struct{}

type Cinema struct{}

/*func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}*/

// Setup and the client
func runClient(service micro.Service) {
	// Create new greeter client
	cine := proto.NewCinemaService("cinema", service.Client())

	// Call the greeter
	rsp, err := cine.Req(context.TODO(), &proto.CinemaRequest{Column: "Name", Value: "testKino"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response
	for _, cd := range rsp.Data {
		fmt.Printf("Name: %s, Rows: %d, RowLength:%d\n", cd.Name, cd.Rows, cd.RowLength)
	}

}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("cinema"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),

		// Setup some flags. Specify --run_client to run the client

		// Add runtime flags
		// We could do this below too
		micro.Flags(cli.BoolFlag{
			Name:  "run_client",
			Usage: "Launch the client",
		}),
	)

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	service.Init(
		// Add runtime action
		// We could actually do this above
		micro.Action(func(c *cli.Context) {
			runClient(service)
			os.Exit(0)
		}),
	)
}
