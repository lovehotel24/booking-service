package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	Id            uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	RoomId        uuid.UUID `gorm:"type:uuid;"`
	UserId        uuid.UUID `gorm:"type:uuid;"`
	BookStartDate time.Time
	BookEndDate   time.Time
	CheckInTime   time.Time
	CheckOutTime  time.Time
	PaymentStatus bool
}

func (book *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	book.Id, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	book.CheckInTime = setTime("13:00")
	book.CheckOutTime = setTime("11:00")
	return nil
}

func setTime(t string) time.Time {
	Layout := "15:04"
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println(err)
	}

	some, err := time.ParseInLocation(Layout, t, location)
	if err != nil {
		fmt.Println(err)
	}
	return some
}
