package complaint

import (
	"sipmas-api/src/apps/users"
	"gorm.io/gorm"
)

type ComplaintModel struct {
	gorm.Model

  UserID    int
	User			users.UserModel 
	Title			string					`gorm:"notNull"`
	Content		string					`gorm:"notNull"`
	Image			[]byte
	Status		string					`gorm:"default:PROSES"`
}