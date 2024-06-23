package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"main/internal/address/dto"
	"main/internal/address/service"
	"main/pkg/config"
	"main/pkg/redis"
	"main/pkg/response"
	"main/pkg/utils"
)

// Address
// address
type AddressHandler struct {
	cache   redis.IRedis
	service service.IAddressService
}

func NewAddressHandler(
	cache redis.IRedis,
	service service.IAddressService,
) *AddressHandler {
	return &AddressHandler{
		cache:   cache,
		service: service,
	}
}

// GetAddressByID godoc
//
//	@Summary	Get Address by id
//	@Tags		Address
//	@Produce	json
//	@Param		id	path	string	true	"Address ID"
//	@Success	200	{object}	dto.Address
//	@Router		/address/{id} [get]
func (p *AddressHandler) GetAddressByID(c *gin.Context) {

	AddressId := c.Param("id")
	Address, err := p.service.GetAddressByID(c, AddressId)
	if err != nil {
		logger.Error("Failed to get Address detail: ", err)
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	var res dto.Address
	cacheKey := c.Request.URL.RequestURI()
	err2 := p.cache.Get(cacheKey, &res)
	if err2 == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}
	utils.Copy(&res, &Address)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(cacheKey, res, config.AddressCachingTime.Abs())
}

// ListAddress godoc
//
//		@Summary	Get list Address
//		@Tags		Address
//		@Produce	json
//	 @Param id_user header string true "id_user"
//	@Param		name	path	string					true	"name"
//	@Param		page	path	string					true	"page"
//	@Param		limit	path	string					true	"limit"
//
// @Success	200	{object}	dto.ListAddressRes
// @Router		/address  [get]
func (p *AddressHandler) ListAddresses(c *gin.Context) {
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

	// var req dto.ListAddressReq
	req := dto.ListAddressReq{
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

	var res dto.ListAddressRes
	cacheKey := c.Request.URL.RequestURI()
	err := p.cache.Get(cacheKey, &res)
	if err == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}

	Addresses, pagination, err := p.service.ListAddresses(c, &req)
	if err != nil {
		logger.Error("Failed to get list Address: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	utils.Copy(&res.Addresses, &Addresses)
	res.Pagination = pagination
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(cacheKey, res, config.AddressCachingTime)
}

// CreateAddress godoc
//
//	@Summary	create Address
//	@Tags		Address
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body	dto.CreateAddressReq	true	"Body"
//	@Success	200	{object}	dto.Address
//	@Router		/address [post]
func (p *AddressHandler) CreateAddress(c *gin.Context) {
	var req dto.CreateAddressReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	Address, err := p.service.Create(c, &req)
	if err != nil {
		logger.Error("Failed to create Address", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Address
	utils.Copy(&res, &Address)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*Address*")
}

// UpdateAddress godoc
//
//	@Summary	Update Address
//	@Tags		Address
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body	dto.UpdateAddressReq	true	"Body"
//	@Success	200	{object}	dto.Address
//	@Router		/address/{id} [put]
func (p *AddressHandler) UpdateAddress(c *gin.Context) {
	//	@Param		id	path	string					true	"Address ID"
	// AddressId := c.Param("id")
	var req dto.UpdateAddressReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	Address, err := p.service.Update(c, req.ID, &req)
	if err != nil {
		logger.Error("Failed to Update Address", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Address
	utils.Copy(&res, &Address)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*Address*")
}

// DeleteAddress godoc
//
//	@Summary	Delete Address
//	@Tags		Address
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body	dto.DeleteAddressReq	true	"Body"
//	@Success	200	{object}	dto.Address
//	@Router		/address/{id} [Delete]
func (p *AddressHandler) DeleteAddress(c *gin.Context) {
	//	@Param		id	path	string					true	"Address ID"
	// AddressId := c.Param("id")
	var req dto.DeleteAddressReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	Address, err := p.service.Delete(c, req.ID, &req)
	if err != nil {
		logger.Error("Failed to Delete Address", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Address
	utils.Copy(&res, &Address)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*Address*")
}
