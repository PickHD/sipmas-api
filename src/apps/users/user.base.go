package users

import (
  "gorm.io/gorm"
  "github.com/go-redis/redis/v8"
)

type UserBaseHandler struct {
  DB *gorm.DB
  rdsDB *redis.Client
}

func NewUserBaseHandler(db *gorm.DB,rds *redis.Client) *UserBaseHandler {
  return &UserBaseHandler{
    DB:db,
    rdsDB: rds,
  }
}