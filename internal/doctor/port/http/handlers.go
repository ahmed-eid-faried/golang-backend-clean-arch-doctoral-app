package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"main/internal/doctor/dto"
	"main/internal/doctor/service"
	"main/pkg/config"
	"main/pkg/redis"
	"main/pkg/response"
	"main/pkg/utils"
)

// Doctor
// Doctor
type DoctorHandler struct {
	cache   redis.IRedis
	service service.IDoctorService
}

func NewDoctorHandler(
	cache redis.IRedis,
	service service.IDoctorService,
) *DoctorHandler {
	return &DoctorHandler{
		cache:   cache,
		service: service,
	}
}

// GetDoctorByID godoc
//
//	@Summary	Get Doctor by id
//	@Tags		Doctor
//	@Produce	json
//	@Param		id	path	string	true	"Doctor ID"
//	@Success	200	{object}	dto.Doctor
//	@Router		/doctor/{id} [get]
func (p *DoctorHandler) GetDoctorByID(c *gin.Context) {

	DoctorId := c.Param("id")
	Doctor, err := p.service.GetDoctorByID(c, DoctorId)
	if err != nil {
		logger.Error("Failed to get Doctor detail: ", err)
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	var res dto.Doctor
	cacheKey := c.Request.URL.RequestURI()
	err2 := p.cache.Get(cacheKey, &res)
	if err2 == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}
	utils.Copy(&res, &Doctor)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(cacheKey, res, config.DoctorCachingTime.Abs())
}

// ListDoctor godoc
// @Summary	ListDoctors
// @Tags		Doctor
// @Produce	json
// @Security	ApiKeyAuth
// @Param		search	query	string	false	"Search"
// @Param		id_user	query	string	false	"ID User"
// @Param		page	query	int64	false	"Page number"
// @Param		limit	query	int64	false	"Limit per page"
// @Param		order_by	query	string	false	"Order by field"
// @Param		order_desc	query	bool	false	"Order descending"
// @Success	200	{object}	dto.ListDoctorRes
// @Router		/doctor/list_doctors [get]
func (p *DoctorHandler) ListDoctors(c *gin.Context) {
	var req dto.ListDoctorReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to get query params", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	var res dto.ListDoctorRes
	cacheKey := c.Request.URL.RequestURI()
	err := p.cache.Get(cacheKey, &res)
	if err == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}

	Doctors, pagination, err := p.service.ListDoctors(c, &req)
	if err != nil {
		logger.Error("Failed to get list of Doctors: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	utils.Copy(&res.Doctors, &Doctors)
	res.Pagination = pagination
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(cacheKey, res, config.DoctorCachingTime)
}

// CreateDoctor godoc
//
//	@Summary	create Doctor
//	@Tags		Doctor
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body	dto.CreateDoctorReq	true	"Body"
//	@Success	200	{object}	dto.Doctor
//	@Router		/doctor [post]
func (p *DoctorHandler) CreateDoctor(c *gin.Context) {
	var req dto.CreateDoctorReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	Doctor, err := p.service.Create(c, &req)
	if err != nil {
		logger.Error("Failed to create Doctor", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Doctor
	utils.Copy(&res, &Doctor)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*Doctor*")
}

// UpdateDoctor godoc
//
//	@Summary	Update Doctor
//	@Tags		Doctor
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body	dto.UpdateDoctorReq	true	"Body"
//	@Success	200	{object}	dto.Doctor
//	@Router		/doctor/{id} [put]
func (p *DoctorHandler) UpdateDoctor(c *gin.Context) {
	//	@Param		id	path	string					true	"Doctor ID"
	// DoctorId := c.Param("id")
	var req dto.UpdateDoctorReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	Doctor, err := p.service.Update(c, req.ID, &req)
	if err != nil {
		logger.Error("Failed to Update Doctor", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Doctor
	utils.Copy(&res, &Doctor)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*Doctor*")
}

// DeleteDoctor godoc
//
//	@Summary	Delete Doctor
//	@Tags		Doctor
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body	dto.DeleteDoctorReq	true	"Body"
//	@Success	200	{object}	dto.Doctor
//	@Router		/doctor/{id} [Delete]
func (p *DoctorHandler) DeleteDoctor(c *gin.Context) {
	//	@Param		id	path	string					true	"Doctor ID"
	// DoctorId := c.Param("id")
	var req dto.DeleteDoctorReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	Doctor, err := p.service.Delete(c, req.ID, &req)
	if err != nil {
		logger.Error("Failed to Delete Doctor", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Doctor
	utils.Copy(&res, &Doctor)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*Doctor*")
}
