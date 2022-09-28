package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Server struct {
	ID        uuid.UUID `json:"id" validate:"required,uuid" `
	Name      *string   `gorm:"uniqueIndex" json:"name" validate:"required"`
	Status    *string   `json:"status" validate:"required"`
	Ipv4      *string   `json:"ipvd4" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

func (s Server) value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Server) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &s)
}

func MigrateServer(db *gorm.DB) error {
	err := db.AutoMigrate(&Server{})
	return err
}
