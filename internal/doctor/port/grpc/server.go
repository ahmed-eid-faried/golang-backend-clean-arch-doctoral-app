package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"

	"main/internal/doctor/repository"
	"main/internal/doctor/service"
	"main/pkg/dbs"
	"main/pkg/redis"
	pb "main/proto/gen/go/doctor"
)

func RegisterHandlers(svr *grpc.Server, db dbs.IDatabase, validator validation.Validation, cache redis.IRedis) {
	DoctorRepo := repository.NewDoctorRepository(db)
	DoctorSvc := service.NewDoctorService(validator, DoctorRepo)
	DoctorHandler := NewDoctorHandler(cache, DoctorSvc)

	pb.RegisterDoctorServiceServer(svr, DoctorHandler)
}
