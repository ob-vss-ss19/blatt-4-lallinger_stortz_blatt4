package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4/proto"
)

type User struct {
	// use map as set
	users map[string]bool
}

func (me *User) CreateUser(ctx context.Context, req *proto.UserData, rsp *proto.Response) error {
	if me.users==nil {
		me.users=make(map[string]bool)
	}
	if _, ok := me.users[req.Name]; ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("User %s already exists.", req.Name)
		return nil
	}

	me.users[req.Name]=true
	rsp.Success=true
	rsp.Message=fmt.Sprintf("Created User %s.",req.Name)
	return nil
}
func (me *User) DeleteUser(ctx context.Context, req *proto.UserData, rsp *proto.Response) error {
	if _, ok := me.users[req.Name]; !ok {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("User %s does not exist.", req.Name)
		return nil
	}

	delete(me.users, req.Name)
	rsp.Success=true
	rsp.Message=fmt.Sprintf("Deleted User %s.",req.Name)
	return nil
}
func (me *User) GetUsers(ctx context.Context, req *proto.UserRequest, rsp *proto.UserResponse) error {
	for k := range me.users {
		rsp.Users= append(rsp.Users, &proto.UserData{Name:k})
	}
	return nil
}

func main() {
	service := micro.NewService(micro.Name("user"))
	service.Init()
	proto.RegisterUserHandler(service.Server(), new(User))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
