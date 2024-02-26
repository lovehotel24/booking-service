package controllers

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"github.com/lovehotel24/booking-service/pkg/models"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	DB  *gorm.DB
	Log *logrus.Logger
}

func NewUserService(db *gorm.DB, log *logrus.Logger) *UserService {
	userService := &UserService{
		DB:  db,
		Log: log,
	}

	return userService
}

func (u *UserService) getUserById(userID interface{}) (models.User, error) {
	var user models.User
	if err := u.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Error("failed to get user by userId")
		return models.User{}, err
	}
	return user, nil
}

func (u *UserService) getUserByPhone(phone string) (models.User, error) {
	var user models.User
	if err := u.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"phone": phone,
		}).Error("Failed to get booking by phone")
		return models.User{}, err
	}
	return user, nil
}

func (u *UserService) CreateUser(ctx context.Context, userReq *userpb.CreateUserRequest) (*emptypb.Empty, error) {

	req := userReq.GetUser()

	id, err := uuid.Parse(req.GetId().GetValue())
	if err != nil {
		u.Log.WithError(err).Error("failed to parse user uuid")
		return nil, err
	}

	newUser := models.User{
		Id:    id,
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Role:  req.GetRole(),
	}

	if err := u.DB.Create(&newUser).Error; err != nil {
		u.Log.WithError(err).Error("failed to create user")
		return nil, err
	}

	u.Log.WithFields(logrus.Fields{
		"userID": newUser.Id.String(),
		"name":   newUser.Name,
	}).Info("user created successfully")

	return &emptypb.Empty{}, nil
}

func (u *UserService) GetUser(ctx context.Context, userReq *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := u.getUserById(userReq.GetId().GetValue())
	if err != nil {
		u.Log.WithError(err).Error("failed to get user by id")
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
		u.Log.WithError(err).Error("failed to get user by id")
		return &userpb.UpdateUserResponse{}, err
	}

	req.GetRole()

	if req.GetName() != "" {
		user.Name = req.GetName()
	}

	//todo
	if req.GetRole() != "" && req.GetRole() != user.Role && user.Role == "ADMIN" {
		user.Role = req.GetRole()
	}

	if err := u.DB.Save(&user).Error; err != nil {
		u.Log.WithError(err).Errorf("failed to update user id: %s", user.Id)
		return nil, err
	}

	u.Log.WithFields(logrus.Fields{
		"userID": user.Id.String(),
		"name":   user.Name,
	}).Info("user updated successfully")

	return &userpb.UpdateUserResponse{
		User: userReq.GetUser(),
	}, nil
}

func (u *UserService) DeleteUser(ctx context.Context, userReq *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
	var user models.User
	if err := u.DB.Where("id = ?", userReq.GetId().GetValue()).Delete(&user).Error; err != nil {
		u.Log.WithError(err).Errorf("failed to delete user id: %s", userReq.GetId().GetValue())
		return nil, err
	}

	u.Log.WithFields(logrus.Fields{
		"userID": userReq.GetId().GetValue(),
	}).Info("user deleted successfully")

	return &emptypb.Empty{}, nil
}

func (u *UserService) GetAllUsers(ctx context.Context, userReq *userpb.GetAllUserRequest) (*userpb.GetAllUserResponse, error) {

	limit := userReq.GetLimit()
	offset := userReq.GetOffset()

	var users []models.User
	var allUsers []*userpb.User

	if err := u.DB.Limit(int(limit)).Offset(int(offset)).Find(&users).Error; err != nil {
		u.Log.WithError(err).Error("failed to get users")
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
