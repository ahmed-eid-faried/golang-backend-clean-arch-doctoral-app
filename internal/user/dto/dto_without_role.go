package dto

type KUpdateUserReq struct {
	ID          string `json:"id" gorm:"unique;not null;index;primary_key"`
	NewPassword string `json:"new_password" validate:"required,password"`
	Password    string `json:"password" validate:"required,password"`
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}
type KLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}
type KRegisterReq struct {
	Password    string `json:"password" validate:"required,password"`
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}
