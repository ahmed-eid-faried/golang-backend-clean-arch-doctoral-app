package service

import (
	"context"
	"errors"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	googleOauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"main/internal/user/dto"
	"main/internal/user/model"
	"main/internal/user/repository"
	"main/pkg/jtoken"
	"main/pkg/paging"
	"main/pkg/utils"
)

//go:generate mockery --name=IUserService
type IUserService interface {
	Login(ctx context.Context, req *dto.LoginReq) (*model.User, string, string, error)
	Register(ctx context.Context, req *dto.RegisterReq) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	RefreshToken(ctx context.Context, userID string) (string, error)
	VerifyEmail(ctx context.Context, request dto.VerifyEmailRequest) (dto.VerifyResponse, error)
	VerifyPhoneNumber(ctx context.Context, request dto.VerifyPhoneNumberRequest) (dto.VerifyResponse, error)
	ResendVerfiyCodePhone(ctx context.Context, request dto.ResendVerifyPhoneNumberRequest) (dto.VerifyResponse, error)
	ResendVerfiyCodeEmail(ctx context.Context, request dto.ResendVerifyEmailRequest) (dto.VerifyResponse, error)
	ListUsers(ctx context.Context, request dto.ListUsersReq) ([]*model.User, *paging.Pagination, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserReq) error
	Delete(ctx context.Context, id string, req *dto.DeleteUserReq) (*model.User, error)
	LoginWithGoogle(ctx context.Context, code string) (*model.User, string, string, error)
	LoginWithFacebook(ctx context.Context, code string) (*model.User, string, string, error)
}

type UserService struct {
	validator     validation.Validation
	repo          repository.IUserRepository
	oauthConfig   *oauth2.Config
	fbOauthConfig *oauth2.Config
}

