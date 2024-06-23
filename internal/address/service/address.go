package service

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"main/internal/address/dto"
	"main/internal/address/model"
	"main/internal/address/repository"
	"main/pkg/paging"
	"main/pkg/utils"
)

//go:generate mockery --name=IAddressService
type IAddressService interface {
	ListAddresses(c context.Context, req *dto.ListAddressReq) ([]*model.Address, *paging.Pagination, error)
	GetAddressByID(ctx context.Context, id string) (*model.Address, error)
	Create(ctx context.Context, req *dto.CreateAddressReq) (*model.Address, error)
	Delete(ctx context.Context, id string, req *dto.DeleteAddressReq) (*model.Address, error)
	Update(ctx context.Context, id string, req *dto.UpdateAddressReq) (*model.Address, error)
}

type AddressService struct {
	validator validation.Validation
	repo      repository.IAddressRepository
}

func NewAddressService(
	validator validation.Validation,
	repo repository.IAddressRepository,
) *AddressService {
	return &AddressService{
		validator: validator,
		repo:      repo,
	}
}

func (p *AddressService) GetAddressByID(ctx context.Context, id string) (*model.Address, error) {
	Address, err := p.repo.GetAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return Address, nil
}

func (p *AddressService) ListAddresses(ctx context.Context, req *dto.ListAddressReq) ([]*model.Address, *paging.Pagination, error) {
	Addresss, pagination, err := p.repo.ListAddresses(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return Addresss, pagination, nil
}

func (p *AddressService) Create(ctx context.Context, req *dto.CreateAddressReq) (*model.Address, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	var Address model.Address
	utils.Copy(&Address, req)

	err := p.repo.Create(ctx, &Address)
	if err != nil {
		logger.Errorf("Create fail, error: %s", err)
		return nil, err
	}

	return &Address, nil
}

func (p *AddressService) Update(ctx context.Context, id string, req *dto.UpdateAddressReq) (*model.Address, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	Address, err := p.repo.GetAddressByID(ctx, id)
	if err != nil {
		logger.Errorf("Update.GetAddressByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	utils.Copy(Address, req)
	err = p.repo.Update(ctx, Address)
	if err != nil {
		logger.Errorf("Update fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return Address, nil
}

func (p *AddressService) Delete(ctx context.Context, id string, req *dto.DeleteAddressReq) (*model.Address, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	Address, err := p.repo.GetAddressByID(ctx, id)
	if err != nil {
		logger.Errorf("Delete.GetAddressByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	utils.Copy(Address, req)
	err = p.repo.Delete(ctx, Address)
	if err != nil {
		logger.Errorf("Delete fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return Address, nil
}
