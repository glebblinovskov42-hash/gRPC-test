package main

import (
	"context"
	pb "example/grpc_test/user"
	"example/storage"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	storage *storage.Storage
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return s.storage.GetUser(ctx, req.Id)
}

func (s *server) DeleteUser(ctx context.Context, req *pb.UserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.storage.DeleteUser(ctx, req.Id); err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{Id: req.Id}, nil
}

func (s *server) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	if err := s.storage.Transfer(ctx, req.FromId, req.ToId, req.Amount); err != nil {
		return nil, err
	}
	return &pb.TransferResponse{Succes: true, Message: "operation succes"}, nil
}

func main() {
	pool, err := pgxpool.New(context.Background(),
		"postgres://postgres:sos11982@localhost:5432/grpc_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println("database connected")

	st := storage.NewStorage(pool)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic("error")
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{storage: st})
	fmt.Println("server started on :50051")

	go func() {
		if err := s.Serve(lis); err != nil {
			panic("error")
		}
	}()

	<-ctx.Done()
	fmt.Println("app stop correctly")
}
