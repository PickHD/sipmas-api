package middlewares

import (
  "net/http"
	u "sipmas-api/src/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc{

  return func(c *gin.Context) {
    err:=u.ValidateJWT(c.Request)
    if err!=nil{
      u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
      c.Abort()
      return
    }
    c.Next()
  }
}

