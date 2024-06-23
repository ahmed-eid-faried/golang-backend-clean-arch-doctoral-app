package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"main/internal/user/dto"
	"main/internal/user/service"
	"main/pkg/redis"
	"main/pkg/response"
	"main/pkg/utils"
)

type UserHandler struct {
	cache   redis.IRedis
	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetMe godoc
//
//	@Summary	get my profile
//	@Tags		users
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Success	200	{object}	dto.User
//	@Router		/auth/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	user, err := h.service.GetUserByID(c, userID)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.User
	utils.Copy(&res, &user)
	response.JSON(c, http.StatusOK, res)
}

// GetMe godoc
//
//	@Summary	get my profile
//	@Tags		users
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body		dto.RefreshTokenReq	true	"Body"
//	@Success	200	{object}	dto.RefreshTokenRes
//	@Router		/auth/refresh-token [get]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	accessToken, err := h.service.RefreshToken(c, userID)
	if err != nil {
		logger.Error("Failed to refresh token", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	res := dto.RefreshTokenRes{
		AccessToken: accessToken,
	}
	response.JSON(c, http.StatusOK, res)
}

// VerfiyCodeEmail godoc
//
//	@Summary	Verfiy Code for Email
//	@Tags		users
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body	dto.VerifyEmailRequest	true	"Body"
//	@Success	200	{object}	dto.VerifyResponse
//	@Router		/auth/verfiy-code-email [put]
func (h *UserHandler) VerfiyCodeEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	resp, err := h.service.VerifyEmail(c, req)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, resp)
}

// VerfiyCodePhoneNumber godoc
//
//	@Summary	Verfiy Code for PhoneNumber
//	@Tags		users
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body	dto.VerifyPhoneNumberRequest	true	"Body"
//	@Success	200	{object}	dto.VerifyResponse
//	@Router		/auth/verfiy-code-phone-number [put]
func (h *UserHandler) VerfiyCodePhoneNumber(c *gin.Context) {
	var req dto.VerifyPhoneNumberRequest
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	resp, err := h.service.VerifyPhoneNumber(c, req)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, resp)
}

// VerfiyCodePhoneNumber godoc
//
//	@Summary	Verfiy Code for PhoneNumber
//	@Tags		users
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body	dto.ResendVerifyPhoneNumberRequest	true	"Body"
//	@Success	200	{object}	dto.VerifyResponse
//	@Router		/auth/resend-verfiy-code-phone-number [put]
func (h *UserHandler) VerfiyCodePhoneNumberResend(c *gin.Context) {
	var req dto.ResendVerifyPhoneNumberRequest
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	resp, err := h.service.ResendVerfiyCodePhone(c, req)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, resp)
}

// VerfiyCodeEmail godoc
//
//	@Summary	Verfiy Code for Email
//	@Tags		users
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body	dto.ResendVerifyEmailRequest	true	"Body"
//	@Success	200	{object}	dto.VerifyResponse
//	@Router		/auth/resend-verfiy-code-email [put]
func (h *UserHandler) VerfiyCodeEmailResend(c *gin.Context) {
	var req dto.ResendVerifyEmailRequest
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	resp, err := h.service.ResendVerfiyCodeEmail(c, req)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, resp)
}

// LoginWithGoogleHandler handles Google OAuth2 login flow

// @Summary Login with Google
//
//	@Tags		users
//
// @Description Logs user in with Google OAuth2
// @Accept json
// @Produce json
// @Param code query string true "Authorization code received from Google"
//
//	@Success	200	{object}	dto.LoginRes
//
// @Router /auth/google/callback [get]
func (h *UserHandler) LoginWithGoogleHandler(c *gin.Context) {

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}
	// var user dto.LoginRes
	// Use a user service to handle login or registration
	user, accessToken, refreshToken, err := h.service.LoginWithGoogle(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login or register user"})
		return
	}

	var res dto.LoginRes
	utils.Copy(&res.User, &user)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	response.JSON(c, http.StatusOK, res)

	c.JSON(http.StatusOK, res)
}

// HandleFacebookCallback handles Facebook OAuth2 callback

// @Router /auth/facebook/callback [get]
// @Summary Handles Facebook OAuth2 callback
//
//	@Tags		users
//
// @Param code query string true "Authorization Code from Facebook"
// @Produce json
//
//	@Success	200	{object}	dto.LoginRes
func (h *UserHandler) HandleFacebookCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state or code parameter"})
		return
	}

	// Use a user service to handle login or registration
	user, accessToken, refreshToken, err := h.service.LoginWithFacebook(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res dto.LoginRes
	utils.Copy(&res.User, &user)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	response.JSON(c, http.StatusOK, res)

	c.JSON(http.StatusOK, res)
}
