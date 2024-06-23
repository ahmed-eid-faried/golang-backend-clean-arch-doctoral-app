package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"

	"main/internal/address/repository"
	"main/internal/address/service"
	"main/pkg/dbs"
	"main/pkg/redis"
	pb "main/proto/gen/go/address"
)

func RegisterHandlers(svr *grpc.Server, db dbs.IDatabase, validator validation.Validation, cache redis.IRedis) {
	AddressRepo := repository.NewAddressRepository(db)
	AddressSvc := service.NewAddressService(validator, AddressRepo)
	AddressHandler := NewAddressHandler(cache, AddressSvc)

	pb.RegisterAddressServiceServer(svr, AddressHandler)
}
