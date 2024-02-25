package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/lovehotel24/booking-service/pkg/configs"
	"github.com/lovehotel24/booking-service/pkg/routers"
)

type API struct {
	DB *gorm.DB
}

func (a API) CreateBooking(ctx context.Context, request routers.CreateBookingRequestObject) (routers.CreateBookingResponseObject, error) {

	booking := routers.Booking{
		RoomId:        request.Body.RoomId,
		UserId:        request.Body.UserId,
		BookStartDate: time.Now(),
		BookEndDate:   time.Now().Add(12 * time.Hour),
		PaymentStatus: true,
	}
	configs.DB.Create(&booking)
	response := routers.CreateBooking200JSONResponse{}
	response.Id = &booking.Id
	return response, nil
}

func (a API) GetBookingById(ctx context.Context, request routers.GetBookingByIdRequestObject) (routers.GetBookingByIdResponseObject, error) {

	booking := routers.Booking{}

	a.DB.Where("id = ?", request.BookId).First(&booking)
	return routers.GetBookingById200JSONResponse(booking), nil
}

func NewApp(db *gorm.DB) *fiber.App {
	api := &API{DB: db}
	app := fiber.New()

	server := routers.NewStrictHandler(api, nil)

	routers.RegisterHandlers(app, server)

	return app
}
