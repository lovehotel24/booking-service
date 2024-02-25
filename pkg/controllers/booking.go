package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/lovehotel24/booking-service/pkg/routers"
)

type API struct {
	DB *gorm.DB
}

func NewApp(db *gorm.DB, log *logrus.Logger) *fiber.App {
	api := &API{DB: db}
	app := fiber.New()
	app.Use(LoggerMiddleware(log))

	server := routers.NewStrictHandler(api, nil)
	routers.RegisterHandlers(app, server)

	return app
}

func (a API) CreateBooking(ctx context.Context, request routers.CreateBookingRequestObject) (routers.CreateBookingResponseObject, error) {
	errMsg := fmt.Sprintf("failed to create booking")

	booking := routers.Booking{
		RoomId:        request.Body.RoomId,
		UserId:        request.Body.UserId,
		BookStartDate: request.Body.BookStartDate,
		BookEndDate:   request.Body.BookEndDate,
		PaymentStatus: request.Body.PaymentStatus,
	}

	if err := a.DB.Create(&booking).Error; err != nil {
		return routers.CreateBookingdefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	}

	response := routers.CreateBooking200JSONResponse{}
	response.Id = &booking.Id

	return response, nil
}

func (a API) GetBookingById(ctx context.Context, request routers.GetBookingByIdRequestObject) (routers.GetBookingByIdResponseObject, error) {
	booking := routers.Booking{}
	errMsg := fmt.Sprintf("failed to get booking id: %s", request.BookId)

	if err := a.DB.Where("id = ?", request.BookId).First(&booking).Error; err != nil {
		return routers.GetBookingByIddefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	}

	return routers.GetBookingById200JSONResponse(booking), nil
}

func (a API) GetAllBooking(ctx context.Context, request routers.GetAllBookingRequestObject) (routers.GetAllBookingResponseObject, error) {
	var book []routers.Booking
	errMsg := fmt.Sprintf("failed to get booking")

	limit := 10
	offSet := 1

	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	if request.Params.Offset != nil {
		offSet = *request.Params.Offset
	}

	if err := a.DB.Limit(limit).Offset(offSet).Find(&book).Error; err != nil {
		return routers.GetAllBookingdefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	}

	return routers.GetAllBooking200JSONResponse(book), nil
}

func (a API) GetBookingByUserId(ctx context.Context, request routers.GetBookingByUserIdRequestObject) (routers.GetBookingByUserIdResponseObject, error) {
	userId := request.UserId
	var book []routers.Booking
	errMsg := fmt.Sprintf("failed to get booking")

	if err := a.DB.Where("user_id = ?", userId).Find(&book).Error; err != nil {
		return routers.GetBookingByUserIddefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	}

	return routers.GetBookingByUserId200JSONResponse(book), nil
}

func (a API) DeleteBookingById(ctx context.Context, request routers.DeleteBookingByIdRequestObject) (routers.DeleteBookingByIdResponseObject, error) {
	bookId, _ := uuid.Parse(request.BookId)
	errMsg := fmt.Sprintf("failed to delete booking id: %s", bookId)

	book, err := a.getBookingById(bookId)
	if err != nil {
		return routers.DeleteBookingByIddefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	}

	if err := a.DB.Delete(&book).Error; err != nil {
		return routers.DeleteBookingByIddefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	}

	return routers.DeleteBookingById204Response{}, nil
}

func (a API) UpdateBookingById(ctx context.Context, request routers.UpdateBookingByIdRequestObject) (routers.UpdateBookingByIdResponseObject, error) {
	bookId, _ := uuid.Parse(request.BookId)
	errMsg := fmt.Sprintf("failed to update booking id: %s", bookId)

	book, err := a.getBookingById(bookId)
	if err != nil {
		return routers.UpdateBookingByIddefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	}

	if request.Body.RoomId != uuid.Nil {
		book.RoomId = request.Body.RoomId
	}

	if request.Body.PaymentStatus {
		book.PaymentStatus = request.Body.PaymentStatus
	}

	if date, err := request.Body.BookStartDate.Value(); err != nil {
		return routers.UpdateBookingByIddefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	} else if date != "" {
		book.BookStartDate = request.Body.BookStartDate
	}

	if date, err := request.Body.BookEndDate.Value(); err != nil {
		return routers.UpdateBookingByIddefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	} else if date != "" {
		book.BookEndDate = request.Body.BookEndDate
	}

	if err := a.DB.Save(&book).Error; err != nil {
		return routers.UpdateBookingByIddefaultJSONResponse{Body: routers.Error{Message: &errMsg}, StatusCode: http.StatusBadRequest}, err
	}

	return routers.UpdateBookingById200JSONResponse{Id: &bookId}, nil
}

func (a API) getBookingById(bookId interface{}) (routers.Booking, error) {
	var book routers.Booking

	if err := a.DB.Where("id = ?", bookId).First(&book).Error; err != nil {
		return routers.Booking{}, err
	}

	return book, nil
}

func LoggerMiddleware(log *logrus.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		log.WithFields(logrus.Fields{
			"method":   c.Method(),
			"path":     c.Path(),
			"duration": fmt.Sprintf("%v", time.Since(start)),
		}).Info("Request handled")

		return err
	}
}
