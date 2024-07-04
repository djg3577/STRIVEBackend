package user

import (
    "context"

    pb "github.com/djg3577/STRIVEBackend/api/user"
)

type UserService struct {
    pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
    return &UserService{}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    // Logic for creating a user
    return &pb.CreateUserResponse{UserId: "123"}, nil
}
