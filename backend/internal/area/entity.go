package area

import (
	"time"

	"gorm.io/gorm"
)

const StatusActive = "active"

type Area struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:64;not null;uniqueIndex" json:"name"`
	Slug      string         `gorm:"size:64;not null;uniqueIndex" json:"slug"`
	SortOrder int            `gorm:"not null;default:0;index" json:"sort_order"`
	Status    string         `gorm:"size:16;not null;default:active;index" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Area) TableName() string {
	return "areas"
}

type AreaResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (a Area) ToResponse() AreaResponse {
	return AreaResponse{
		ID:   a.ID,
		Name: a.Name,
		Slug: a.Slug,
	}
}
