package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	"main/internal/doctor/repository"
	"main/internal/doctor/service"
	"main/pkg/dbs"
	"main/pkg/middleware"
	"main/pkg/redis"
)

func Routes(r *gin.RouterGroup, sqlDB dbs.IDatabase, validator validation.Validation, cache redis.IRedis) {
	doctorRepo := repository.NewDoctorRepository(sqlDB)
	doctorSvc := service.NewDoctorService(validator, doctorRepo)
	doctorHandler := NewDoctorHandler(cache, doctorSvc)

	authMiddleware := middleware.JWTAuth()
	doctorRoute := r.Group("/doctor")
	{
		doctorRoute.GET("", doctorHandler.ListDoctors)
		doctorRoute.GET("/:id", doctorHandler.GetDoctorByID)
		doctorRoute.POST("", authMiddleware, doctorHandler.CreateDoctor)
		doctorRoute.PUT("/:id", authMiddleware, doctorHandler.UpdateDoctor)
		doctorRoute.DELETE("/:id", authMiddleware, doctorHandler.DeleteDoctor)
	}
}
