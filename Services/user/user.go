package user

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type User struct {
	// use map as set
	users map[string]bool
}

func deleteReservations(user string) {
	var client client.Client
	reservation := proto.NewReservationService("reservation", client)

	rsp, err := reservation.GetReservations(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range rsp.Data {
		if v.User == user {
			resp, err := reservation.DeleteReservation(context.TODO(), &proto.ReservationData{ReservationID: v.ReservationID})
			if err != nil || !resp.Success {
				fmt.Println(err)
				return
			}
		}
	}
}

func (me *User) CreateUser(ctx context.Context, req *proto.UserData, rsp *proto.Response) error {
	if me.users == nil {
		me.users = make(map[string]bool)
	}
	if _, ok := me.users[req.Name]; ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("User %s already exists.", req.Name)
		return nil
	}

	me.users[req.Name] = true
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Created User %s.", req.Name)
	return nil
}
func (me *User) DeleteUser(ctx context.Context, req *proto.UserData, rsp *proto.Response) error {
	if _, ok := me.users[req.Name]; !ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("User %s does not exist.", req.Name)
		return nil
	}

	deleteReservations(req.Name)
	delete(me.users, req.Name)
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Deleted User %s.", req.Name)
	return nil
}
func (me *User) GetUsers(ctx context.Context, req *proto.Request, rsp *proto.UserResponse) error {
	for k := range me.users {
		rsp.Users = append(rsp.Users, &proto.UserData{Name: k})
	}
	return nil
}

func StartUserService(ctx context.Context, test bool) {
	var port int64
	port = 8096
	//if test {
	//	reader := rand.Reader
	//	rsp, _ := rand.Int(reader, big.NewInt(1000))
	//	port = 1024 + 4 + rsp.Int64()
	//}

	service := micro.NewService(
		micro.Name("user"),
		micro.Address(fmt.Sprintf(":%v", port)),
		micro.Context(ctx),
	)

	if !test {
		service.Init()
	}
	proto.RegisterUserHandler(service.Server(), new(User))

	fmt.Println("Starting user service")
	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
