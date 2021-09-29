package client

import (
	"context"
	"errors"

	"gateway/user/userdb"
	"gateway/user/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Age      int32  `json:"age"`
	Greeting string `json:"greeting"`
}

type UserClient struct {
}

var (
	userGrpcServiceClient userpb.UserServiceClient
	userGrpcService       = userdb.GetEnv("USER_GRPC_SERVICE")
)

func prepareUserGrpcClient(c *context.Context) error {

	conn, err := grpc.DialContext(*c, userGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...)

	if err != nil {
		userGrpcServiceClient = nil
		return errors.New("connection to user gRPC service failed")
	}

	if userGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	userGrpcServiceClient = userpb.NewUserServiceClient(conn)
	return nil
}

func (uc *UserClient) GetUsers(c *context.Context) (*[]User, error) {

	if err := prepareUserGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := userGrpcServiceClient.GetUsers(*c, &userpb.GetUsersRequest{})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}

	var users []User
	for _, u := range res.GetUsers() {
		users = append(users, User{Id: u.Id, Name: u.Name, Age: u.Age, Greeting: u.Greeting})
	}
	return &users, nil
}
