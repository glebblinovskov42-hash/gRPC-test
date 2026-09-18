package storage

import (
	"context"

	pb "example/grpc_test/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) GetUser(ctx context.Context, id int32) (*pb.UserResponse, error) {
	var name string
	var age int32
	var balance int32

	err := s.pool.QueryRow(ctx,
		"SELECT name, age, balance FROM users WHERE id = $1", id,
	).Scan(&name, &age, &balance)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{Name: name, Age: age, Balance: balance}, nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int32) error {
	_, err := s.pool.Exec(ctx,
		"DELETE FROM users WHERE id = $1", id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Transfer(ctx context.Context, fromId, toId, amount int32) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		"UPDATE users SET balance = balance - $1 WHERE id = $2", amount, fromId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		"UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
