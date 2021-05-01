package auth

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	mUser "sipmas-api/src/apps/users"

	jwt "github.com/dgrijalva/jwt-go"
	hermes "github.com/matcornic/hermes/v2"
	goMail "gopkg.in/gomail.v2"

	"github.com/go-redis/redis/v8"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB,vSignupUser *SignupValidation) (mUser.UserModel,error) {

  //!Fill user & address struct with given valid signup user
  newUser:= mUser.UserModel {
    FullName: vSignupUser.FullName,
    Ktp:vSignupUser.Ktp,
    Email:vSignupUser.Email,
    Password: vSignupUser.Password,
    Age:vSignupUser.Age,
    Address: mUser.AddressModel {  
      FullAddress: vSignupUser.FullAddress,
      City:vSignupUser.City,
      SubDistrict: vSignupUser.SubDistrict,
      PostalCode:vSignupUser.PostalCode,  
    },
    Phone:vSignupUser.Phone,
  }

  //!Check if user is already exists in database or not
  if err:=db.First(&newUser, "email = ?", newUser.Email).Error;err!=nil{

    //!Create user itself 
    db.Create(&newUser)

    //!Return a user model 
    return newUser,nil

  }else{
    return mUser.UserModel{},err
  }

}

func GenerateConfirmToken(email string,rds *redis.Client,db *gorm.DB) (string,error) {
  
  var getUser mUser.UserModel

  //!Setup context 
  ctx:=context.Background()

  //!Make sure user is existed with given email & Check is user already verified or not
  if err:=db.Where("email= ? AND is_verified=?",email,false).First(&getUser).Error;err!=nil{
    return "",errors.New("akun anda sudah terverifikasi, silahkan untuk melanjutkan masuk")
  }

  //!Encode email into a hex string
  encodedHashMail:=hex.EncodeToString([]byte(email))

  //!Store encoded email into redis, expired in 5 mins.
  if err := rds.Set(ctx, fmt.Sprintf("tokenForEncodedEmail:%s",email), encodedHashMail, time.Minute * 5).Err(); err!=nil {
    return "",err
  }

  return encodedHashMail,nil
}

func SendConfirmEmail(token string,validUser *mUser.UserModel,rds *redis.Client,db *gorm.DB) error {

  //!Decode & Verify Token 
  decodedMail,err:=decodeAndVerifyToken(token,rds,db)
  if err!=nil{
    return err
  }

  //!Setup GoMail Here 
  m := goMail.NewMessage()

  m.SetHeader("From", os.Getenv("SMTP_USER"))
  m.SetHeader("To",decodedMail)
  m.SetAddressHeader("Cc",os.Getenv("SMTP_USER"), "Administrator")
  m.SetHeader("Subject", "Konfirmasi Akun SIPMAS")

  //!Setup Hermes here
  //?Configure hermes by setting a theme and your product info
  h := hermes.Hermes{
    Product: hermes.Product{
      //?Appears in header & footer of e-mails
      Name: "SIPMAS",
      Link: "#",
      //?Optional product logo & Copyright
      Logo: "https://iili.io/qQB8yx.png",
      Copyright: fmt.Sprintf("Copyright © %d SIPMAS. All rights reserved.",time.Now().Year()),
    },
  }

  setupEmail := hermes.Email{
    Body:hermes.Body{
      Name:validUser.FullName,
      Greeting: "Hai",
      Signature: "Hormat Kami",
      Intros: []string{
        "Sebelum kamu melanjutkan masuk, pastikan untuk konfirmasi akunmu terlebih dahulu ya!",
      },
      Actions: []hermes.Action{
        {
          Instructions: "Silahkan untuk klik tombol dibawah ini :",
          Button: hermes.Button{
              Color: "#22BC66",
              Text:  "Konfirmasi Akun",
              Link:  fmt.Sprintf("http://localhost:35401/api/v1/auth/konfirmasi?token=%s",token),
          },
        },
      },
      Outros: []string{
        "Butuh pertanyaan atau bantuan ? cukup balas email ini. Kami siap membantu dan melayani anda.",
      },
    },
  }

  //!Parse setupMail into pure HTML
  emailBody, err := h.GenerateHTML(setupEmail)
  if err != nil {
    return err
  }

  //!Set body Mail with HTML email
  m.SetBody("text/html",emailBody)

  //!Setup Dialer 
  d := goMail.NewDialer(os.Getenv("SMTP_HOST"),587,os.Getenv("SMTP_USER"),os.Getenv("SMTP_PASS"))

  //! Send the email
  if err := d.DialAndSend(m); err != nil {
    return err
  }

  return nil
  
}

