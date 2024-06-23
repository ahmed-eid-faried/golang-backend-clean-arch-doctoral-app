package grpc

import (
	"fmt"
	"net"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// cartGRPC "main/internal/cart/port/grpc"
	addressGRPC "main/internal/address/port/grpc"
	userGRPC "main/internal/user/port/grpc"
	"main/pkg/config"
	"main/pkg/dbs"
	"main/pkg/middleware"
	"main/pkg/redis"
)

type Server struct {
	engine        *grpc.Server
	cfg           *config.Schema
	validator     validation.Validation
	db            dbs.IDatabase
	cache         redis.IRedis
	oauthConfig   *oauth2.Config
	fbOauthConfig *oauth2.Config
}

func NewServer(validator validation.Validation, db dbs.IDatabase, cache redis.IRedis, oauthConfig *oauth2.Config, fbOauthConfig *oauth2.Config) *Server {
	interceptor := middleware.NewAuthInterceptor(config.AuthIgnoreMethods)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Unary(),
		),
	)

	return &Server{
		engine:        grpcServer,
		cfg:           config.GetConfig(),
		validator:     validator,
		db:            db,
		cache:         cache,
		oauthConfig:   oauthConfig,
		fbOauthConfig: fbOauthConfig,
	}
}

func (s Server) Run() error {
	userGRPC.RegisterHandlers(s.engine, s.db, s.validator, s.oauthConfig, s.fbOauthConfig)
	addressGRPC.RegisterHandlers(s.engine, s.db, s.validator, s.cache)
	// cartGRPC.RegisterHandlers(s.engine, s.db, s.validator)

	reflection.Register(s.engine)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GrpcPort))
	logger.Info("GRPC server is listening on PORT: ", s.cfg.GrpcPort)
	if err != nil {
		logger.Error("Failed to listen: ", err)
		return err
	}

	// Start grpc server
	err = s.engine.Serve(lis)
	if err != nil {
		logger.Fatal("Failed to serve grpc: ", err)
		return err
	}

	return nil
}
