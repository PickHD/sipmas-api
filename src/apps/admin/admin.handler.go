package admin

import (
	"fmt"
	"net/http"
  "strconv"

	u "sipmas-api/src/utils"

	"github.com/gin-gonic/gin"
)

func (h *AdminBaseHandler) DashboardHandler(c *gin.Context) {

  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  _,err=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  getStatUserAndComplaint,err:=CountAllUserAndComplaint(h.postDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("CountAllUserAndComplaint Error : %s",err.Error()),err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusOK,"Selamat Datang Admin!",nil,getStatUserAndComplaint,c)

}

func (h *AdminBaseHandler) ManageHandler(c *gin.Context) {

  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  _,err=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  manageLinks:=map[string]interface{}{
    "manageUserLink":"http://localhost:35401/api/v1/admin/kelola/users",
    "manageComplaintLink":"http://localhost:35401/api/v1/admin/kelola/pengaduan",
  }

  u.ResponseFormatter(http.StatusOK,"Daftar Kelola",nil,manageLinks,c)

}

func (h *AdminBaseHandler) ManageGetAllComplaintHandler(c *gin.Context) {
  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  _,err=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  complaints,err:=FetchAllComplaint(h.postDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("FetchAllComplaint Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{
    "allComplaints":complaints,
  }

  u.ResponseFormatter(http.StatusOK,"Semua Pengaduan",nil,payloadData,c)
}

func (h *AdminBaseHandler) ManageGetOneComplaintHandler(c *gin.Context) {

  getCompID,err:=strconv.Atoi(c.Param("id"))
  if err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }
  
  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  _,err=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  getOneComplaint,err:=FetchOneComplaint(h.postDB,uint(getCompID))
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("FetchOneComplaint Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{
    "specComplaint":getOneComplaint,
  }

  u.ResponseFormatter(http.StatusOK,"Detail Pengaduan",nil,payloadData,c)

}

func (h *AdminBaseHandler) ManageUpdateComplaintHandler(c *gin.Context) {
  var validComplaint ManageComplaintValidation

  if err:=c.ShouldBindJSON(&validComplaint);err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }

  getCompID,err:=strconv.Atoi(c.Param("id"))
  if err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }

  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  _,err=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  _,err=UpdateComplaint(h.postDB,uint(getCompID),validComplaint)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("UpdateComplaint Error : %s",err.Error()),err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusOK,"Pengaduan berhasil di perbaharui",nil,nil,c)

}

func (h *AdminBaseHandler) ManageGetAllUserHandler(c *gin.Context) {
  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  _,err=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  users,err:=FetchAllUser(h.postDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("FetchAllUser Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{
    "allUsers":users,
  }

  u.ResponseFormatter(http.StatusOK,"Semua Pengguna",nil,payloadData,c)
}

func (h *AdminBaseHandler) ManageGetOneUserHandler(c *gin.Context) {
  
  getUserID,err:=strconv.Atoi(c.Param("id"))
  if err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }
  
  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  _,err=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  getOneUser,err:=FetchOneUser(h.postDB,uint(getUserID))
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("FetchOneComplaint Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{
    "specUser":getOneUser,
  }

  u.ResponseFormatter(http.StatusOK,"Detail Pengguna",nil,payloadData,c)
}

func (h *AdminBaseHandler) ManageUpdateUserHandler(c *gin.Context) {
  var validUser ManageUserValidation

  if err:=c.ShouldBindJSON(&validUser);err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }

  getUserID,err:=strconv.Atoi(c.Param("id"))
  if err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }

  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  _,err=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  _,err=UpdateUser(h.postDB,uint(getUserID),validUser)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("UpdateUser Error : %s",err.Error()),err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusOK,"Pengguna berhasil di perbaharui",nil,nil,c)
}

