package users

import "gorm.io/gorm"

type UserBaseHandler struct {
  DB *gorm.DB
}

func NewUserBaseHandler(db *gorm.DB) *UserBaseHandler {
  return &UserBaseHandler{
    DB:db,
  }
}