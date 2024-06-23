package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"main/internal/user/dto"
	"main/internal/user/model"
	"main/pkg/config"
	"main/pkg/dbs"
	"main/pkg/paging"
)

//go:generate mockery --name=IUserRepository
type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindByEmailAndVerifyCode(ctx context.Context, email, verifyCode string) (*model.User, error)
	FindByPhoneAndVerifyCode(ctx context.Context, PhoneNumber, verifyCode string) (*model.User, error)
	FindByPhone(ctx context.Context, PhoneNumber string) (*model.User, error)
	FindByEmail(ctx context.Context, PhoneNumber string) (*model.User, error)
	ListUsers(ctx context.Context, req dto.ListUsersReq) ([]*model.User, *paging.Pagination, error)
	UpdatePhone(ctx context.Context, user *model.User) error
	UpdateEmail(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, User *model.User) error
	FindOrCreateByGoogleID(ctx context.Context, googleID, email, name string) (*model.User, error)
	FindOrCreateByFacebookID(ctx context.Context, facebookID, email, name string) (*model.User, error)
}

type UserRepo struct {
	db dbs.IDatabase
}

// repo repository.IUserRepository
func NewUserRepository(db dbs.IDatabase) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.Create(ctx, user)
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.Update(ctx, user)
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.FindById(ctx, id, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := dbs.NewQuery("email = ?", email)
	if err := r.db.FindOne(ctx, &user, dbs.WithQuery(query)); err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *UserRepo) FindByEmailAndVerifyCode(ctx context.Context, email, verifyCode string) (*model.User, error) {
	var user model.User
	query := "SELECT id, email, verify_code_email, approve FROM users WHERE email = ? AND verify_code_email = ?"

	result := r.db.GetDB().WithContext(ctx).Raw(query, email, verifyCode).Scan(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := "SELECT id, email, verify_code_email, approve FROM users WHERE email = ?"

	result := r.db.GetDB().WithContext(ctx).Raw(query, email).Scan(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
func (r *UserRepo) UpdateEmail(ctx context.Context, user *model.User) error {
	query := "UPDATE users SET approve = ? WHERE email = ?"
	err := r.db.Exec(ctx, query, user.ApproveEmail, user.Email)
	return err
}
func (r *UserRepo) FindByPhoneAndVerifyCode(ctx context.Context, PhoneNumber, verifyCode string) (*model.User, error) {
	var user model.User
	query := "SELECT id, phone_number, verify_code_phone_number, approve FROM users WHERE phone_number = ? AND verify_code_phone_number = ?"

	result := r.db.GetDB().WithContext(ctx).Raw(query, PhoneNumber, verifyCode).Scan(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepo) FindByPhone(ctx context.Context, PhoneNumber string) (*model.User, error) {
	var user model.User
	query := "SELECT id, phone_number, verify_code_phone_number, approve FROM users WHERE phone_number = ? "

	result := r.db.GetDB().WithContext(ctx).Raw(query, PhoneNumber).Scan(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepo) UpdatePhone(ctx context.Context, user *model.User) error {
	query := "UPDATE users SET approve = ? WHERE phone_number = ?"
	err := r.db.Exec(ctx, query, user.ApprovePhoneNumber, user.PhoneNumber)
	return err
}

func (r *UserRepo) ListUsers(ctx context.Context, req dto.ListUsersReq) ([]*model.User, *paging.Pagination, error) {
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
	if err := r.db.Count(ctx, &model.User{}, &total, dbs.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.New(req.Page, req.Limit, total)

	var Users []*model.User
	if err := r.db.Find(
		ctx,
		&Users,
		dbs.WithQuery(query...),
		dbs.WithLimit(int(pagination.Limit)),
		dbs.WithOffset(int(pagination.Skip)),
		// dbs.WithOrder(order),
	); err != nil {
		return nil, nil, err
	}

	return Users, pagination, nil
}
func (r *UserRepo) Delete(ctx context.Context, User *model.User) error {
	return r.db.Delete(ctx, User)
}

func (r *UserRepo) FindOrCreateByGoogleID(ctx context.Context, googleID, email, name string) (*model.User, error) {
	var user model.User
	err := r.db.GetDB().WithContext(ctx).Where("id = ?", googleID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = model.User{
			ID:    googleID,
			Email: email,
			Name:  name,
		}
		if err := r.db.Create(ctx, &user); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindOrCreateByFacebookID(ctx context.Context, facebookID, email, name string) (*model.User, error) {
	var user model.User
	err := r.db.GetDB().WithContext(ctx).Where("id = ?", facebookID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = model.User{
			ID:    facebookID,
			Email: email,
			Name:  name,
		}
		if err := r.db.Create(ctx, &user); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
