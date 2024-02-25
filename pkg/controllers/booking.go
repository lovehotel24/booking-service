package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
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
		BookStartDate: datatypes.Date(time.Now()),
		BookEndDate:   datatypes.Date(time.Now().Add(12 * time.Hour)),
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

func (a API) GetAllBooking(ctx context.Context, request routers.GetAllBookingRequestObject) (routers.GetAllBookingResponseObject, error) {
	var book []routers.Booking

	a.DB.Find(&book)
	return routers.GetAllBooking200JSONResponse(book), nil
}

func (a API) GetBookingByUserId(ctx context.Context, request routers.GetBookingByUserIdRequestObject) (routers.GetBookingByUserIdResponseObject, error) {
	userId := request.UserId
	var book []routers.Booking
	configs.DB.Where("user_id = ?", userId).First(&book)
	return routers.GetBookingByUserId200JSONResponse(book), nil
}

func (a API) DeleteBookingById(ctx context.Context, request routers.DeleteBookingByIdRequestObject) (routers.DeleteBookingByIdResponseObject, error) {
	bookId, _ := uuid.Parse(request.BookId)
	book := getBookingById(bookId)
	a.DB.Delete(&book)
	return routers.DeleteBookingById204Response{}, nil
}

func (a API) UpdateBookingById(ctx context.Context, request routers.UpdateBookingByIdRequestObject) (routers.UpdateBookingByIdResponseObject, error) {
	bookId, _ := uuid.Parse(request.BookId)
	book := getBookingById(bookId)

	if request.Body.RoomId != uuid.Nil {
		book.RoomId = request.Body.RoomId
	}

	a.DB.Save(&book)

	return routers.UpdateBookingById200JSONResponse{Id: &bookId}, nil

}

func getBookingById(bookId interface{}) routers.Booking {
	var book routers.Booking
	configs.DB.Where("id = ?", bookId).First(&book)
	return book
}

func NewApp(db *gorm.DB) *fiber.App {
	api := &API{DB: db}
	app := fiber.New()

	server := routers.NewStrictHandler(api, nil)

	routers.RegisterHandlers(app, server)

	return app
}
