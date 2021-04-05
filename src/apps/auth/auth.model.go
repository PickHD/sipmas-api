package auth

type SignupValidation struct {
	FullName 		string	`json:"fullName" binding:"required"`
	Ktp					string	`json:"ktp" binding:"required"`
	Email		  	string	`json:"email" binding:"required"`
	Password		string	`json:"password" binding:"required"`
	Age					int			`json:"age" binding:"required"`
	FullAddress	string	`json:"fullAddress" binding:"required"`
	City 				string	`json:"city" binding:"required"`
	SubDistrict	string	`json:"subDistrict" binding:"required"`
	PostalCode	string	`json:"postalCode" binding:"required"`
	Phone				string	`json:"phone" binding:"required"`
}

type SigninValidation struct {
	Ktp				string	`json:"ktp" binding:"required"`
	Password	string	`json:"password" binding:"required"`
}