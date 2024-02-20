package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id    uuid.UUID `json:"id" gorm:"primary_key;type:uuid;"`
	Name  string    `json:"name"`
	Phone string    `json:"phone" gorm:"<-:create;uniqueIndex"`
	Role  string    `json:"role"`
}

//func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
//	user.Id, err = uuid.NewUUID()
//	if err != nil {
//		return err
//	}
//	return nil
//}
