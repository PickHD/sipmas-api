package users

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	u "sipmas-api/src/utils"

	"github.com/gin-gonic/gin"
)

func (h *UserBaseHandler) DashboardHandler(c *gin.Context) {

  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  userID,err:=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  fetchComplaintCounts,err:=CountComplaintUser(uint(userID),h.DB)
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Count Complaint Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{
    "profileLink":"http://localhost:35401/api/v1/users/profile",
    "complaintLink":"http://localhost:35401/api/v1/users/pengaduan",
    "complaintUserStatus":map[string]int{
      "PROSES":fetchComplaintCounts[0],
      "DITERMA":fetchComplaintCounts[1],
      "DITOLAK":fetchComplaintCounts[2],
    },
  }

  u.ResponseFormatter(http.StatusOK,"Dashboard Anda",nil,payloadData,c)
}

func (h *UserBaseHandler) GetProfileHandler(c *gin.Context) {
  
  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  userID,err:=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  getUser,err:=FetchUserProfile(uint(userID),h.DB)
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Fetch User Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{"userData":getUser}

  u.ResponseFormatter(http.StatusOK,"Profil Anda",nil,payloadData,c)

}

func (h *UserBaseHandler) UpdateProfileHandler(c *gin.Context) {
  var validUser UpdateUserValidation

  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }

  userID,err:=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  if err:=c.ShouldBindJSON(&validUser);err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }

  _,err=EditProfile(uint(userID),h.DB,validUser)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Edit Profile Error : %s",err.Error()),err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusOK,"Data anda sudah diperbaharui",nil,nil,c)

}

func (h *UserBaseHandler) CreateComplaintHandler(c *gin.Context) {

  var validComplaint ComplaintValidation

  if err:=c.ShouldBindJSON(&validComplaint);err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }

  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }
  
  userID,err:=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  getUser,err:=FetchUserProfile(uint(userID),h.DB)
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Fetch User Error : %s",err.Error()),err,nil,c)
    return
  }

  newComplaint,err:=CreateComplaint(h.DB,getUser,&validComplaint)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Create Complaint Error : %s",err.Error()),err,nil,c)
  }

  payloadData:=map[string]interface{}{
    "complaintData":newComplaint,
  }

  u.ResponseFormatter(http.StatusCreated,"Pengaduan Berhasil Dibuat.",nil,payloadData,c)
}

func (h *UserBaseHandler) UploadComplaintImageHandler(c *gin.Context) {

  getCompID,err:=strconv.Atoi(c.Param("id"))
  if err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }
  
  file,err:=c.FormFile("file")
  if err!=nil{
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Get File Error : %s",err.Error()),err,nil,c)
    return
  }

  currTime := time.Now()

  fileName:=filepath.Base(file.Filename)
  filePath := "./src/upload/"+currTime.String()+"_"+fileName

  if err:=c.SaveUploadedFile(file,filePath);err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("SaveUploadedFile Error : %s",err.Error()),err,nil,c)
    return
  }

  getFileName,err := UpdateImageComplaint(h.DB,uint(getCompID),fileName)
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Update Image Complaint Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{
    "imageName":getFileName,
    "imagePath":filePath,
  }


  u.ResponseFormatter(http.StatusOK,"Unggah Foto Berhasil",nil,payloadData,c)

}

func (h *UserBaseHandler) GetAllComplaintHandler(c *gin.Context) {
  
  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",nil,nil,c)
    return
  }
  
  userID,err:=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  complaints,err:=FetchAllComplaint(h.DB,uint(userID))
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Fetch All Complaint Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{
    "allComplaint":complaints,
  }

  u.ResponseFormatter(http.StatusOK,"Daftar Semua Pengaduan Anda",nil,payloadData,c)

}

func (h *UserBaseHandler) GetOneComplaintHandler(c *gin.Context) {

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

  userID,err:=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  getOneComplaint,err:=FetchOneComplaint(h.DB,uint(userID),uint(getCompID))
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Fetch One Complaint Error : %s",err.Error()),err,nil,c)
    return
  }

  payloadData:=map[string]interface{}{
    "complaintData":getOneComplaint,
  }

  u.ResponseFormatter(http.StatusOK,"Pengaduan Detil Anda",nil,payloadData,c)
}

func (h *UserBaseHandler) UpdateComplaintHandler(c *gin.Context) {
  var validComplaint UpdateComplaintValidation

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

  userID,err:=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  _,err=UpdateComplaint(h.DB,uint(userID),uint(getCompID),validComplaint)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Update Complaint Error : %s",err.Error()),err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusOK,"Pengaduan Anda Berhasil Di Perbaharui",nil,nil,c)
}

func (h *UserBaseHandler) DeleteComplaintHandler(c *gin.Context) {

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

  userID,err:=u.FetchAuth(ad,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"User Not Authorized",err,nil,c)
    return
  }

  _,err=DeleteComplaint(h.DB,uint(userID),uint(getCompID))
  if err!=nil{
    u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Delete Complaint Error : %s",err.Error()),err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusNonAuthoritativeInfo,"Pengaduan Anda Berhasil DiHapus",nil,nil,c)

}