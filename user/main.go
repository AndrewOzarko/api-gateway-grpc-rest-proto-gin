package main

import (
	"context"
	"gateway/user/userdb"
	"gateway/user/userpb"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	timeout      = time.Second
	mongo_client *mongo.Client
)

type server struct {
	userpb.UnimplementedUserServiceServer
}

func (s *server) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	log.Println("Called GetUsers")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	users, err := userdb.Find(mongo_client, c)

	if err != nil {
		return nil, error_response(err)
	}

	var resp userpb.GetUsersResponse

	for _, d := range *users {
		resp.Users = append(resp.Users, &userpb.GetUserResponse{Id: d.ID.Hex(), Name: d.Name, Age: d.Age, Greeting: d.Greeting})
	}

	return &resp, nil

}

func error_response(err error) error {
	log.Println("ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func main() {
	log.Println("User Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

	mongo_client, err = userdb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer mongo_client.Disconnect(context.Background())

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
}
