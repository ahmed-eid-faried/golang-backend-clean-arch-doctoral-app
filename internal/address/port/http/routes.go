package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	"main/internal/address/repository"
	"main/internal/address/service"
	"main/pkg/dbs"
	"main/pkg/middleware"
	"main/pkg/redis"
)

func Routes(r *gin.RouterGroup, sqlDB dbs.IDatabase, validator validation.Validation, cache redis.IRedis) {
	addressRepo := repository.NewAddressRepository(sqlDB)
	addressSvc := service.NewAddressService(validator, addressRepo)
	addressHandler := NewAddressHandler(cache, addressSvc)

	authMiddleware := middleware.JWTAuth()
	AddressRoute := r.Group("/address")
	{
		AddressRoute.GET("", addressHandler.ListAddresses)
		AddressRoute.GET("/:id", addressHandler.GetAddressByID)
		AddressRoute.POST("", authMiddleware, addressHandler.CreateAddress)
		AddressRoute.PUT("/:id", authMiddleware, addressHandler.UpdateAddress)
		AddressRoute.DELETE("/:id", authMiddleware, addressHandler.DeleteAddress)
	}
}
