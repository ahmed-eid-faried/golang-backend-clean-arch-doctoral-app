package repository

import (
	"context"

	"main/internal/address/dto"
	"main/internal/address/model"
	"main/pkg/config"
	"main/pkg/dbs"
	"main/pkg/paging"
)

//go:generate mockery --name=IAddressRepository
type IAddressRepository interface {
	Create(ctx context.Context, Address *model.Address) error
	Delete(ctx context.Context, Address *model.Address) error
	Update(ctx context.Context, Address *model.Address) error
	ListAddresses(ctx context.Context, req *dto.ListAddressReq) ([]*model.Address, *paging.Pagination, error)
	GetAddressByID(ctx context.Context, id string) (*model.Address, error)
}

type AddressRepo struct {
	db dbs.IDatabase
}

// ListAddresses implements IAddressRepository.
// func (r *AddressRepo) ListAddresses(ctx context.Context, req *dto.ListAddressReq) ([]*model.Address, *paging.Pagination, error) {
// 	panic("unimplemented")
// }

func NewAddressRepository(db dbs.IDatabase) *AddressRepo {
	return &AddressRepo{db: db}
}

func (r *AddressRepo) ListAddresses(ctx context.Context, req *dto.ListAddressReq) ([]*model.Address, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := make([]dbs.Query, 0)
	if req.Name != "" {
		query = append(query, dbs.NewQuery("name LIKE ?", "%"+req.Name+"%"))
	}
	// if req.Code != "" {
	// 	query = append(query, dbs.NewQuery("code = ?", req.Code))
	// }

	// order := "created_at"
	// if req.OrderBy != "" {
	// 	order = req.OrderBy
	// 	if req.OrderDesc {
	// 		order += " DESC"
	// 	}
	// }

	var total int64
	if err := r.db.Count(ctx, &model.Address{}, &total, dbs.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.New(req.Page, req.Limit, total)

	var Addresss []*model.Address
	if err := r.db.Find(
		ctx,
		&Addresss,
		dbs.WithQuery(query...),
		dbs.WithLimit(int(pagination.Limit)),
		dbs.WithOffset(int(pagination.Skip)),
		// dbs.WithOrder(order),
	); err != nil {
		return nil, nil, err
	}

	return Addresss, pagination, nil
}

func (r *AddressRepo) GetAddressByID(ctx context.Context, id string) (*model.Address, error) {
	var Address model.Address
	if err := r.db.FindById(ctx, id, &Address); err != nil {
		return nil, err
	}
	return &Address, nil
}

func (r *AddressRepo) Create(ctx context.Context, Address *model.Address) error {
	return r.db.Create(ctx, Address)
}

func (r *AddressRepo) Update(ctx context.Context, Address *model.Address) error {
	return r.db.Update(ctx, Address)
}
func (r *AddressRepo) Delete(ctx context.Context, Address *model.Address) error {
	return r.db.Delete(ctx, Address)
}
