package controllers

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lovehotel24/booking-service/pkg/configs"
	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"github.com/lovehotel24/booking-service/pkg/models"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
}

func getUserById(userID interface{}) models.User {
	var user models.User
	configs.DB.Where("id = ?", userID).First(&user)
	return user
}

func getUserByPhone(phone string) models.User {
	var user models.User
	configs.DB.Where("phone = ?", phone).First(&user)
	return user
}

func (u *UserService) CreateUser(ctx context.Context, userReq *userpb.CreateUserRequest) (*emptypb.Empty, error) {

	req := userReq.GetUser()

	Id, err := uuid.Parse(req.GetId().GetValue())
	if err != nil {
		return nil, err
	}

	newUser := models.User{
		Id:    Id,
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Role:  req.GetRole(),
	}

	return &emptypb.Empty{}, configs.DB.Create(&newUser).Error
}

func (u *UserService) GetUser(ctx context.Context, userReq *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {

	user := getUserById(userReq.GetId().GetValue())

	userInfo := userpb.User{
		Id:    &userpb.UUID{Value: user.Id.String()},
		Name:  user.Name,
		Phone: user.Phone,
		Role:  user.Role,
	}

	return &userpb.GetUserResponse{
		User: &userInfo,
	}, nil
}

func (u *UserService) UpdateUser(ctx context.Context, userReq *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {

	req := userReq.GetUser()
	user := getUserById(req.GetId().GetValue())

	req.GetRole()

	if req.GetName() != "" {
		user.Name = req.GetName()
	}

	if req.GetRole() != "" && req.GetRole() != user.Role && user.Role == "ADMIN" {
		user.Role = req.GetRole()
	}

	return &userpb.UpdateUserResponse{
		User: userReq.GetUser(),
	}, configs.DB.Save(&user).Error
}

func (u *UserService) DeleteUser(ctx context.Context, userReq *userpb.DeleteUserRequest) (*emptypb.Empty, error) {

	var user models.User
	return &emptypb.Empty{}, configs.DB.Where("id = ?", userReq.GetId().GetValue()).Delete(&user).Error
}

func (u *UserService) GetAllUsers(ctx context.Context, userReq *userpb.GetAllUserRequest) (*userpb.GetAllUserResponse, error) {

	limit := userReq.GetLimit()
	offset := userReq.GetOffset()

	var users []models.User
	var allUsers []*userpb.User

	configs.DB.Limit(int(limit)).Offset(int(offset)).Find(&users)

	for _, v := range users {
		user := &userpb.User{
			Id:    &userpb.UUID{Value: v.Id.String()},
			Name:  v.Name,
			Phone: v.Phone,
			Role:  v.Role,
		}
		allUsers = append(allUsers, user)
	}

	return &userpb.GetAllUserResponse{Users: allUsers}, nil
}
