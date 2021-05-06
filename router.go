package main

import (
	"fmt"
	"net/http"
	"time"
  
	"sipmas-api/src/apps/auth"
	"sipmas-api/src/apps/users"
	u "sipmas-api/src/utils"
  c "sipmas-api/src/config"
  m "sipmas-api/src/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Router returning a *gin.Engine
func Router() *gin.Engine {
  //!Setup New Gin 
	r:=gin.New()

  //!Setup Max Memory Uploading 
  r.MaxMultipartMemory = 8 << 20

  //!Setup CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{"GET","POST","PUT","DELETE"},
		AllowHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

  //!Setup Default Logger & Recovery From Gin 
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

  //!Connecting To Postgres DB & Redis DB 
  db,err:=c.Connect()
	if err!=nil{
		fmt.Printf("u.Connect() failed with %s\n",err)
	}
  rds,err:=c.ConnectRedis()
  if err!=nil{
    fmt.Printf("u.ConnectRedis() failed with %s\n",err)
  }

  //!Setup Route For Pinging Server 
	r.GET("/ping",func(c *gin.Context){
		u.ResponseFormatter(http.StatusOK,"Pong",nil,nil,c)
	})

  //!Setup Router For Not Found Route Handling
	r.NoRoute(func(c *gin.Context){
		u.ResponseFormatter(http.StatusNotFound,"Route not found.",nil,nil,c)
	})

  //!Setup New Base Handler From Each Apps 
  hAuth:=auth.NewAuthBaseHandler(db,rds)
  hUser:=users.NewUserBaseHandler(db,rds)

  //!Setup Super Nested Group Route For API's
	superGroupv1:=r.Group("/api/v1")
	{
    //!Auth Group Section 
		auth:=superGroupv1.Group("/auth")
		{
			auth.POST("/daftar",hAuth.SignupHandler)
			auth.GET("/konfirmasi",hAuth.ConfirmAccountHandler)
			auth.POST("/masuk",hAuth.SigninHandler)
			auth.GET("/keluar",m.JWTAuthMiddleware(),hAuth.SignoutHandler)
			auth.POST("/token/refresh",hAuth.TokenHandler)
		}

    //!Users Group Section
		users:=superGroupv1.Group("/users")
		{
			users.GET("/dashboard",m.JWTAuthMiddleware(),hUser.DashboardHandler)
			users.GET("/profile",m.JWTAuthMiddleware(),hUser.ProfileHandler)

			users.POST("/pengaduan",m.JWTAuthMiddleware(),hUser.CreateComplaintHandler)
      users.POST("/pengaduan/:id/unggah-foto",m.JWTAuthMiddleware(),hUser.UploadComplaintImageHandler)
      users.GET("/pengaduan",m.JWTAuthMiddleware(),hUser.GetAllComplaintHandler)
			users.GET("/pengaduan/:id",m.JWTAuthMiddleware(),hUser.GetOneComplaintHandler)
			users.PUT("/pengaduan/:id",m.JWTAuthMiddleware(),hUser.UpdateComplaintHandler)
			users.DELETE("/pengaduan/:id",m.JWTAuthMiddleware(),hUser.DeleteComplaintHandler)
      
		}

    //!Admin Group Section 
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