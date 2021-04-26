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
      u.ResponseFormatter(http.StatusUnauthorized,"Unauthorized",err,nil,c)
      c.Abort()
      return
    }
    c.Next()
  }
}

