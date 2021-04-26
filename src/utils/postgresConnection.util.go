package utils

import (
	"fmt"
	"os"
	"sipmas-api/src/apps/complaint"
	"sipmas-api/src/apps/users"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB,error) {

	dsn:= fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",os.Getenv("DB_HOST"),os.Getenv("DB_USER"),os.Getenv("DB_PASS"),os.Getenv("DB_NAME"),os.Getenv("POSTGRES_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err!=nil{
		return nil,err
	}

	fmt.Println("Database Connected Successfully !")

	db.AutoMigrate(&users.UserModel{},&users.AddressModel{},&complaint.ComplaintModel{})

	fmt.Println("Database Migrated Successfully !")

	return db,nil
}