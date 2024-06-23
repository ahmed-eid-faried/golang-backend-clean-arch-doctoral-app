package service

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"main/internal/doctor/dto"
	"main/internal/doctor/model"
	"main/internal/doctor/repository"
	"main/pkg/paging"
	"main/pkg/utils"
)

//go:generate mockery --name=IDoctorService
type IDoctorService interface {
	ListDoctors(c context.Context, req *dto.ListDoctorReq) ([]*model.Doctor, *paging.Pagination, error)
	GetDoctorByID(ctx context.Context, id string) (*model.Doctor, error)
	Create(ctx context.Context, req *dto.CreateDoctorReq) (*model.Doctor, error)
	Delete(ctx context.Context, id string, req *dto.DeleteDoctorReq) (*model.Doctor, error)
	Update(ctx context.Context, id string, req *dto.UpdateDoctorReq) (*model.Doctor, error)
}

type DoctorService struct {
	validator validation.Validation
	repo      repository.IDoctorRepository
}

func NewDoctorService(
	validator validation.Validation,
	repo repository.IDoctorRepository,
) *DoctorService {
	return &DoctorService{
		validator: validator,
		repo:      repo,
	}
}

func (p *DoctorService) GetDoctorByID(ctx context.Context, id string) (*model.Doctor, error) {
	Doctor, err := p.repo.GetDoctorByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return Doctor, nil
}

func (p *DoctorService) ListDoctors(ctx context.Context, req *dto.ListDoctorReq) ([]*model.Doctor, *paging.Pagination, error) {
	Doctors, pagination, err := p.repo.ListDoctors(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return Doctors, pagination, nil
}

func (p *DoctorService) Create(ctx context.Context, req *dto.CreateDoctorReq) (*model.Doctor, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	var doctor model.Doctor
	utils.Copy(&doctor, req)
	doctor.BeforeCreate()
	err := p.repo.Create(ctx, &doctor)
	if err != nil {
		logger.Errorf("Create fail, error: %s", err)
		return nil, err
	}

	return &doctor, nil
}

func (p *DoctorService) Update(ctx context.Context, id string, req *dto.UpdateDoctorReq) (*model.Doctor, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	Doctor, err := p.repo.GetDoctorByID(ctx, id)
	if err != nil {
		logger.Errorf("Update.GetDoctorByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	utils.Copy(Doctor, req)
	err = p.repo.Update(ctx, Doctor)
	if err != nil {
		logger.Errorf("Update fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return Doctor, nil
}

func (p *DoctorService) Delete(ctx context.Context, id string, req *dto.DeleteDoctorReq) (*model.Doctor, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	Doctor, err := p.repo.GetDoctorByID(ctx, id)
	if err != nil {
		logger.Errorf("Delete.GetDoctorByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	utils.Copy(Doctor, req)
	err = p.repo.Delete(ctx, Doctor)
	if err != nil {
		logger.Errorf("Delete fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return Doctor, nil
}