func ConfirmAccToken(token string,db *gorm.DB,rds *redis.Client) error {  

  var confirmUser mUser.UserModel

  //!Decode & Verify Token 
  decodedMail,err:=decodeAndVerifyToken(token,rds,db)
  if err!=nil{
    return err
  }

  //!Update user field is_verified from false to true 
  db.Model(&confirmUser).Where("email = ?", decodedMail).Update("is_verified", true)
  return nil

}

func VerifyUser(db *gorm.DB,validUser *SigninValidation) (mUser.UserModel,error){

  var getUser mUser.UserModel

  //!Find user exists or not and verified
  if err:=db.Where("email = ? AND is_verified = ?",&validUser.Email,true).First(&getUser).Error;err!=nil{

    //!If not verified 
    if !getUser.IsVerified {
      return mUser.UserModel{},err
    }
    return mUser.UserModel{},err 

  }

  //!Check the password is correct or not 
  if err:=bcrypt.CompareHashAndPassword([]byte(getUser.Password),[]byte(validUser.Password));err!=nil{
    return mUser.UserModel{},errors.New("invalid password") 
  }

  //!Update user field last_login_at to now 
  getUser.LastLoginAt=time.Now()
  db.Save(&getUser)

  return getUser,nil
  
}

func CreateJWT(userId uint64) (*TokenDetails,error) {

  var err error
  
  //!Create Token Details 
  td := &TokenDetails{}

  //?Access token expires in 5 min 
  td.AtExpires = time.Now().Add(time.Minute * 5).Unix()
  //?Generate uuid for each unique access token 
  td.AccessUuid = uuid.NewV4().String()

  //?Refresh token expires in 1 day  
  td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
  //?Generate uuid for each unique refresh token 
  td.RefreshUuid = uuid.NewV4().String()


  //!Create Access Token Claims
  atClaims := jwt.MapClaims{}
  atClaims["authorized"] = true
  atClaims["access_uuid"] = td.AccessUuid
  atClaims["user_id"] = userId
  atClaims["exp"] = td.AtExpires

  //?New jwt with signed method + refresh token claims 
  at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

  //?Signed with secret 
  td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACC_JWT_SECRET")))
  if err != nil {
    return nil, err
  }

  //!Create Refresh Token Claims
  rtClaims:=jwt.MapClaims{}
  rtClaims["refresh_uuid"] = td.RefreshUuid
  rtClaims["user_id"] = userId
  rtClaims["exp"] = td.RtExpires

  //?New jwt with signed method + access token claims 
  rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

  //?Signed with secret 
  td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REF_JWT_SECRET")))
  if err != nil {
    return nil, err
  }

  return td,nil

}

func CreateAuth(userId uint64,td *TokenDetails,rds *redis.Client) error {

  //!Setup Context 
  ctx:=context.Background()

  //!Converting token expires to unix time 
  at:=time.Unix(td.AtExpires,0)
  rt:=time.Unix(td.RtExpires,0)
  now:=time.Now()

  //!Set access uuid as a key & userId as a value into redis 
  errAccess:=rds.Set(ctx,td.AccessUuid,userId,at.Sub(now)).Err()
  if errAccess!=nil{
    return errAccess
  }
  //!Set refresh uuid as a key & userId as a value into redis 
  errRefresh:=rds.Set(ctx,td.RefreshUuid,userId,rt.Sub(now)).Err()
  if errAccess!=nil{
    return errRefresh
  }

  return nil
}

/*UNEXPORTED FUNC SECTION*/ 

func decodeAndVerifyToken(token string,rds *redis.Client,db *gorm.DB) (string,error){

  var getUser mUser.UserModel

  //!Setup context 
  ctx:=context.Background()

  //!Decode hex string into raw email string
  decodedMail,err:=hex.DecodeString(token)
  if err!=nil {
    return "",err
  }

  //!Check the email is associated with specific user in database or not
  if err:=db.Where("email = ?",string(decodedMail)).First(&getUser).Error;err!=nil{
    return "",err
  } 

  //!Checking the value in redis is exists or not 
  if err:= rds.Get(ctx,fmt.Sprintf("tokenForEncodedEmail:%s",string(decodedMail))).Err();err!=nil{
    
    //?If nil, re-create confirm token & send email again to user
    newToken,err:=GenerateConfirmToken(string(decodedMail),rds,db)
    if err!=nil{
      return "",err
    }
    err = SendConfirmEmail(newToken,&getUser,rds,db)
    if err!=nil{
      return "",err
    }

    return "",errors.New("token anda sudah kadaluarsa, silahkan untuk cek email konfirmasi ulang dari kami kembali")

  }

  return string(decodedMail),nil
}
