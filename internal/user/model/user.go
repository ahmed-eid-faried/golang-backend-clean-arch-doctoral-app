package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"main/pkg/utils"
)

// UserRole represents the role of a user
type UserRole string

// Constants for user roles
const (
	UserRoleAdmin  UserRole = "admin"  // Administrator role
	UserRoleDoctor UserRole = "doctor" // Doctor role
	UserRoleClient UserRole = "client" // Client role
)

// User represents a user in the system
type User struct {
	ID                    string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at" gorm:"index"`
	Password              string     `json:"password"`
	Role                  UserRole   `json:"role"`
	Email                 string     `json:"email" gorm:"unique;not null;index:idx_user_email"`
	Name                  string     `json:"name"`
	PhoneNumber           string     `json:"phone_number"`
	VerifyCodeEmail       int        `json:"verify_code_email"`
	VerifyCodePhoneNumber int        `json:"verify_code_phone_number"`
	ApproveEmail          bool       `json:"approve_email"`
	ApprovePhoneNumber    bool       `json:"approve_phone_number"`
}

// BeforeCreate is a hook that is called before creating a new user
// func (user *User) BeforeCreate(tx *gorm.DB) error {
func (user *User) BeforeCreate() error {
	// Generate a unique ID for the user
	user.ID = uuid.New().String()

	// Hash and salt the password for security
	user.Password = utils.HashAndSalt([]byte(user.Password))

	// Set the default role to customer if not specified
	if user.Role == "" {
		user.Role = UserRoleClient
	}

	// Generate verification codes
	user.VerifyCodePhoneNumber = utils.GenerateRandomCode()
	user.VerifyCodeEmail = utils.GenerateRandomCode()
	user.ApproveEmail = false
	user.ApprovePhoneNumber = false

	// Send verification codes
	err := SendVerificationCodes(user)
	if err != nil {
		return err
	}
	return nil
}

// BeforeCreate is a hook that is called before creating a new user
func (user *User) BeforeUpdate() error {
	// Generate verification codes
	user.VerifyCodeEmail = utils.GenerateRandomCode()
	user.VerifyCodePhoneNumber = utils.GenerateRandomCode()
	// Send verification codes
	err := SendVerificationCodes(user)
	if err != nil {
		return err
	}
	return nil
}

// BeforeCreate is a hook that is called before creating a new user
func (user *User) BeforeUpdateVerificationEmail() error {
	// Generate verification codes
	user.VerifyCodeEmail = utils.GenerateRandomCode()
	// Send verification codes
	err := SendVerificationCodesEmail(user)
	if err != nil {
		return err
	}
	return nil
}

// BeforeCreate is a hook that is called before creating a new user
func (user *User) BeforeUpdateVerificationPhone() error {
	// Generate verification codes
	user.VerifyCodePhoneNumber = utils.GenerateRandomCode()
	// Send verification codes
	err := SendVerificationCodesPhone(user)
	if err != nil {
		return err
	}
	return nil
}

// sendVerificationCodes sends the verification codes to the user's phone and email
func SendVerificationCodes(user *User) error {
	// Send email verification code
	err := SendVerificationCodesEmail(user)
	if err != nil {
		return err
	}

	// Send phone verification code
	err = SendVerificationCodesPhone(user)
	if err != nil {
		return err
	}

	return nil
}

// Send phone verification code
func SendVerificationCodesPhone(user *User) error {

	phoneMessage := fmt.Sprintf("Your verification code is %d", user.VerifyCodePhoneNumber)
	err := utils.SendSMS(user.PhoneNumber, phoneMessage)
	if err != nil {
		return fmt.Errorf("failed to send SMS verification code: %w", err)
	}

	return nil
}

// Send email verification code
func SendVerificationCodesEmail(user *User) error {
	emailSubject := "Your Email Verification Code"
	emailPlainText := fmt.Sprintf("Your verification code is %d", user.VerifyCodeEmail)
	emailHTMLContent := fmt.Sprintf("<strong>Your verification code is %d</strong>", user.VerifyCodeEmail)
	err := utils.SendEmail(user.Name, user.Email, emailSubject, emailPlainText, emailHTMLContent)
	if err != nil {
		return fmt.Errorf("failed to send email verification code: %w", err)
	}

	return nil
}
