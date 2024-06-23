package dto

import (
	"time"

	"main/internal/user/model"
	"main/pkg/paging"
)

type KUser struct {
	ID                    string         `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             *time.Time     `json:"deleted_at" gorm:"index"`
	Password              string         `json:"password"`
	Role                  model.UserRole `json:"role"`
	Email                 string         `json:"email" gorm:"unique;not null;index:idx_user_email"`
	Name                  string         `json:"name"`
	PhoneNumber           string         `json:"phone_number"`
	VerifyCodeEmail       int            `json:"verify_code_email"`
	VerifyCodePhoneNumber int            `json:"verify_code_phone_number"`
	ApproveEmail          bool           `json:"approve_email"`
	ApprovePhoneNumber    bool           `json:"approve_phone_number"`
}
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterReq struct {
	Password    string         `json:"password" validate:"required,password"`
	Role        model.UserRole `json:"role"`
	Email       string         `json:"email" validate:"required,email"`
	Name        string         `json:"name"`
	PhoneNumber string         `json:"phone_number"`
}

type RegisterRes struct {
	User User `json:"user"`
}

type LoginReq struct {
	Email    string         `json:"email" validate:"required,email"`
	Password string         `json:"password" validate:"required,password"`
	Role     model.UserRole `json:"role"`
}

type LoginRes struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenRes struct {
	AccessToken string `json:"access_token"`
}

type UpdateUserReq struct {
	ID          string         `json:"id" gorm:"unique;not null;index;primary_key"`
	NewPassword string         `json:"new_password" validate:"required,password"`
	Password    string         `json:"password" validate:"required,password"`
	Role        model.UserRole `json:"role"`
	Email       string         `json:"email" validate:"required,email"`
	Name        string         `json:"name"`
	PhoneNumber string         `json:"phone_number"`
}

type UpdateUserRes struct {
	Message string `json:"message"`
}

//***************************************************************************\\
//***************************************************************************\\

type VerifyPhoneNumberRequest struct {
	PhoneNumber           string `json:"phone_number"`
	VerifyCodePhoneNumber string `json:"verify_code_phone_number"`
}
type ResendVerifyPhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// ***************************************************************************\\
type VerifyEmailRequest struct {
	Email           string `json:"email"`
	VerifyCodeEmail string `json:"verify_code_email"`
}
type ResendVerifyEmailRequest struct {
	Email string `json:"email"`
}

// ***************************************************************************\\
type VerifyResponse struct {
	Message string `json:"message"`
}

// ***************************************************************************\\
type ListUsersReq struct {
	// Name of the user
	// example: "Home"
	Name string `json:"name,omitempty" form:"name"`
	// User ID associated with the address
	// example: "67890"
	IDUser string `json:"id_user"`
	// Page number for pagination
	// example: 1
	Page int64 `json:"-" form:"page"`
	// Limit number of items per page
	// example: 10
	Limit int64 `json:"-" form:"limit"`
}
type ListUsersRes struct {
	// List of Users
	// example: [{"id_Users":"12345","id_user":"67890","name":"Home","city":"San Francisco","street":"Market Street","lat":"37.7749","long":"-122.4194"}]
	Users []*KUser `json:"Users"`
	// Pagination info
	Pagination *paging.Pagination `json:"pagination"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// DeleteAddressReq represents the request body for deleting an address.
// swagger:model DeleteAddressReq
type DeleteUserReq struct {
	// ID of the address
	// example: "12345"
	ID string `json:"id"`
}

//***************************************************************************\\
//***************************************************************************\\
