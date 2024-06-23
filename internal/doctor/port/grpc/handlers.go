package grpc

import (
	"context"
	"time"

	"github.com/quangdangfit/gocommon/logger"

	"main/internal/doctor/dto"
	"main/internal/doctor/service"
	"main/pkg/config"
	"main/pkg/redis"
	"main/pkg/utils"
	pb "main/proto/gen/go/doctor"
)

type DoctorHandler struct {
	cache   redis.IRedis
	service service.IDoctorService
	pb.UnimplementedDoctorServiceServer
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

// // Convert time.Time to string
//
//	func formatTimeToString(t time.Time) string {
//		return t.Format(time.RFC3339)
//	}

func (h *DoctorHandler) GetDoctorByID(ctx context.Context, req *pb.GetDoctorByIDRequest) (*pb.DoctorResponse, error) {
	var res dto.Doctor
	cacheKey := "Doctor_" + req.Id
	err := h.cache.Get(cacheKey, &res)
	if err == nil {
		return &pb.DoctorResponse{Doctor: &pb.Doctor{
			Id:         res.ID,
			IdUser:     res.IDUser,
			Name:       res.Name,
			Image:      res.Image,
			Price:      res.Price,
			Specialist: res.Specalist,
			Experience: int32(res.Experience),
		}}, nil
	}

	Doctor, err := h.service.GetDoctorByID(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to get Doctor detail: ", err)
		return nil, err
	}

	utils.Copy(&res, &Doctor)
	_ = h.cache.SetWithExpiration(cacheKey, res, config.DoctorCachingTime.Abs())
	return &pb.DoctorResponse{Doctor: &pb.Doctor{
		Id:         res.ID,
		IdUser:     res.IDUser,
		Name:       res.Name,
		Image:      res.Image,
		Price:      res.Price,
		Specialist: res.Specalist,
		Experience: int32(res.Experience),
	}}, nil
}

func (h *DoctorHandler) ListDoctors(ctx context.Context, req *pb.ListDoctorReq) (*pb.ListDoctorRes, error) {
	var res dto.ListDoctorRes
	cacheKey := "Doctors_list"
	err := h.cache.Get(cacheKey, &res)
	if err == nil {
		var pbDoctors []*pb.Doctor
		for _, addr := range res.Doctors {
			pbDoctors = append(pbDoctors, &pb.Doctor{
				Id:         addr.ID,
				IdUser:     addr.IDUser,
				Name:       addr.Name,
				Image:      addr.Image,
				Price:      addr.Price,
				Specialist: addr.Specalist,
				Experience: int32(addr.Experience),
			})
		}
		return &pb.ListDoctorRes{Doctors: pbDoctors}, nil
	}

	Doctors, pagination, err := h.service.ListDoctors(ctx, &dto.ListDoctorReq{})
	if err != nil {
		logger.Error("Failed to get list of Doctors: ", err)
		return nil, err
	}

	utils.Copy(&res.Doctors, &Doctors)
	res.Pagination = pagination
	_ = h.cache.SetWithExpiration(cacheKey, res, time.Hour) // Adjust caching time as needed

	var pbDoctors []*pb.Doctor
	for _, addr := range Doctors {
		pbDoctors = append(pbDoctors, &pb.Doctor{
			Id:         addr.ID,
			IdUser:     addr.IDUser,
			Name:       addr.Name,
			Image:      addr.Image,
			Price:      addr.Price,
			Specialist: addr.Specalist,
			Experience: int32(addr.Experience),
		})
	}
	return &pb.ListDoctorRes{Doctors: pbDoctors}, nil
}

func (h *DoctorHandler) CreateDoctor(ctx context.Context, req *pb.CreateDoctorReq) (*pb.DoctorResponse, error) {
	var DoctorDTO dto.CreateDoctorReq
	DoctorDTO.IDUser = req.IdUser
	DoctorDTO.Name = req.Name
	DoctorDTO.Image = req.Image
	DoctorDTO.Price = req.Price
	DoctorDTO.Specalist = req.Specialist
	DoctorDTO.Experience = int(req.Experience)

	Doctor, err := h.service.Create(ctx, &DoctorDTO)
	if err != nil {
		logger.Error("Failed to create Doctor: ", err)
		return nil, err
	}

	var res dto.Doctor
	utils.Copy(&res, &Doctor)
	_ = h.cache.RemovePattern("*Doctor*")
	return &pb.DoctorResponse{Doctor: &pb.Doctor{
		Id:         res.ID,
		IdUser:     res.IDUser,
		Name:       res.Name,
		Image:      res.Image,
		Price:      res.Price,
		Specialist: res.Specalist,
		Experience: int32(res.Experience),
	}}, nil
}

func (h *DoctorHandler) UpdateDoctor(ctx context.Context, req *pb.UpdateDoctorReq) (*pb.DoctorResponse, error) {
	var DoctorDTO dto.UpdateDoctorReq
	DoctorDTO.IDUser = req.IdUser
	DoctorDTO.Name = req.Name
	DoctorDTO.Image = req.Image
	DoctorDTO.Price = req.Price
	DoctorDTO.Specalist = req.Specialist
	DoctorDTO.Experience = int(req.Experience)

	Doctor, err := h.service.Update(ctx, req.Id, &DoctorDTO)
	if err != nil {
		logger.Error("Failed to update Doctor: ", err)
		return nil, err
	}

	var res dto.Doctor
	utils.Copy(&res, &Doctor)
	_ = h.cache.RemovePattern("*Doctor*")
	return &pb.DoctorResponse{Doctor: &pb.Doctor{
		Id:         res.ID,
		IdUser:     res.IDUser,
		Name:       res.Name,
		Image:      res.Image,
		Price:      res.Price,
		Specialist: res.Specalist,
		Experience: int32(res.Experience),
	}}, nil
}

func (h *DoctorHandler) DeleteDoctor(ctx context.Context, req *pb.DeleteDoctorReq) (*pb.DoctorResponse, error) {
	Doctor, err := h.service.Delete(ctx, req.Id, &dto.DeleteDoctorReq{})
	if err != nil {
		logger.Error("Failed to delete Doctor: ", err)
		return nil, err
	}

	var res dto.Doctor
	utils.Copy(&res, &Doctor)
	_ = h.cache.RemovePattern("*Doctor*")
	return &pb.DoctorResponse{Doctor: &pb.Doctor{
		Id:         res.ID,
		IdUser:     res.IDUser,
		Name:       res.Name,
		Image:      res.Image,
		Price:      res.Price,
		Specialist: res.Specalist,
		Experience: int32(res.Experience),
	}}, nil
}
