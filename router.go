package main

import (
	"fmt"
	"net/http"
	"sipmas-api/src/apps/auth"
	u "sipmas-api/src/utils"
  m "sipmas-api/src/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r:=gin.New()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{"GET","POST","PUT","DELETE"},
		AllowHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

  db,err:=u.Connect()
	if err!=nil{
		fmt.Printf("u.Connect() failed with %s\n",err)
	}

  rds,err:=u.ConnectRedis()
  if err!=nil{
    fmt.Printf("u.ConnectRedis() failed with %s\n",err)
  }

	r.GET("/ping",func(c *gin.Context){
		u.ResponseFormatter(http.StatusOK,"Pong",nil,nil,c)
	})

	r.NoRoute(func(c *gin.Context){
		u.ResponseFormatter(http.StatusNotFound,"Route not found.",nil,nil,c)
	})

  hAuth:=auth.NewAuthBaseHandler(db,rds)

	superGroupv1:=r.Group("/api/v1")
	{
		auth:=superGroupv1.Group("/auth")
			{
				auth.POST("/daftar",hAuth.SignupHandler)
				auth.GET("/konfirmasi",hAuth.ConfirmAccountHandler)
				auth.POST("/masuk",hAuth.SigninHandler)
				auth.GET("/keluar",m.JWTAuthMiddleware(),hAuth.SignoutHandler)
				auth.POST("/token/refresh",hAuth.TokenHandler)
			}

		// users:=superGroupv1.Group("/users")
		// {
		// 	users.GET("/dashboard")
		// 	users.GET("/profile")
		// 	users.POST("/suspensasi")

		// 	users.POST("/pengaduan")
		// 	users.GET("/pengaduan")
		// 	users.GET("/pengaduan/:id")
		// 	users.PUT("/pengaduan/:id")
		// 	users.DELETE("/pengaduan/:id")
		// }

		// admin:=superGroupv1.Group("/admin")
		// {
		// 	admin.GET("/dashboard")
			
		// 	admin.POST("/kelola/pengaduan")
		// 	admin.GET("/kelola/pengaduan")
		// 	admin.GET("/kelola/pengaduan/:id")
		// 	admin.PUT("/kelola/pengaduan/:id")
		// 	admin.DELETE("/kelola/pengaduan/:id")

		// 	admin.POST("/kelola/users")
		// 	admin.GET("/kelola/users")
		// 	admin.GET("/kelola/users/:id")
		// 	admin.PUT("/kelola/users/:id")
		// 	admin.DELETE("/kelola/users/:id")
		// }
	}

	return r
}