package middlewares

import (
	"net/http"

	user "sipmas-api/src/apps/users"
	u "sipmas-api/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func AdminOnlyMiddleware(db *gorm.DB,rds *redis.Client) gin.HandlerFunc {
  var checkAdmin user.UserModel

  return func(c *gin.Context) {
  
    ad,err:=u.ExtractTokenMetadata(c.Request)
    if err!=nil{
      u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
      c.Abort()
      return
    }

    userID,err:=u.FetchAuth(ad,rds)
    if err!=nil{
      u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
      c.Abort()
      return
    }

    if err=db.First(&checkAdmin,"roles=? AND id=?","ADMIN",uint(userID)).Error;err!=nil{
      u.ResponseFormatter(http.StatusForbidden,"Hanya Admin yang bisa akses route ini, silahkan untuk kembali.",err,nil,c)
      c.Abort()
      return
    }
    
    c.Next()
  }


}