package users

import ( 
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model

	FullName 		string				 `gorm:"notNull"`
	Ktp					string				 `gorm:"notNull;index"`
	Password		string				 `gorm:"notNull"`
	Email		  	string				 `gorm:"notNull"`
	Age					int					   `gorm:"notNull"`
	AddressID		int
	Address			AddressModel	 `gorm:"foreignKey:AddressID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Phone				string				 `gorm:"notNull"`
	IsActive 		bool					 `gorm:"notNull;default:true"`
	LastLoginAt	time.Time			 `gorm:"notNull"`
	RolesID			int
	Roles				UserRolesModel `gorm:"foreignKey:RolesID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type AddressModel struct {
	gorm.Model

	FullAddress	string	`gorm:"notNull"`
	City 				string	`gorm:"notNull"`
	SubDistrict	string	`gorm:"notNull"`
	PostalCode	string	`gorm:"notNull"`
}

type UserRolesModel struct {
	gorm.Model

	RolesName		string	`gorm:"notNull"`
}