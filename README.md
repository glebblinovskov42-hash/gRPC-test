# gRPC User Service

gRPC-сервис для управления пользователями на Go.

## Стек
- Go, gRPC, Protobuf
- PostgreSQL, pgx
- Docker, docker-compose
- Тесты (testing)

## Запуск
make service-run

## Методы
- GetUser
- DeleteUser
- Transfer (транзакции, FOR UPDATE)

## Тесты
make unit-test
