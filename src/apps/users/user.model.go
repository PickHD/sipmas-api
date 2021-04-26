package users

import ( 
	"time"

	"gorm.io/gorm"
  "golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	gorm.Model

	FullName 		string				 `gorm:"notNull"`
	Ktp					string				 `gorm:"notNull"`
	Password		string				 `gorm:"notNull"`
	Email		  	string				 `gorm:"notNull;index"`
	Age					int					   `gorm:"notNull"`
  AddressID   int 
	Address			AddressModel	 
	Phone				string				 `gorm:"notNull"`
	IsActive 		bool					 `gorm:"notNull;default:true"`
  IsVerified  bool           `gorm:"notNull;default:false"`
	LastLoginAt	time.Time
	Roles				string         `gorm:"notNull;default:USER"`
}

type AddressModel struct {
	gorm.Model

	FullAddress	string	`gorm:"notNull"`
	City 				string	`gorm:"notNull"`
	SubDistrict	string	`gorm:"notNull"`
	PostalCode	string	`gorm:"notNull"`
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
  return u.HashPassword(tx)
}

func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
  if tx.Statement.Changed("Password") {
    return u.HashPassword(tx)
  }

  return nil
}

func (u *UserModel) HashPassword(tx *gorm.DB) error {
  var newPass string

  switch u := tx.Statement.Dest.(type) {

    case map[string]interface{}:
      newPass = u["password"].(string)
    case *UserModel:
      newPass = u.Password
    case []*UserModel:
      newPass = u[tx.Statement.CurDestIndex].Password

  }

  b, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
  if err != nil {
      return err
  }
  tx.Statement.SetColumn("password", b)

  return nil
}