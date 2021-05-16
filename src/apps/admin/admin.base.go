package admin

import (
  "gorm.io/gorm"
  "github.com/go-redis/redis/v8"
)

type AdminBaseHandler struct {
  postDB *gorm.DB
  rdsDB *redis.Client
}

func NewAdminBaseHandler(db *gorm.DB,rds *redis.Client) *AdminBaseHandler {
  return &AdminBaseHandler{
    postDB: db,
    rdsDB: rds,
  }
}