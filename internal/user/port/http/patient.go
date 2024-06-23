package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"main/internal/user/dto"
	"main/internal/user/model"
	"main/pkg/response"
	"main/pkg/utils"
)

// Login , Register, GetMe, RefreshToken, UpdateUser, VerfiyCode FOR PHONE

// Login godoc
//
//	@Summary	Login
//	@Tags		users-patient
//	@Produce	json
//	@Param		_	body		dto.KLoginReq	true	"Body"
//	@Success	200	{object}	dto.LoginRes
//	@Router		/auth-patient/login [post]
func (h *UserHandler) LoginPatient(c *gin.Context) {
	var req dto.KLoginReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body ", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	req2 := dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
		Role:     model.UserRoleClient,
	}
	user, accessToken, refreshToken, err := h.service.Login(c, &req2)
	if err != nil {
		logger.Error("Failed to login ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.LoginRes
	utils.Copy(&res.User, &user)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	response.JSON(c, http.StatusOK, res)
}

// Register godoc
//
//	@Summary	Register new user
//	@Tags		users-patient
//	@Produce	json
//	@Param		_	body		dto.KRegisterReq	true	"Body"
//	@Success	200	{object}	dto.RegisterRes
//	@Router		/auth-patient/register [post]
func (h *UserHandler) RegisterPatient(c *gin.Context) {
	var req dto.KRegisterReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	req2 := dto.RegisterReq{
		Email:       req.Email,
		Password:    req.Password,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Role:        model.UserRoleClient,
	}
	user, err := h.service.Register(c, &req2)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.RegisterRes
	utils.Copy(&res.User, &user)
	response.JSON(c, http.StatusOK, res)
}

// UpdateUser godoc
//
//	@Summary	changes the password
//	@Tags		users-patient
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body	dto.KUpdateUserReq	true	"Body"
//	@Success	200	{object}	dto.UpdateUserRes
//	@Router		/auth-patient/update-user [put]
func (h *UserHandler) UpdatePatient(c *gin.Context) {
	var req dto.KUpdateUserReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	req2 := dto.UpdateUserReq{
		ID:          req.ID,
		Email:       req.Email,
		Password:    req.Password,
		NewPassword: req.NewPassword,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Role:        model.UserRoleClient,
	}
	userID := c.GetString("userId")
	err := h.service.UpdateUser(c, userID, &req2)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
