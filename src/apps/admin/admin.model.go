package admin

type ManageUserValidation struct {
  IsActive bool   `json:"isActive"`
}

type ManageComplaintValidation struct {
  Status  string   `json:"status"`
}