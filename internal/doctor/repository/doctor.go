package repository

import (
	"context"
	"fmt"

	"main/internal/doctor/dto"
	"main/internal/doctor/model"
	"main/pkg/config"
	"main/pkg/dbs"
	"main/pkg/paging"
)

//go:generate mockery --name=IDoctorRepository
type IDoctorRepository interface {
	Create(ctx context.Context, Doctor *model.Doctor) error
	Delete(ctx context.Context, Doctor *model.Doctor) error
	Update(ctx context.Context, Doctor *model.Doctor) error
	ListDoctors(ctx context.Context, req *dto.ListDoctorReq) ([]*model.Doctor, *paging.Pagination, error)
	GetDoctorByID(ctx context.Context, id string) (*model.Doctor, error)
}

type DoctorRepo struct {
	db dbs.IDatabase
}

// ListDoctors implements IDoctorRepository.
// func (r *DoctorRepo) ListDoctors(ctx context.Context, req *dto.ListDoctorReq) ([]*model.Doctor, *paging.Pagination, error) {
// 	panic("unimplemented")
// }

func NewDoctorRepository(db dbs.IDatabase) *DoctorRepo {
	return &DoctorRepo{db: db}
}

func (r *DoctorRepo) ListDoctors(ctx context.Context, req *dto.ListDoctorReq) ([]*model.Doctor, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := make([]dbs.Query, 0)
	if req.Search != "" || req.Search != " " {
		query = append(query, dbs.NewQuery("name LIKE ?", "%"+req.Search+"%"))
	} else {
		query = append(query, dbs.NewQuery(""))
	}
	// if req.Code != "" {
	// 	query = append(query, dbs.NewQuery("code = ?", req.Code))
	// }

	// Construct order criteria
	orderList := req.OrderList
	order := ""
	OrderDes := ""
	for _, i := range orderList {
		if i.OrderDesc {
			OrderDes = " DESC"
		} else {
			OrderDes = ""

		}
		order += (i.OrderBy + OrderDes)

	}

	var total int64
	if err := r.db.Count(ctx, &model.Doctor{}, &total, dbs.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.New(req.Page, req.Limit, total)

	var Doctors []*model.Doctor
	if err := r.db.Find(
		ctx,
		&Doctors,
		dbs.WithQuery(query...),
		dbs.WithLimit(int(pagination.Limit)),
		dbs.WithOffset(int(pagination.Skip)),
		dbs.WithOrder(order),
	); err != nil {
		return nil, nil, err
	}
	print(Doctors)

	return Doctors, pagination, nil
}

func (r *DoctorRepo) GetDoctorByID(ctx context.Context, id string) (*model.Doctor, error) {
	var Doctor model.Doctor
	if err := r.db.FindById(ctx, id, &Doctor); err != nil {
		return nil, err
	}
	return &Doctor, nil
}

func (r *DoctorRepo) Create(ctx context.Context, doctor *model.Doctor) error {
	var existingDoctor model.Doctor
	query := dbs.NewQuery("id_user = ?", doctor.IDUser)
	if err := r.db.FindOne(ctx, &existingDoctor, dbs.WithQuery(query)).Error; err == nil {
		return fmt.Errorf("doctor with id_user %s already exists", doctor.IDUser)
	} else {
		// return fmt.Errorf("we have errors :-  %s", err)
		return r.db.Create(ctx, doctor)

	}
}

func (r *DoctorRepo) Update(ctx context.Context, Doctor *model.Doctor) error {
	return r.db.Update(ctx, Doctor)
}
func (r *DoctorRepo) Delete(ctx context.Context, Doctor *model.Doctor) error {
	return r.db.Delete(ctx, Doctor)
}
