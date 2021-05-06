package users

import ( 
	"time"
  "errors"

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

type ComplaintModel struct {
	gorm.Model

  UserID    int
	User			UserModel 
	Title			string					`gorm:"notNull"`
	Content		string					`gorm:"notNull"`
  ImageName string
	Status		string					`gorm:"default:PROSES"`
}

type ComplaintValidation struct {
  Title    string   `json:"title" binding:"required"`
  Content  string   `json:"content" binding:"required"`
}

type UpdateComplaintValidation struct {
  Title       string    `json:"title"`
  Content     string    `json:"content"`
}

//!USER MODEL HOOKS 
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
  return u.HashPassword(tx)
}

func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {

  if tx.Statement.Changed("Password") {
    return u.HashPassword(tx)
  }

  if u.Roles == "ADMIN" {
    return errors.New("admin not allowed to update")
  }

  if tx.Statement.Changed("Role") {
    return errors.New("role not allowed to change")
  }

  return nil
}

func (u *UserModel) BeforeDelete(tx *gorm.DB) error {
  if u.Roles == "USER" {
    return errors.New("user not allowed to delete")
  }

  return nil
}

//!COMPLAINT MODEL HOOKS
func (c *ComplaintModel) BeforeCreate(tx *gorm.DB) error {

  if c.User.Roles == "ADMIN" {
    return errors.New("admin not allowed to create")
  }

  return nil
}

func (c *ComplaintModel) BeforeUpdate(tx *gorm.DB) error {

  if c.User.Roles == "USER" && tx.Statement.Changed("Status") {
    return errors.New("user not allowed to change status")
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