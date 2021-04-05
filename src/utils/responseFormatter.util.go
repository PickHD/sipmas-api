package utils

import (
	"github.com/gin-gonic/gin"
)

func ResponseFormatter(code int,message string,err error ,result []interface{},c *gin.Context){
	if code < 400 {
		c.JSON(code,gin.H{
			"code":code,
			"success":true,
			"message":message,
			"error":nil,
			"data":result,
		})
		return
	}
	c.JSON(code,gin.H{
		"code":code,
		"success":false,
		"message":message,
		"error":err,
		"data":result,
	})
}