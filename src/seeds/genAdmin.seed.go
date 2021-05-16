package seeds

import (
	"errors"
	"os"

	user "sipmas-api/src/apps/users"

	"gorm.io/gorm"
)

func GenerateAdmin(db *gorm.DB)(error) {

  newAdmin:= user.UserModel {
    FullName: os.Getenv("ADMIN_NAME"),
    NIK:"3273215901120004",
    Email:os.Getenv("ADMIN_EMAIL"),
    Password: os.Getenv("ADMIN_PASS"),
    Age:19,
    Phone:os.Getenv("ADMIN_PHONE"),
    Address: user.AddressModel {
      FullAddress:"Komp PLN GG.V Ciateul no.45",
      City:"Kota Bandung",
      SubDistrict:"Regol",
      PostalCode:"40252",  
    },
    Roles:"ADMIN",
    IsVerified: true,
  }

  if err:=db.First(&newAdmin,"email=?",newAdmin.Email).Error;err!=nil{
    db.Create(&newAdmin)
    return nil
  }else {
    return errors.New("admin already existed")
  }
}