func NewUserService(
	validator validation.Validation,
	oauthConfig *oauth2.Config,
	fbOauthConfig *oauth2.Config, repo repository.IUserRepository) *UserService {

	return &UserService{
		validator:     validator,
		repo:          repo,
		oauthConfig:   oauthConfig,
		fbOauthConfig: fbOauthConfig,
	}
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginReq) (*model.User, string, string, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, "", "", err
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Errorf("Login.GetUserByEmail fail, email: %s, error: %s", req.Email, err)
		return nil, "", "", err
	}

	if !(req.Role == user.Role) {
		return nil, "", "", errors.New("wrong Role")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", "", errors.New("wrong password")
	}

	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	accessToken := jtoken.GenerateAccessToken(tokenData)
	refreshToken := jtoken.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterReq) (*model.User, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	var user model.User
	utils.Copy(&user, &req)
	user.BeforeCreate()
	err := s.repo.Create(ctx, &user)
	if err != nil {
		logger.Errorf("Register.Create fail, email: %s, error: %s", req.Email, err)
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("GetUserByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) RefreshToken(ctx context.Context, userID string) (string, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		logger.Errorf("RefreshToken.GetUserByID fail, id: %s, error: %s", userID, err)
		return "", err
	}

	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	accessToken := jtoken.GenerateAccessToken(tokenData)
	return accessToken, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserReq) error {
	if err := s.validator.ValidateStruct(req); err != nil {
		return err
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("UpdateUser.GetUserByID fail, id: %s, error: %s", id, err)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return errors.New("wrong password")
	}

	user.Password = utils.HashAndSalt([]byte(req.NewPassword))
	err = s.repo.Update(ctx, user)
	if err != nil {
		logger.Errorf("UpdateUser.Update fail, id: %s, error: %s", id, err)
		return err
	}

	return nil
}

func (s *UserService) VerifyEmail(ctx context.Context, request dto.VerifyEmailRequest) (dto.VerifyResponse, error) {
	user, err := s.repo.FindByEmailAndVerifyCode(ctx, request.Email, request.VerifyCodeEmail)
	if err != nil {
		return dto.VerifyResponse{Message: "Verification failed"}, err
	}

	if user == nil {
		return dto.VerifyResponse{Message: "Verify code not correct"}, errors.New("verify code not correct")
	}

	user.ApproveEmail = true
	if err := s.repo.Update(ctx, user); err != nil {
		return dto.VerifyResponse{Message: "Failed to update user"}, err
	}

	return dto.VerifyResponse{Message: "Verification successful"}, nil
}

func (s *UserService) VerifyPhoneNumber(ctx context.Context, request dto.VerifyPhoneNumberRequest) (dto.VerifyResponse, error) {
	user, err := s.repo.FindByPhoneAndVerifyCode(ctx, request.PhoneNumber, request.VerifyCodePhoneNumber)
	if err != nil {
		return dto.VerifyResponse{Message: "Verification failed"}, err
	}

	if user == nil {
		return dto.VerifyResponse{Message: "Verify code not correct"}, errors.New("verify code not correct")
	}

	user.ApprovePhoneNumber = true
	if err := s.repo.Update(ctx, user); err != nil {
		return dto.VerifyResponse{Message: "Failed to update user"}, err
	}

	return dto.VerifyResponse{Message: "Verification successful"}, nil
}

func (s *UserService) ResendVerfiyCodePhone(ctx context.Context, request dto.ResendVerifyPhoneNumberRequest) (dto.VerifyResponse, error) {

	user, err := s.repo.FindByPhone(ctx, request.PhoneNumber)
	if err != nil {
		return dto.VerifyResponse{Message: "Resend Verification failed"}, err
	}

	if user == nil {
		return dto.VerifyResponse{Message: "Resend Verify code not correct"}, errors.New("verify code not correct")
	}
	// sendVerificationCodes sends the verification codes to the user's phone and email
	user.BeforeUpdateVerificationPhone()
	if err := s.repo.Update(ctx, user); err != nil {
		return dto.VerifyResponse{Message: "Failed to Resend user"}, err
	}

	return dto.VerifyResponse{Message: "Resend Verify code is successful"}, nil
}

func (s *UserService) ResendVerfiyCodeEmail(ctx context.Context, request dto.ResendVerifyEmailRequest) (dto.VerifyResponse, error) {

	user, err := s.repo.FindByEmail(ctx, request.Email)
	if err != nil {
		return dto.VerifyResponse{Message: "Resend Verification failed"}, err
	}

	if user == nil {
		return dto.VerifyResponse{Message: "Resend Verify code not correct"}, errors.New("verify code not correct")
	}
	// sendVerificationCodes sends the verification codes to the user's phone and email
	user.BeforeUpdateVerificationEmail()
	if err := s.repo.Update(ctx, user); err != nil {
		return dto.VerifyResponse{Message: "Failed to Resend Resend"}, err
	}

	return dto.VerifyResponse{Message: "Resend Verify code is successful"}, nil
}

func (p *UserService) ListUsers(ctx context.Context, req dto.ListUsersReq) ([]*model.User, *paging.Pagination, error) {
	Userss, pagination, err := p.repo.ListUsers(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return Userss, pagination, nil
}

func (p *UserService) Delete(ctx context.Context, id string, req *dto.DeleteUserReq) (*model.User, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	User, err := p.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("Delete.GetUserByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	utils.Copy(User, req)
	err = p.repo.Delete(ctx, User)
	if err != nil {
		logger.Errorf("Delete fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return User, nil
}

func (uc *UserService) LoginWithGoogle(ctx context.Context, code string) (*model.User, string, string, error) {
	token, err := uc.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, "", "", err
	}

	oauth2Service, err := googleOauth2.NewService(ctx, option.WithTokenSource(uc.oauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, "", "", err
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, "", "", err
	}

	user, err := uc.repo.FindOrCreateByGoogleID(ctx, userInfo.Id, userInfo.Email, userInfo.Name)
	if err != nil {
		return nil, "", "", err
	}
	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	accessToken := jtoken.GenerateAccessToken(tokenData)
	refreshToken := jtoken.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

func (s *UserService) LoginWithFacebook(ctx context.Context, code string) (*model.User, string, string, error) {
	token, err := s.fbOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, "", "", err
	}

	// Assuming Facebook does not have a similar NewService function like Google.
	// You need to use Facebook Graph API to get user information here.

	// Pseudo code:
	userInfo, err := fetchFacebookUserInfo(token)
	if err != nil {
		return nil, "", "", err
	}

	user, err := s.repo.FindOrCreateByFacebookID(ctx, userInfo.Id, userInfo.Email, userInfo.Name)
	if err != nil {
		return nil, "", "", err
	}

	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	accessToken := jtoken.GenerateAccessToken(tokenData)
	refreshToken := jtoken.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

// You need to implement the fetchFacebookUserInfo function to interact with the Facebook Graph API.
func fetchFacebookUserInfo(token *oauth2.Token) (*FacebookUserInfo, error) {
	// Implement your Facebook Graph API call here.
	return &FacebookUserInfo{
		Id:    "example-id",
		Email: "example-email",
		Name:  "example-name",
	}, nil
}

type FacebookUserInfo struct {
	Id    string
	Email string
	Name  string
}
