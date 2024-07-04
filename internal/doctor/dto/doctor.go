package dto

import (
	"main/pkg/paging"
)

// ***************************************************************************\\
// ***************************************************************************\\
// Address DTO represents the structure of address data transfer object.
// swagger:model Doctor
type Doctor struct {
	ID         string  `json:"id_Doctor"`
	IDUser     string  `json:"id_user"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Price      float32 `json:"price"`
	Specalist  string  `json:"specalist"`
	Experience int     `json:"experience"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// CreateDoctorReq represents the request body for creating a new Doctor.
// swagger:model CreateDoctorReq
type CreateDoctorReq struct {
	IDUser     string  `json:"id_user"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Price      float32 `json:"price"`
	Specalist  string  `json:"specalist"`
	Experience int     `json:"experience"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// UpdateDoctorReq represents the request body for updating an existing Doctor.
// swagger:model UpdateDoctorReq
type UpdateDoctorReq struct {
	ID         string  `json:"id_Doctor"`
	IDUser     string  `json:"id_user"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Price      float32 `json:"price"`
	Specalist  string  `json:"specalist"`
	Experience int     `json:"experience"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// ListDoctorReq represents the query parameters for listing Doctors.
// swagger:model ListDoctorReq
// type OrderBy string

// const (
// 	OrderByNone OrderBy = ""
// 	IDUser      OrderBy = "id_user"
// 	Name        OrderBy = "name"
// 	Price       OrderBy = "price"
// 	Specialty   OrderBy = "specalist"
// 	Experience  OrderBy = "experience"
// 	CreatedAt   OrderBy = "created_at"
// )

type ListDoctorReq struct {
	Search    string    `json:"search,omitempty" form:"search"`
	IDUser    string    `json:"id_user,omitempty" form:"id_user"`
	Page      int64     `json:"page,omitempty" form:"page"`
	Limit     int64     `json:"limit,omitempty" form:"limit"`
	OrderList []OrderBy `json:"order_list,omitempty" form:"order_list"`
}
type OrderBy struct {
	OrderBy   string `json:"order_by,omitempty" form:"order_by"`
	OrderDesc bool   `json:"order_desc,omitempty" form:"order_desc"`
}

// ListDoctorRes represents the response body for listing Doctors.
// swagger:model ListDoctorRes
type ListDoctorRes struct {
	// List of Doctors
	// example: [{"id_Doctor":"12345","id_user":"67890","name":"Home","city":"San Francisco","street":"Market Street","lat":"37.7749","long":"-122.4194"}]
	Doctors []*Doctor `json:"Doctors"`
	// Pagination info
	Pagination *paging.Pagination `json:"pagination"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// DeleteDoctorReq represents the request body for deleting an Doctor.
// swagger:model DeleteDoctorReq
type DeleteDoctorReq struct {
	// ID of the Doctor
	// example: "12345"
	ID string `json:"id"`
	// User ID associated with the Doctor
	// example: "67890"
	IDUser string `json:"id_user"`
}

//***************************************************************************\\
//***************************************************************************\\
