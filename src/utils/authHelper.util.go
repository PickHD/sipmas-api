package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
)

type AccessDetails struct {
  AccessUuid string
  UserId   uint64
}

func ExtractJWT(r *http.Request) string {

  //!Extract bearer token from header 
  bearToken:=r.Header.Get("Authorization")

  //!Split bearer token & take only the token 
  strArr:=strings.Split(bearToken," ")

  //!Check if splitted array length is 2 
  if len(strArr) == 2 {
    //!Returning splitted array index 1
    return strArr[1]
  }
  
  return ""
}

func VerifyJWT(r *http.Request) (*jwt.Token,error) {

  //!Call service to extract jwt 
  tokenStr:=ExtractJWT(r)

  //!Check if the token is empty string or not 
  if tokenStr == ""{
    return nil,errors.New("token is empty")
  }

  //!Parsing the token string into jwt token struct 
  token, err := jwt.Parse(tokenStr,func(token *jwt.Token) (interface{}, error){

    //?Check the signing method 
    if _,ok:=token.Method.(*jwt.SigningMethodHMAC);!ok{
      return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    //?Returning byte of access jwt secret
		return []byte(os.Getenv("ACC_JWT_SECRET")), nil
	})

  if err!=nil{
    return nil,err
  }

  return token,nil
}

func ValidateJWT(r *http.Request) error {
  
  //!Call service to verify jwt
  token,err:=VerifyJWT(r)

  if err!=nil{
    return err
  }

  //!Check the token is valid or not
  if _,ok:=token.Claims.(jwt.Claims);!ok && !token.Valid{
    return err
  }
  return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
  
  //!Call service to verify jwt
  token, err := VerifyJWT(r)
  if err != nil {
     return nil, err
  }
  //!Extract token claims
  claims, ok := token.Claims.(jwt.MapClaims)

  //?Check the token claims is existed and valid or not 
  if ok && token.Valid {
    
    //?Extract accessUuid from token claims
    accessUuid, ok := claims["access_uuid"].(string)
     
    if !ok {
      return nil, err
    }

    //?Convert userId from claims into a unsigned integer 
    userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
    if err != nil {
      return nil, err
    }

    //?Return a accessDetail (accessUuid,userId)   
    return &AccessDetails{
      AccessUuid: accessUuid,
      UserId:userId,
    }, nil
  }
  return nil, err
}

func FetchAuth(authD *AccessDetails,rds *redis.Client) (uint64, error) {
  
  //!Setup context 
  ctx:=context.Background()

  //!Get the value userId with key accessUuid from AccessDetails
  userid, err := rds.Get(ctx,authD.AccessUuid).Result()
  if err != nil {
    return 0, err
  }

  //!Convert userid into a unsigned integer 
  userID, _ := strconv.ParseUint(userid, 10, 64)
  return userID, nil
}

func DeleteAuth(givenUuid string,rds *redis.Client) (int64,error) {

  //!Setup context 
  ctx:=context.Background()

  //!Delete the auth with given uuid as a key 
  deleted, err := rds.Del(ctx,givenUuid).Result()
  if err != nil {
    return 0, err
  }
  return deleted, nil
}