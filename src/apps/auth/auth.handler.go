package auth

import (
	"fmt"
	"net/http"
  "strconv"
  "os"
	u "sipmas-api/src/utils"

  jwt "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

func (h *AuthBaseHandler) SignupHandler(c *gin.Context) {
  var validUser SignupValidation

  if err:=c.ShouldBindJSON(&validUser); err!= nil {
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }

  newUser,err:=CreateUser(h.postDB,&validUser)
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Database Error : %s",err.Error()),err,nil,c)
    return
  }
  
  token,err:=GenerateConfirmToken(newUser.Email,h.rdsDB,h.postDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Generate Confirm Token Error : %s",err.Error()),err,nil,c)
    return
  }

  err = SendConfirmEmail(token,&newUser,h.rdsDB,h.postDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Send Email Error : %s",err.Error()),err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusOK,fmt.Sprintf("Konfirmasi email kami sudah terkirim ke akun %s. Mohon untuk cek terlebih dahulu.",newUser.Email),nil,nil,c)

}

func (h *AuthBaseHandler) ConfirmAccountHandler(c *gin.Context) {

  getToken,ok:=c.GetQuery("token")
  if !ok {
    u.ResponseFormatter(http.StatusNotFound,"Token not found.",nil,nil,c)
    return
  }

  if err:=ConfirmAccToken(getToken,h.postDB,h.rdsDB);err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Confirm Account Token Error : %s",err.Error()),err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusOK,"Akun Kamu Berhasil di Verifikasi. Silahkan untuk masuk menggunakan akun tersebut.",nil,nil,c)

}

func (h *AuthBaseHandler) SigninHandler(c *gin.Context) {
  var validUser SigninValidation

  if err:=c.ShouldBindJSON(&validUser); err!= nil {
    u.ResponseFormatter(http.StatusBadRequest,fmt.Sprintf("Validation Error : %s",err.Error()),err,nil,c)
    return
  }

  getUser,err:=VerifyUser(h.postDB,&validUser)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,fmt.Sprintf("Verify User Error : %s",err.Error()),err,nil,c)
    return
  }

  tDetail,err:=CreateJWT(uint64(getUser.ID))
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Create JWT Token Error : %s",err.Error()),err,nil,c)
    return
  }

  err = CreateAuth(uint64(getUser.ID),tDetail,h.rdsDB)
  if err!=nil{
    u.ResponseFormatter(http.StatusInternalServerError,fmt.Sprintf("Create Auth Error : %s",err.Error()),err,nil,c)
    return
  }

  tokens:=map[string]interface{}{
    "accessToken":tDetail.AccessToken,
    "refreshToken":tDetail.RefreshToken,
  }

  u.ResponseFormatter(http.StatusOK,"Berhasil Masuk.",nil,tokens,c)

}

func (h *AuthBaseHandler) TokenHandler(c *gin.Context) {

  mapToken := map[string]string{}

  if err := c.ShouldBindJSON(&mapToken); err != nil {
    u.ResponseFormatter(http.StatusBadRequest,"Refresh Token Not Found",err,nil,c)
    return
  }

  refreshToken := mapToken["refreshToken"]

  token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
  
    //!Make sure that the token method conform to "SigningMethodHMAC"
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(os.Getenv("REF_JWT_SECRET")), nil
  })

  //!if there is an error, the token must have expired
  if err != nil {
    u.ResponseFormatter(http.StatusUnauthorized,"Refresh Token Already Expired",err,nil,c)
    return
  }

  //!is token valid?
  if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
    u.ResponseFormatter(http.StatusUnauthorized,"Refresh Token Invalid",err,nil,c)
    return
  }

  //!Since token is valid, get the uuid:
  claims, ok := token.Claims.(jwt.MapClaims) //?the token claims should conform to MapClaims

  if ok && token.Valid {
    refreshUuid, ok := claims["refresh_uuid"].(string) //?convert the interface to string
    if !ok {
      u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Claiming Refresh token Error : %s",err.Error()),err,nil,c)
      return
    }
    userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
    if err != nil {
      u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Parsing userId Error : %s",err.Error()),err,nil,c)
      return
    }

    //!Delete the previous Refresh Token
    deleted,err:=u.DeleteAuth(refreshUuid,h.rdsDB)
    if err!=nil||deleted==0{
      u.ResponseFormatter(http.StatusUnprocessableEntity,"Auth Already Deleted",err,nil,c)
      return
    }

    //!Create new jwt with current userId 
    tDetail,err:=CreateJWT(userId)
    if err!=nil{
      u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Create JWT Token Error : %s",err.Error()),err,nil,c)
      return
    }

    //!Create new auth
    err = CreateAuth(userId,tDetail,h.rdsDB)
    if err!=nil{
      u.ResponseFormatter(http.StatusUnprocessableEntity,fmt.Sprintf("Create Auth Error : %s",err.Error()),err,nil,c)
      return
    }

    //!Create new map tokens 
    tokens:=map[string]interface{}{
      "accessToken":tDetail.AccessToken,
      "refreshToken":tDetail.RefreshToken,
    }

    //!Return the token 
    u.ResponseFormatter(http.StatusCreated,"Token Berhasil di Refresh",nil,tokens,c)

  }

}

func (h *AuthBaseHandler) SignoutHandler(c *gin.Context) {
  ad,err:=u.ExtractTokenMetadata(c.Request)
  if err!=nil{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",err,nil,c)
    return
  }

  deleted,err:=u.DeleteAuth(ad.AccessUuid,h.rdsDB)
  if err!=nil||deleted==0{
    u.ResponseFormatter(http.StatusUnauthorized,"Token anda sudah kadaluarsa, silahkan untuk mendapatkannya disini http://localhost:35401/api/v1/auth/token/refresh",err,nil,c)
    return
  }

  u.ResponseFormatter(http.StatusOK,"Berhasil keluar.",nil,nil,c)
}

