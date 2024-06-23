package dto

import (
	"main/pkg/paging"
)

// ***************************************************************************\\
// ***************************************************************************\\
// Address DTO represents the structure of address data transfer object.
// swagger:model Address
type Address struct {
	// ID of the address
	// example: "12345"
	ID string `json:"id_address"`
	// User ID associated with the address
	// example: "67890"
	IDUser string `json:"id_user"`
	// Name of the address
	// example: "Home"
	Name string `json:"name"`
	// City of the address
	// example: "San Francisco"
	City string `json:"city"`
	// Street of the address
	// example: "Market Street"
	Street string `json:"street"`
	// Latitude of the address
	// example: "37.7749"
	Lat string `json:"lat"`
	// Longitude of the address
	// example: "-122.4194"
	Long string `json:"long"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// CreateAddressReq represents the request body for creating a new address.
// swagger:model CreateAddressReq
type CreateAddressReq struct {
	// User ID associated with the address
	// example: "67890"
	IDUser string `json:"id_user"`
	// Name of the address
	// example: "Home"
	Name string `json:"name"`
	// City of the address
	// example: "San Francisco"
	City string `json:"city"`
	// Street of the address
	// example: "Market Street"
	Street string `json:"street"`
	// Latitude of the address
	// example: "37.7749"
	Lat string `json:"lat"`
	// Longitude of the address
	// example: "-122.4194"
	Long string `json:"long"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// UpdateAddressReq represents the request body for updating an existing address.
// swagger:model UpdateAddressReq
type UpdateAddressReq struct {
	// ID of the address
	// example: "12345"
	ID string `json:"id"`
	// User ID associated with the address
	// example: "67890"
	IDUser string `json:"id_user"`
	// Name of the address
	// example: "Home"
	Name string `json:"name"`
	// City of the address
	// example: "San Francisco"
	City string `json:"city"`
	// Street of the address
	// example: "Market Street"
	Street string `json:"street"`
	// Latitude of the address
	// example: "37.7749"
	Lat string `json:"lat"`
	// Longitude of the address
	// example: "-122.4194"
	Long string `json:"long"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// ListAddressReq represents the query parameters for listing addresses.
// swagger:model ListAddressReq
type ListAddressReq struct {
	// Name of the address
	// example: "Home"
	Name string `json:"name,omitempty" form:"name"`
	// User ID associated with the address
	// example: "67890"
	IDUser string `json:"id_user"`
	// Page number for pagination
	// example: 1
	Page int64 `json:"-" form:"page"`
	// Limit number of items per page
	// example: 10
	Limit int64 `json:"-" form:"limit"`
}

// ListAddressRes represents the response body for listing addresses.
// swagger:model ListAddressRes
type ListAddressRes struct {
	// List of addresses
	// example: [{"id_address":"12345","id_user":"67890","name":"Home","city":"San Francisco","street":"Market Street","lat":"37.7749","long":"-122.4194"}]
	Addresses []*Address `json:"addresses"`
	// Pagination info
	Pagination *paging.Pagination `json:"pagination"`
}

// ***************************************************************************\\
// ***************************************************************************\\
// DeleteAddressReq represents the request body for deleting an address.
// swagger:model DeleteAddressReq
type DeleteAddressReq struct {
	// ID of the address
	// example: "12345"
	ID string `json:"id"`
	// User ID associated with the address
	// example: "67890"
	IDUser string `json:"id_user"`
}

//***************************************************************************\\
//***************************************************************************\\
