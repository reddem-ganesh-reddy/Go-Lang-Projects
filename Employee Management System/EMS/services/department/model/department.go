package model

import "time"

// Department represents a department in the organization.
type Department struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null;unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the database table name.
func (Department) TableName() string {
	return "departments"
}