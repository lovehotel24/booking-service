package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	Id            uuid.UUID      `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	RoomId        uuid.UUID      `json:"roomId" gorm:"type:uuid;"`
	UserId        uuid.UUID      `json:"userId" gorm:"type:uuid;"`
	BookStartDate time.Time      `json:"bookStartDate" gorm:"type:date"`
	BookEndDate   time.Time      `json:"bookEndDate" gorm:"type:date"`
	CheckInTime   datatypes.Time `json:"checkInTime" gorm:"type:time"`
	CheckOutTime  datatypes.Time `json:"checkOutTime" gorm:"type:time"`
	PaymentStatus bool           `json:"paymentStatus" gorm:"default:false"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	b.CheckInTime = datatypes.NewTime(13, 0, 0, 0)
	b.CheckOutTime = datatypes.NewTime(11, 0, 0, 0)
	return nil
}

func (b *Booking) BeforeSave(tx *gorm.DB) error {
	formattedStartDate := b.BookStartDate.Format("2006-01-02")

	startDate, err := time.Parse("2006-01-02", formattedStartDate)
	if err != nil {
		return err
	}
	b.BookStartDate = startDate

	formattedEndDate := b.BookEndDate.Format("2006-01-02")
	endDate, err := time.Parse("2006-01-02", formattedEndDate)
	if err != nil {
		return err
	}
	b.BookEndDate = endDate
	return nil
}
