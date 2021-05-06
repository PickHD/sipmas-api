package config

import (
  "fmt"
  "os"

  "github.com/go-redis/redis/v8"
)

func ConnectRedis() (*redis.Client,error) {

  opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s:%s/%s",os.Getenv("RDS_USER"),os.Getenv("RDS_PASS"),os.Getenv("RDS_HOST"),os.Getenv("REDIS_PORT"),os.Getenv("RDS_NAME")))

  if err != nil {
    panic(err)
  }

  rdb := redis.NewClient(opt)

  fmt.Println("Redis Connected Successfully !")

  return rdb,nil
}
