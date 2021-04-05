package complaint

import (
	"sipmas-api/src/apps/users"
	"gorm.io/gorm"
)

type ComplaintModel struct {
	gorm.Model

	UserID		int	
	User			users.UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Title			string					`gorm:"notNull"`
	Content		string					`gorm:"notNull"`
	Image			[]byte
	Status		string					`gorm:"default:on process"`
}