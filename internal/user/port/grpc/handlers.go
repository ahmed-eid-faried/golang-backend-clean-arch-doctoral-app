package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"

	"main/internal/user/dto"
	"main/internal/user/model"
	"main/internal/user/service"
	"main/pkg/redis"
	"main/pkg/utils"
	pb "main/proto/gen/go/user"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	cache         redis.IRedis
	service       service.IUserService
	fbOauthConfig *oauth2.Config
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	user, accessToken, refreshToken, err := h.service.Login(ctx, &dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	var res pb.LoginRes
	utils.Copy(&res.User, &user)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	return &res, nil
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	protoRole, err := ConvertProtoToModelUserRole(req.Role)
	if err != nil {
		logger.Error("Failed to convert user role: ", err)
	}
	user, err := h.service.Register(ctx, &dto.RegisterReq{
		Email:       req.Email,
		Password:    req.Password,
		Name:        req.Name,
		Role:        protoRole,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	var res pb.RegisterRes
	utils.Copy(&res.User, &user)
	return &res, nil
}

func (h *UserHandler) GetMe(ctx context.Context, _ *pb.GetMeReq) (*pb.GetMeRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	var res pb.GetMeRes
	utils.Copy(&res.User, &user)
	return &res, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenReq) (*pb.RefreshTokenRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	accessToken, err := h.service.RefreshToken(ctx, userID)
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	res := pb.RefreshTokenRes{
		AccessToken: accessToken,
	}
	return &res, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	protoRole, err2 := ConvertProtoToModelUserRole(req.Role)
	if err2 != nil {
		logger.Error("Failed to convert user role: ", err2)
	}

	err := h.service.UpdateUser(ctx, userID, &dto.UpdateUserReq{
		Password:    req.Password,
		NewPassword: req.NewPassword,
		Email:       req.Email,
		ID:          req.Id,
		Name:        req.Name,
		Role:        protoRole,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	return &pb.UpdateUserRes{}, nil
}

func (h *UserHandler) VerfiyCodePhoneNumber(ctx context.Context, req *pb.VerifyPhoneNumberRequest) (*pb.VerifyResponse, error) {

	_, err := h.service.VerifyPhoneNumber(ctx, dto.VerifyPhoneNumberRequest{
		PhoneNumber:           req.PhoneNumber,
		VerifyCodePhoneNumber: req.VerifyCodePhoneNumber})

	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err

	}

	return &pb.VerifyResponse{}, nil
}
func (h *UserHandler) VerfiyCodeEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyResponse, error) {

	_, err := h.service.VerifyPhoneNumber(ctx, dto.VerifyPhoneNumberRequest{
		PhoneNumber:           req.Email,
		VerifyCodePhoneNumber: req.VerifyCodeEmail})

	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err

	}

	return &pb.VerifyResponse{}, nil
}

func (h *UserHandler) VerfiyCodePhoneNumberResend(ctx context.Context, req *pb.ResendVerifyPhoneNumberRequest) (*pb.VerifyResponse, error) {

	_, err := h.service.VerifyPhoneNumber(ctx, dto.VerifyPhoneNumberRequest{
		PhoneNumber:           req.PhoneNumber,
		VerifyCodePhoneNumber: req.PhoneNumber})

	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err

	}

	return &pb.VerifyResponse{}, nil
}

func (h *UserHandler) VerfiyCodeEmailResend(ctx context.Context, req *pb.ResendVerifyEmailRequest) (*pb.VerifyResponse, error) {

	_, err := h.service.VerifyPhoneNumber(ctx, dto.VerifyPhoneNumberRequest{
		PhoneNumber:           req.Email,
		VerifyCodePhoneNumber: req.Email})

	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err

	}

	return &pb.VerifyResponse{}, nil
}

