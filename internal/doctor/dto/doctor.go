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
type ListDoctorReq struct {
	// Name of the Doctor
	// example: "Home"
	Name string `json:"name,omitempty" form:"name"`
	// User ID associated with the Doctor
	// example: "67890"
	IDUser string `json:"id_user"`
	// Page number for pagination
	// example: 1
	Page int64 `json:"-" form:"page"`
	// Limit number of items per page
	// example: 10
	Limit int64 `json:"-" form:"limit"`
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
