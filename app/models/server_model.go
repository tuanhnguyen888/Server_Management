package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Server struct {
	ID        uuid.UUID `json:"id" validate:"required,uuid" `
	Name      *string   `gorm:"uniqueIndex" json:"name"`
	Status    bool      `json:"status"`
	Ipv4      *string   `json:"ipvd4" `
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"update_at"`
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
