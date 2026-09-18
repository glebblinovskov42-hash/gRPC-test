package main

import (
	"context"
	pb "example/grpc_test/user"
	"example/storage"
	"fmt"
	"net"
	"os"
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

func connect() string {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	return conn
}

func main() {
	conn := connect()

	pool, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println("database connected")

	_, err = pool.Exec(context.Background(), `
		 CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			age INT NOT NULL, 
			balance INT NOT NULL DEFAULT 1000
		)
	`)
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(context.Background(), `
		INSERT INTO users (name, age) VALUES
			('Bob', 20 ),
			('Alice', 30 ),
			('Gunter', 23 )
			ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		panic(err)
	}

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
