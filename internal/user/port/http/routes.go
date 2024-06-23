package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
	"golang.org/x/oauth2"

	"main/internal/user/repository"
	"main/internal/user/service"
	"main/pkg/dbs"
	"main/pkg/middleware"
)

func Routes(r *gin.RouterGroup, sqlDB dbs.IDatabase, validator validation.Validation, fbOauthConfig *oauth2.Config, oauthConfig *oauth2.Config) {
	userRepo := repository.NewUserRepository(sqlDB)
	userSvc := service.NewUserService(validator, oauthConfig, fbOauthConfig, userRepo)
	userHandler := NewUserHandler(userSvc)

	authMiddleware := middleware.JWTAuth()
	refreshAuthMiddleware := middleware.JWTRefresh()

	// GetMe RefreshToken  --  VerfiyCodeEmail  VerfiyCodePhoneNumber  VerfiyCodePhoneNumberResend VerfiyCodeEmailResend
	authRoute := r.Group("/auth")
	{
		authRoute.GET("/google/callback", userHandler.LoginWithGoogleHandler)
		authRoute.GET("/facebook/callback", userHandler.HandleFacebookCallback)
		authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		authRoute.POST("/refresh-token", refreshAuthMiddleware, userHandler.RefreshToken)
		//for doctor or Patient only
		authRoute.PUT("/verfiy-code-email", authMiddleware, userHandler.VerfiyCodeEmail)
		authRoute.PUT("/verfiy-code-phone-number", authMiddleware, userHandler.VerfiyCodePhoneNumber)
		authRoute.PUT("/resend-verfiy-code-phone-number", authMiddleware, userHandler.VerfiyCodePhoneNumberResend)
		authRoute.PUT("/resend-verfiy-code-email", authMiddleware, userHandler.VerfiyCodeEmailResend)
	}

	// ListUsers DeleteAdmin CreateAdmin UpdateAdmin LoginAdmin
	authRouteAdmin := r.Group("/auth-admin")
	{
		authRouteAdmin.POST("/login", userHandler.LoginAdmin)
		authRouteAdmin.POST("/create", userHandler.CreateAdmin)
		authRouteAdmin.PUT("/update", authMiddleware, userHandler.UpdateAdmin)
		authRouteAdmin.GET("/users", authMiddleware, userHandler.ListUsers)
		authRouteAdmin.DELETE("/", authMiddleware, userHandler.DeleteAdmin)
	}

	// LoginDoctor RegisterDoctor UpdateDoctor
	authRouteDoctor := r.Group("/auth-doctor")
	{
		authRouteDoctor.POST("/login", userHandler.LoginDoctor)
		authRouteDoctor.POST("/register", userHandler.RegisterDoctor)
		authRouteDoctor.PUT("/update-user", authMiddleware, userHandler.UpdateDoctor)
	}
	// LoginPatient RegisterPatient UpdatePatient
	authRoutePatient := r.Group("/auth-patient")
	{
		authRoutePatient.POST("/login", userHandler.LoginPatient)
		authRoutePatient.POST("/register", userHandler.RegisterPatient)
		authRoutePatient.PUT("/update-user", authMiddleware, userHandler.UpdatePatient)
	}
}
