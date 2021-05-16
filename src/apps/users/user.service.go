package users

import (
  "errors"
  "reflect"
  
  u "sipmas-api/src/utils"

  "gorm.io/gorm"
)

func CountComplaintUser(UserID uint,db *gorm.DB) ([]int,error) {

  //!Setup variable for holding 3 result count of complaint with 3 condition 
  var rejectCount,acceptCount,defaultCount int64

  if err:=db.Model(&ComplaintModel{}).Where("user_id=? AND status=?",UserID,"PROSES").Count(&defaultCount).Error;err!=nil{
    return nil,err
  }
  if err:=db.Model(&ComplaintModel{}).Where("user_id=? AND status=?",UserID,"DITERIMA").Count(&acceptCount).Error;err!=nil{
    return nil,err
  }
  if err:=db.Model(&ComplaintModel{}).Where("user_id=? AND status=?",UserID,"DITOLAK").Count(&rejectCount).Error;err!=nil{
    return nil,err
  }

  //!Returning array of int 
  return []int{int(defaultCount),int(acceptCount),int(rejectCount)},nil

}

func FetchUserProfile(UserID uint,db *gorm.DB) (UserModel,error) {
  var getUser UserModel

  //!Find user by user id and user roles with joining address table
  if err:=db.Joins("Address").First(&getUser,"user_models.id=? AND user_models.roles=?",UserID,"USER").Error;err!=nil{
    return UserModel{},errors.New("user not found / admin cannot retrieve privacy data from user")
  }

  return getUser,nil
}

func EditProfile(UserID uint,db *gorm.DB,validUser UpdateUserValidation) (UserModel,error) {

  //!Returning a reflect value (from a user model validation struct)
  v :=reflect.ValueOf(validUser)
  typeOfV:=v.Type()

  //!Create blank map 
  inputData:=map[string]interface{}{}

  //!Loop through number of fields of struct 
  for i:=0; i < v.NumField(); i++{

    //!fill the map with field name from reflect val , current value with interface{}
    inputData[typeOfV.Field(i).Name]=v.Field(i).Interface()

    //!If the fields nil , delete that fields
    if inputData[typeOfV.Field(i).Name]== "" {
      delete(inputData,typeOfV.Field(i).Name)
    }

    //!Check if inputData has property Age or not 
    if val,ok:=inputData["Age"];ok{
      //!Check if that property is have a zero value or not
      if u.IsZeroOfUnderlyingType(val){
        //!If yes, delete that property and value from inputData 
        delete(inputData,"Age")
      }
    }
  }

  //!Send the prev map as updated fields, along with condition where id user & role "USER" is existed 
  result:=db.Model(&UserModel{}).Where("id=? AND roles=?",UserID,"USER").Updates(inputData)
  
  //!If the no rows affected / 0, return error
  if result.RowsAffected == 0 {
    return UserModel{},errors.New("tidak ada yang diupdate")
  }

  return UserModel{},nil

}

func CreateComplaint(db *gorm.DB,validUser UserModel,validComplaint *ComplaintValidation)(ComplaintModel,error) {

  //!Fill new struct of complaint model
  newComplaint:=ComplaintModel{
    Title: validComplaint.Title,
    Content:validComplaint.Content,
    User: validUser,
  }

  //!Create new complaint 
  if err:=db.Create(&newComplaint).Error;err!=nil{
    return ComplaintModel{},err
  }

  return newComplaint,nil

}

func UpdateImageComplaint(db *gorm.DB,CompID uint,fileName string) (string,error) {

  //!Update complaint field image_name with given fileName in parameters 
  if err:=db.Model(&ComplaintModel{}).Where("id=?",CompID).Update("image_name",fileName).Error;err!=nil{
    return "",err
  }

  return fileName,nil
}

func FetchAllComplaint(db *gorm.DB,UserID uint) ([]ComplaintModel,error) {

  var getAllComplaint []ComplaintModel

  //!Find all complaint based on user_id with joins user table & address table 
  if err:=db.Model(&ComplaintModel{}).Joins("User").Preload("User.Address").Find(&getAllComplaint,"user_id=?",UserID).Error;err!=nil{
    return []ComplaintModel{},err
  }

  return getAllComplaint,nil

}

func FetchOneComplaint(db *gorm.DB,UserID uint,CompID uint) (ComplaintModel,error) {
  
  var getOneComplaint ComplaintModel

  //!Find specific complaint based on user_id & complaint id. Joins with user table & address table 
  if err:=db.Model(&ComplaintModel{}).Joins("User").Preload("User.Address").First(&getOneComplaint,"user_id=? AND complaint_models.id=?",UserID,CompID).Error;err!=nil{
    return ComplaintModel{},err
  }

  return getOneComplaint,nil

}

func UpdateComplaint(db *gorm.DB,UserID uint,CompID uint,validComplaint UpdateComplaintValidation) (ComplaintModel,error) {
  
  //!Returning a reflect value (from a complaint model validation struct)
  v :=reflect.ValueOf(validComplaint)
  typeOfV:=v.Type()

  //!Create blank map 
  inputData:=map[string]interface{}{}

  //!Loop through number of fields of struct 
  for i:=0; i < v.NumField(); i++{

    //!fill the map with field name from reflect val , current value with interface{}
    inputData[typeOfV.Field(i).Name]=v.Field(i).Interface()

    //!If the fields nil , delete that fields
    if inputData[typeOfV.Field(i).Name]== "" {
      delete(inputData,typeOfV.Field(i).Name)
    }

  }

  //!Send the prev map as updated fields, along with condition where user_id & complaint id is existed 
  result:=db.Model(&ComplaintModel{}).Joins("User").Where("user_id=? AND complaint_models.id=?",UserID,CompID).Updates(inputData)
  
  //!If the no rows affected / 0, return error
  if result.RowsAffected == 0 {
    return ComplaintModel{},errors.New("tidak ada yang diupdate")
  }

  return ComplaintModel{},nil

}

func DeleteComplaint(db *gorm.DB,UserID uint,CompID uint) (ComplaintModel,error) {

  var getComplaint ComplaintModel

  //!Delete a complaint with condition where user_id=? & complaint id is existed 
  result:=db.Model(&ComplaintModel{}).Joins("User").Where("user_id=? AND complaint_models.id=?",UserID,CompID).Delete(&getComplaint)

  //!If the no rows affected / 0, return error
  if result.RowsAffected == 0 {
    return ComplaintModel{},errors.New("tidak ada yang dihapus")
  }

  return getComplaint,nil
}