package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"

	"main/internal/user/repository"
	"main/internal/user/service"
	"main/pkg/dbs"
	pb "main/proto/gen/go/user"
)

func RegisterHandlers(svr *grpc.Server, db dbs.IDatabase, validator validation.Validation, oauthConfig *oauth2.Config, fbOauthConfig *oauth2.Config) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(validator, oauthConfig, fbOauthConfig, userRepo)
	userHandler := NewUserHandler(userSvc)

	pb.RegisterUserServiceServer(svr, userHandler)
}
