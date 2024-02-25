package controllers

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"github.com/lovehotel24/booking-service/pkg/models"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	userService := &UserService{
		DB: db,
	}

	return userService
}

func (u *UserService) getUserById(userID interface{}) (models.User, error) {
	var user models.User
	if err := u.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *UserService) getUserByPhone(phone string) (models.User, error) {
	var user models.User
	if err := u.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *UserService) CreateUser(ctx context.Context, userReq *userpb.CreateUserRequest) (*emptypb.Empty, error) {

	req := userReq.GetUser()

	id, err := uuid.Parse(req.GetId().GetValue())
	if err != nil {
		return nil, err
	}

	newUser := models.User{
		Id:    id,
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Role:  req.GetRole(),
	}

	return &emptypb.Empty{}, u.DB.Create(&newUser).Error
}

func (u *UserService) GetUser(ctx context.Context, userReq *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := u.getUserById(userReq.GetId().GetValue())
	if err != nil {
		return &userpb.GetUserResponse{}, err
	}

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

	user, err := u.getUserById(req.GetId().GetValue())
	if err != nil {
		return &userpb.UpdateUserResponse{}, err
	}

	req.GetRole()

	if req.GetName() != "" {
		user.Name = req.GetName()
	}

	if req.GetRole() != "" && req.GetRole() != user.Role && user.Role == "ADMIN" {
		user.Role = req.GetRole()
	}

	return &userpb.UpdateUserResponse{
		User: userReq.GetUser(),
	}, u.DB.Save(&user).Error
}

func (u *UserService) DeleteUser(ctx context.Context, userReq *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
	var user models.User
	return &emptypb.Empty{}, u.DB.Where("id = ?", userReq.GetId().GetValue()).Delete(&user).Error
}

func (u *UserService) GetAllUsers(ctx context.Context, userReq *userpb.GetAllUserRequest) (*userpb.GetAllUserResponse, error) {

	limit := userReq.GetLimit()
	offset := userReq.GetOffset()

	var users []models.User
	var allUsers []*userpb.User

	if err := u.DB.Limit(int(limit)).Offset(int(offset)).Find(&allUsers).Error; err != nil {
		return &userpb.GetAllUserResponse{}, err
	}

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
