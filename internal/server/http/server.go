package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/oauth2"

	_ "main/docs"
	// orderHttp "main/internal/order/port/http"
	addressHttp "main/internal/address/port/http"
	doctorHttp "main/internal/doctor/port/http"
	userHttp "main/internal/user/port/http"
	"main/pkg/config"
	"main/pkg/dbs"
	"main/pkg/redis"
	"main/pkg/response"
)

type Server struct {
	engine        *gin.Engine
	cfg           *config.Schema
	validator     validation.Validation
	db            dbs.IDatabase
	cache         redis.IRedis
	oauthConfig   *oauth2.Config
	fbOauthConfig *oauth2.Config
}

func NewServer(validator validation.Validation, db dbs.IDatabase, cache redis.IRedis, fbOauthConfig *oauth2.Config,
	oauthConfig *oauth2.Config) *Server {
	return &Server{
		engine:        gin.Default(),
		cfg:           config.GetConfig(),
		validator:     validator,
		db:            db,
		cache:         cache,
		oauthConfig:   oauthConfig,
		fbOauthConfig: fbOauthConfig,
	}
}

func (s Server) Run() error {
	_ = s.engine.SetTrustedProxies(nil)
	if s.cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Handle 404 errors
	s.engine.NoRoute(func(c *gin.Context) {
		// Set the status code to 404
		// c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
		c.File("./templates/404.html")
	})
	s.engine.GET("/health", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, nil)
		return
	})

	// Start http server
	logger.Info("HTTP server is listening on PORT: ", s.cfg.HttpPort)
	if err := s.engine.Run(fmt.Sprintf(":%d", s.cfg.HttpPort)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}

func (s Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s Server) MapRoutes() error {
	v1 := s.engine.Group("/api/v1")
	userHttp.Routes(v1, s.db, s.validator, s.fbOauthConfig, s.oauthConfig)
	addressHttp.Routes(v1, s.db, s.validator, s.cache)
	doctorHttp.Routes(v1, s.db, s.validator, s.cache)
	// orderHttp.Routes(v1, s.db, s.validator)
	return nil
}
