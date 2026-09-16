package main

import (
	"context"
	pb "example/grpc_test/user"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("err")
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &pb.UserRequest{Id: 1})
	if err != nil {
		panic("err")
	}

	fmt.Printf("User: %s, Age: %d", resp.Name, resp.Age)
}
