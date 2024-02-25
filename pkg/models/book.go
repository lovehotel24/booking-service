package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	Id            uuid.UUID      `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	RoomId        uuid.UUID      `json:"roomId" gorm:"type:uuid;"`
	UserId        uuid.UUID      `json:"userId" gorm:"type:uuid;"`
	BookStartDate datatypes.Date `json:"bookStartDate" gorm:"type:date"`
	BookEndDate   datatypes.Date `json:"bookEndDate" gorm:"type:date"`
	CheckInTime   datatypes.Time `json:"checkInTime" gorm:"type:time"`
	CheckOutTime  datatypes.Time `json:"checkOutTime" gorm:"type:time"`
	PaymentStatus bool           `json:"paymentStatus" gorm:"default:true"`
}

func (book *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	book.Id, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	book.CheckInTime = datatypes.NewTime(13, 0, 0, 0)
	book.CheckOutTime = datatypes.NewTime(11, 0, 0, 0)
	return nil
}
