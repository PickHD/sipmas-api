package admin

import (
  "reflect"
  "errors"

	user "sipmas-api/src/apps/users"

	"gorm.io/gorm"
)

func CountAllUserAndComplaint(db *gorm.DB) (map[string]interface{},error) {
  //TODO : 1. Count all user active and not active separately, execpt user with ADMIN roles. 2. count all complaint with 3 condition (PROSES,DITOLAK,DITERIMA) separately. 3. return all with map[string]interface{}

  var (
    activeUserCount       int64
    inactiveUserCount     int64
    processComplaintCount int64
    acceptComplaintCount  int64
    rejectComplaintCount  int64
  )

  if err:=db.Model(&user.UserModel{}).Where("is_active=? AND roles=?",true,"USER").Count(&activeUserCount).Error;err!=nil{
    return nil,err
  }
  if err:=db.Model(&user.UserModel{}).Where("is_active=? AND roles=?",false,"USER").Count(&inactiveUserCount).Error;err!=nil{
    return nil,err
  }
  if err:=db.Model(&user.ComplaintModel{}).Where("status=?","PROSES").Count(&processComplaintCount).Error;err!=nil{
    return nil,err
  }
  if err:=db.Model(&user.ComplaintModel{}).Where("status=?","DITERIMA").Count(&acceptComplaintCount).Error;err!=nil{
    return nil,err
  }
  if err:=db.Model(&user.ComplaintModel{}).Where("status=?","DITOLAK").Count(&rejectComplaintCount).Error;err!=nil{
    return nil,err
  }

  return map[string]interface{}{
    "statusUserAndComplaint":map[string]int64{
      "totalActiveUser":activeUserCount,
      "totalInactiveUser":inactiveUserCount,
      "totalProcessComplaint":processComplaintCount,
      "totalAcceptComplaint":acceptComplaintCount,
      "totalRejectComplaint":rejectComplaintCount,
    },
    "manageLink":"http://localhost:35401/api/v1/admin/kelola",
  },nil

}

func FetchAllComplaint(db *gorm.DB) ([]user.ComplaintModel,error) {
  //TODO : 1. Get all complaint with preload users table & address table, return a nil if not found. 2. returning a array of complaint.

  var getAllComplaint []user.ComplaintModel

  //!Find all complaint with joins user table & address table 
  if err:=db.Model(&user.ComplaintModel{}).Joins("User").Preload("User.Address").Find(&getAllComplaint).Error;err!=nil{
    return []user.ComplaintModel{},err
  }

  return getAllComplaint,nil
}

func FetchOneComplaint(db *gorm.DB,CompID uint) (user.ComplaintModel,error) {
  //TODO : 1. Get spec. one complaint (where id and user_id) with preload users table & address table, return error if not found. 2. returning a complaint.
  
  var getOneComplaint user.ComplaintModel

  //!Find one spec. complaint with joins user table & address table 
  if err:=db.Model(&user.ComplaintModel{}).Joins("User").Preload("User.Address").First(&getOneComplaint,"complaint_models.id=?",CompID).Error;err!=nil{
    return user.ComplaintModel{},err
  }

  return getOneComplaint,nil
}

func UpdateComplaint(db *gorm.DB,CompID uint,validComplaint ManageComplaintValidation) (user.ComplaintModel,error) {
  //TODO : 1. Get spec. one complaint (where id and user_id), return error if not found. 2. update complaint (in this case, ADMIN only can update status) and if the rows is not affected, return error, 3. returning a complaint.

  //!Returning a reflect value (from a complaint model validation struct)
  v :=reflect.ValueOf(validComplaint)
  typeOfV:=v.Type()

  //!Create blank map 
  inputData:=map[string]interface{}{}

  //!Loop through number of fields of struct 
  for i:=0; i < v.NumField(); i++{

    //!fill the map with field name from reflect val, current value with interface{}
    inputData[typeOfV.Field(i).Name]=v.Field(i).Interface()

    //!If the fields nil , delete that fields
    if inputData[typeOfV.Field(i).Name]== "" {
      delete(inputData,typeOfV.Field(i).Name)
    }

    if inputData["Status"] == "" { 
      inputData["Status"]="PROSES"
    } 

  }

  //!Send the prev map as updated fields, along with condition where user_id is existed 
  result:=db.Model(&user.ComplaintModel{}).Where("complaint_models.id=?",CompID).Updates(inputData)
  
  //!If the no rows affected / 0, return error
  if result.RowsAffected == 0 {
    return user.ComplaintModel{},errors.New("tidak ada yang diupdate")
  }

  return user.ComplaintModel{},nil
}

func FetchAllUser(db *gorm.DB) ([]user.UserModel,error) {
  //TODO : 1. Get all user (where user is not ADMIN and user is active) with preload address table, return a nil if not found. 2. returning a array of user.

  var getAllUser []user.UserModel

  //!Find all users with joins address table 
  if err:=db.Model(&user.UserModel{}).Joins("Address").Find(&getAllUser,"roles=? AND is_active=?","USER",true).Error;err!=nil{
    return []user.UserModel{},err
  }

  return getAllUser,nil
}

func FetchOneUser(db *gorm.DB,UserID uint) (user.UserModel,error) {
  //TODO : 1. Get spec. one user (where id or user is not ADMIN or user is active) with preload address table, return error if not found. 2. returning a user.

  var getOneUser user.UserModel

  //!Find all users with joins address table 
  if err:=db.Model(&user.UserModel{}).Joins("Address").First(&getOneUser,"roles=? AND is_active=?","USER",true).Error;err!=nil{
    return user.UserModel{},err
  }

  return getOneUser,nil
}

func UpdateUser(db *gorm.DB,UserID uint,validUser ManageUserValidation) (user.UserModel,error) {
  //TODO : 1. Get spec. one user (where id or user is not ADMIN or user is active), return error if not found. 2. update user (in this case, ADMIN only can update isActive) and if the rows is not affected, return error, 3. returning a user.

  //!Returning a reflect value (from a user model validation struct)
  v :=reflect.ValueOf(validUser)
  typeOfV:=v.Type()

  //!Create blank map 
  inputData:=map[string]interface{}{}

  //!Loop through number of fields of struct 
  for i:=0; i < v.NumField(); i++{

    //!fill the map with field name from reflect val, current value with interface{}
    inputData[typeOfV.Field(i).Name]=v.Field(i).Interface()

    //!If the fields nil , delete that fields
    if inputData[typeOfV.Field(i).Name]== "" {
      delete(inputData,typeOfV.Field(i).Name)
    }

    if inputData["IsActive"] == "" { 
      inputData["IsActive"]=true
    } 

  }

  //!Send the prev map as updated fields, along with condition where complaint id is existed 
  result:=db.Model(&user.ComplaintModel{}).Where("user_models.id=?",UserID).Updates(inputData)
  
  //!If the no rows affected / 0, return error
  if result.RowsAffected == 0 {
    return user.UserModel{},errors.New("tidak ada yang diupdate")
  }

  return user.UserModel{},nil
}

