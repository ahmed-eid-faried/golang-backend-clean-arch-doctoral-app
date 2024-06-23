package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"main/internal/user/dto"
	"main/pkg/config"
	"main/pkg/response"
	"main/pkg/utils"
)

// ListUsers DeleteAdmin CreateAdmin UpdateAdmin LoginAdmin

// ListUsers godoc
//
//		@Summary	Get list Users
//		@Tags		users-admin
//		@Produce	json
//	@Security	ApiKeyAuth
//	 @Param id_user header string true "id_user"
//	@Param		name	path	string					true	"name"
//	@Param		page	path	string					true	"page"
//	@Param		limit	path	string					true	"limit"
//
// @Success	200	{object}	dto.ListUsersRes
// @Router		/auth-admin/users  [get]
func (p *UserHandler) ListUsers(c *gin.Context) {
	Name := c.Param("name")
	IDUser := c.Param("id_user")
	PageStr := c.Param("page")
	LimitStr := c.Param("limit")

	Page, err1 := strconv.ParseInt(PageStr, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{Code: http.StatusBadRequest, Message: "Invalid page number"})
		return
	}

	Limit, err2 := strconv.ParseInt(LimitStr, 10, 64)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{Code: http.StatusBadRequest, Message: "Invalid limit number"})
		return
	}

	// var req dto.ListUsersReq
	req := dto.ListUsersReq{
		Name:   Name,
		IDUser: IDUser,
		Page:   Page,
		Limit:  Limit,
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to parse request query: ", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	var res dto.ListUsersRes
	cacheKey := c.Request.URL.RequestURI()
	err := p.cache.Get(cacheKey, &res)
	if err == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}

	Users, pagination, err := p.service.ListUsers(c, req)
	if err != nil {
		logger.Error("Failed to get list Users: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	utils.Copy(&res.Users, &Users)
	res.Pagination = pagination
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(cacheKey, res, config.UsersCachingTime)
}

// DeleteUser godoc
//
//	@Summary	Delete User
//	@Tags		users-admin
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body	dto.DeleteUserReq	true	"Body"
//	@Success	200	{object}	dto.User
//	@Router		/auth-admin/{id} [Delete]
func (p *UserHandler) DeleteAdmin(c *gin.Context) {
	var req dto.DeleteUserReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	User, err := p.service.Delete(c, req.ID, &req)
	if err != nil {
		logger.Error("Failed to Delete User", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.User
	utils.Copy(&res, &User)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*User*")
}

// Create godoc
//
//	@Summary	Create new user
//	@Tags		users-admin
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body		dto.RegisterReq	true	"Body"
//	@Success	200	{object}	dto.RegisterRes
//	@Router		/auth-admin/create [post]
func (h *UserHandler) CreateAdmin(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, err := h.service.Register(c, &req)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.RegisterRes
	utils.Copy(&res.User, &user)
	response.JSON(c, http.StatusOK, res)
}

// UpdateAdmin godoc
//
//	@Summary	Update user
//	@Tags		users-admin
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body	dto.UpdateUserReq	true	"Body"
//	@Success	200	{object}	dto.UpdateUserRes
//	@Router		/auth-admin/update [put]
func (h *UserHandler) UpdateAdmin(c *gin.Context) {
	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	userID := c.GetString("userId")
	err := h.service.UpdateUser(c, userID, &req)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}

// LoginAdmin godoc
//
//	@Summary	LoginAdmin
//	@Tags		users-admin
//	@Produce	json
//	@Param		_	body		dto.LoginReq	true	"Body"
//	@Success	200	{object}	dto.LoginRes
//	@Router		/auth-admin/login [post]
func (h *UserHandler) LoginAdmin(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body ", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c, &req)
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
