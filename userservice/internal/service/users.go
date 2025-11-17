package service

import (
	"context"
	v1 "userservice/api/userservice/v1"
	"userservice/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServiceServer
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.UserReply, error) {
    user, err := s.uc.Create(ctx, req)
    if err != nil {
        return nil, err
    }
    return &v1.UserReply{
        Id:    user.Id,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.UserReply, error) {
    user, err := s.uc.Get(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return &v1.UserReply{
        Id:    user.Id,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UserReply, error) {
    user, err := s.uc.Update(ctx, req)
    if err != nil {
        return nil, err
    }
    return &v1.UserReply{
        Id:    user.Id,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}


func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	success, err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteUserReply{Success: success}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	users, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.ListUsersReply{Users: users}, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *v1.LoginUserRequest) (*v1.AuthReply, error) {
	return s.uc.Login(ctx, req)
}

