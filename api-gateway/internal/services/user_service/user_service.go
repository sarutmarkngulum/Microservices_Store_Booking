package services

import (
	"context"
	"time"

	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	"go.uber.org/zap"
)

type UserService interface {
	Register(ctx context.Context, req *UserRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*Token, error)
	GetAllUsers(ctx context.Context) (*UserList, error)
	GetUser(ctx context.Context, req *UserIdentifier) (*UserResponse, error)
	UpdateUserRole(ctx context.Context, req *UpdateRoleRequest) (*UserResponse, error)
	UpdateUserPassword(ctx context.Context, req *UpdatePasswordRequest) (*UserResponse, error)
	DeleteUser(ctx context.Context, req *UserIdentifier) (*Empty, error)
}

type userService struct {
	userClient UserServiceClient
}

func NewUserService(userClient UserServiceClient) UserService {
	return &userService{userClient: userClient}
}

func (s *userService) Register(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.userClient.Register(ctx, req)
	if err != nil {
		logs.Error("Error calling Register: %v", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *userService) Login(ctx context.Context, req *LoginRequest) (*Token, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.userClient.Login(ctx, req)
	if err != nil {
		logs.Error("Error calling Login: %v", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *userService) GetAllUsers(ctx context.Context) (*UserList, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.userClient.GetAllUsers(ctx, &Empty{})
	if err != nil {
		logs.Error("Error calling GetAllUsers: %v", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *userService) GetUser(ctx context.Context, req *UserIdentifier) (*UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.userClient.GetUser(ctx, req)
	if err != nil {
		logs.Error("Error calling GetUser: %v", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *userService) UpdateUserRole(ctx context.Context, req *UpdateRoleRequest) (*UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.userClient.UpdateUserRole(ctx, req)
	if err != nil {
		logs.Error("Error calling UpdateUserRole: %v", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *userService) UpdateUserPassword(ctx context.Context, req *UpdatePasswordRequest) (*UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.userClient.UpdateUserPassword(ctx, req)
	if err != nil {
		logs.Error("Error calling UpdateUserRole: %v", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *UserIdentifier) (*Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.userClient.DeleteUser(ctx, req)
	if err != nil {
		logs.Error("Error calling DeleteUser: %v", zap.Error(err))
		return nil, err
	}

	return res, nil
}
