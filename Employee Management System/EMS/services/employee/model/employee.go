
package model

import "time"

type Employee struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	DepartmentID uint      `gorm:"not null;index" json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Employee) TableName() string {
	return "employees"
}
