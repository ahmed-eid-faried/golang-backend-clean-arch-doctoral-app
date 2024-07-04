package model

import (
	"time"

	"github.com/google/uuid"
)

// Doctor represents the domain model for an Doctor.
type Doctor struct {
	ID         string     `json:"id_Doctor"`
	IDUser     string     `json:"id_user"`
	Name       string     `json:"name"`
	Image      string     `json:"image"`
	Price      float32    `json:"price"`
	Specalist  string     `json:"specalist"`
	Experience int        `json:"experience"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
}

func (m *Doctor) BeforeCreate() error {
	m.ID = uuid.New().String()
	m.CreatedAt = time.Now()
	return nil
}

// func (m *Doctor) BeforeUpdate(tx *gorm.DB) error {
// 	m.UpdatedAt = time.Now()
// 	return nil
// }
