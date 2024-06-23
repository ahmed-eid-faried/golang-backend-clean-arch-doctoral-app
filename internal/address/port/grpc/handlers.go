package grpc

import (
	"context"
	"time"

	"github.com/quangdangfit/gocommon/logger"

	"main/internal/address/dto"
	"main/internal/address/service"
	"main/pkg/config"
	"main/pkg/redis"
	"main/pkg/utils"
	pb "main/proto/gen/go/address"
)

type AddressHandler struct {
	cache   redis.IRedis
	service service.IAddressService
	pb.UnimplementedAddressServiceServer
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

// // Convert time.Time to string
//
//	func formatTimeToString(t time.Time) string {
//		return t.Format(time.RFC3339)
//	}

func (h *AddressHandler) GetAddressByID(ctx context.Context, req *pb.GetAddressByIDRequest) (*pb.AddressResponse, error) {
	var res dto.Address
	cacheKey := "address_" + req.Id
	err := h.cache.Get(cacheKey, &res)
	if err == nil {
		return &pb.AddressResponse{Address: &pb.Address{
			IdAddress: res.ID,
			IdUser:    res.IDUser,
			Name:      res.Name,
			City:      res.City,
			Street:    res.Street,
			Lat:       res.Lat,
			Long:      res.Long,

			// CreatedAt: formatTimeToString(res.CreatedAt),
			// UpdatedAt: formatTimeToString(res.UpdatedAt),
		}}, nil
	}

	address, err := h.service.GetAddressByID(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to get address detail: ", err)
		return nil, err
	}

	utils.Copy(&res, &address)
	_ = h.cache.SetWithExpiration(cacheKey, res, config.AddressCachingTime.Abs())
	return &pb.AddressResponse{Address: &pb.Address{
		IdAddress: res.ID,
		IdUser:    res.IDUser,
		Name:      res.Name,
		City:      res.City,
		Street:    res.Street,
		Lat:       res.Lat,
		Long:      res.Long,
	}}, nil
}

func (h *AddressHandler) ListAddresses(ctx context.Context, req *pb.ListAddressesRequest) (*pb.ListAddressesResponse, error) {
	var res dto.ListAddressRes
	cacheKey := "addresses_list"
	err := h.cache.Get(cacheKey, &res)
	if err == nil {
		var pbAddresses []*pb.Address
		for _, addr := range res.Addresses {
			pbAddresses = append(pbAddresses, &pb.Address{
				IdAddress: addr.ID,
				IdUser:    addr.IDUser,
				Name:      addr.Name,
				City:      addr.City,
				Street:    addr.Street,
				Lat:       addr.Lat,
				Long:      addr.Long,
			})
		}
		return &pb.ListAddressesResponse{Addresses: pbAddresses}, nil
	}

	addresses, pagination, err := h.service.ListAddresses(ctx, &dto.ListAddressReq{})
	if err != nil {
		logger.Error("Failed to get list of addresses: ", err)
		return nil, err
	}

	utils.Copy(&res.Addresses, &addresses)
	res.Pagination = pagination
	_ = h.cache.SetWithExpiration(cacheKey, res, time.Hour) // Adjust caching time as needed

	var pbAddresses []*pb.Address
	for _, addr := range addresses {
		pbAddresses = append(pbAddresses, &pb.Address{
			IdAddress: addr.ID,
			IdUser:    addr.IDUser,
			Name:      addr.Name,
			City:      addr.City,
			Street:    addr.Street,
			Lat:       addr.Lat,
			Long:      addr.Long,
		})
	}
	return &pb.ListAddressesResponse{Addresses: pbAddresses}, nil
}

func (h *AddressHandler) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.AddressResponse, error) {
	var addressDTO dto.CreateAddressReq
	addressDTO.IDUser = req.Request.IdUser
	addressDTO.Name = req.Request.Name
	addressDTO.City = req.Request.City
	addressDTO.Street = req.Request.Street
	addressDTO.Lat = req.Request.Lat
	addressDTO.Long = req.Request.Long

	address, err := h.service.Create(ctx, &addressDTO)
	if err != nil {
		logger.Error("Failed to create address: ", err)
		return nil, err
	}

	var res dto.Address
	utils.Copy(&res, &address)
	_ = h.cache.RemovePattern("*address*")
	return &pb.AddressResponse{Address: &pb.Address{
		IdAddress: res.ID,
		IdUser:    res.IDUser,
		Name:      res.Name,
		City:      res.City,
		Street:    res.Street,
		Lat:       res.Lat,
		Long:      res.Long,
	}}, nil
}

func (h *AddressHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.AddressResponse, error) {
	var addressDTO dto.UpdateAddressReq
	addressDTO.Name = req.Request.Name
	addressDTO.City = req.Request.City
	addressDTO.Street = req.Request.Street
	addressDTO.Lat = req.Request.Lat
	addressDTO.Long = req.Request.Long

	address, err := h.service.Update(ctx, req.Id, &addressDTO)
	if err != nil {
		logger.Error("Failed to update address: ", err)
		return nil, err
	}

	var res dto.Address
	utils.Copy(&res, &address)
	_ = h.cache.RemovePattern("*address*")
	return &pb.AddressResponse{Address: &pb.Address{
		IdAddress: res.ID,
		IdUser:    res.IDUser,
		Name:      res.Name,
		City:      res.City,
		Street:    res.Street,
		Lat:       res.Lat,
		Long:      res.Long,
	}}, nil
}

func (h *AddressHandler) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.AddressResponse, error) {
	address, err := h.service.Delete(ctx, req.Id, &dto.DeleteAddressReq{})
	if err != nil {
		logger.Error("Failed to delete address: ", err)
		return nil, err
	}

	var res dto.Address
	utils.Copy(&res, &address)
	_ = h.cache.RemovePattern("*address*")
	return &pb.AddressResponse{Address: &pb.Address{
		IdAddress: req.Request.Id,
		IdUser:    req.Request.IdUser,
	}}, nil
}
