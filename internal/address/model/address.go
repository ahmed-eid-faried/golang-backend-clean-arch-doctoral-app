package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Address represents the domain model for an address.
type Address struct {
	ID        string    `json:"id_address"`
	IDUser    string    `json:"id_user"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	Street    string    `json:"street"`
	Lat       string    `json:"lat"`
	Long      string    `json:"long"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Address) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New().String()
	m.CreatedAt = time.Now()
	return nil
}

func (m *Address) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
