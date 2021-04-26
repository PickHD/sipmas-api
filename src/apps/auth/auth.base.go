package auth

import (
  "gorm.io/gorm"
  "github.com/go-redis/redis/v8"
)

type AuthBaseHandler struct {
  postDB *gorm.DB
  rdsDB *redis.Client
}

func NewAuthBaseHandler(db *gorm.DB,rds *redis.Client) *AuthBaseHandler {
  return &AuthBaseHandler{
    postDB: db,
    rdsDB: rds,
  }
}