// ConvertModelUserRoleToProto converts a model.UserRole to pb.UserRole
func ConvertModelUserRoleToProto(role model.UserRole) (pb.UserRole, error) {
	switch role {
	case model.UserRoleClient:
		return pb.UserRole_ROLE_USER, nil
	case model.UserRoleDoctor:
		return pb.UserRole_ROLE_DOCTOR, nil
	case model.UserRoleAdmin:
		return pb.UserRole_ROLE_ADMIN, nil
	default:
		return pb.UserRole_ROLE_UNKNOWN, fmt.Errorf("unknown role: %v", role)
	}
}

// ConvertProtoToModelUserRole converts a pb.UserRole to model.UserRole
func ConvertProtoToModelUserRole(role pb.UserRole) (model.UserRole, error) {
	switch role {
	case pb.UserRole_ROLE_USER:
		return model.UserRoleClient, nil
	case pb.UserRole_ROLE_DOCTOR:
		return model.UserRoleDoctor, nil
	case pb.UserRole_ROLE_ADMIN:
		return model.UserRoleAdmin, nil
	default:
		return model.UserRoleClient, fmt.Errorf("unknown role: %v", role)
	}
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var res dto.ListUsersRes
	cacheKey := "users_list"
	err := h.cache.Get(cacheKey, &res)
	if err == nil {
		var pbUsers []*pb.User
		for _, addr := range res.Users {
			protoRole, err := ConvertModelUserRoleToProto(addr.Role)
			if err != nil {
				logger.Error("Failed to convert user role: ", err)
				continue // or handle the error as needed
			}
			var deletedAtProto *timestamppb.Timestamp
			if addr.DeletedAt != nil {
				deletedAtProto = timestamppb.New(*addr.DeletedAt)
			}
			pbUsers = append(pbUsers, &pb.User{
				Id:                    addr.ID,
				CreatedAt:             timestamppb.New(addr.CreatedAt),
				UpdatedAt:             timestamppb.New(addr.UpdatedAt),
				DeletedAt:             deletedAtProto,
				Password:              addr.Password,
				Role:                  protoRole,
				Name:                  addr.Name,
				Email:                 addr.Email,
				PhoneNumber:           addr.PhoneNumber,
				VerifyCodeEmail:       int32(addr.VerifyCodeEmail),
				VerifyCodePhoneNumber: int32(addr.VerifyCodePhoneNumber),
				ApproveEmail:          addr.ApproveEmail,
				ApprovePhoneNumber:    addr.ApprovePhoneNumber,
			})
		}
		return &pb.ListUsersResponse{Users: pbUsers}, nil
	}

	Users, pagination, err := h.service.ListUsers(ctx, dto.ListUsersReq{})
	if err != nil {
		logger.Error("Failed to get list of Users: ", err)
		return nil, err
	}

	utils.Copy(&res.Users, &Users)
	res.Pagination = pagination
	_ = h.cache.SetWithExpiration(cacheKey, res, time.Hour) // Adjust caching time as needed

	var pbUsers []*pb.User
	for _, addr := range Users {
		protoRole, err := ConvertModelUserRoleToProto(addr.Role)
		if err != nil {
			logger.Error("Failed to convert user role: ", err)
			continue // or handle the error as needed
		}
		var deletedAtProto *timestamppb.Timestamp
		if addr.DeletedAt != nil {
			deletedAtProto = timestamppb.New(*addr.DeletedAt)
		}
		pbUsers = append(pbUsers, &pb.User{
			Id:                    addr.ID,
			CreatedAt:             timestamppb.New(addr.CreatedAt),
			UpdatedAt:             timestamppb.New(addr.UpdatedAt),
			DeletedAt:             deletedAtProto,
			Password:              addr.Password,
			Role:                  protoRole,
			Name:                  addr.Name,
			Email:                 addr.Email,
			PhoneNumber:           addr.PhoneNumber,
			VerifyCodeEmail:       int32(addr.VerifyCodeEmail),
			VerifyCodePhoneNumber: int32(addr.VerifyCodePhoneNumber),
			ApproveEmail:          addr.ApproveEmail,
			ApprovePhoneNumber:    addr.ApprovePhoneNumber,
		})
	}
	return &pb.ListUsersResponse{Users: pbUsers}, nil
}
