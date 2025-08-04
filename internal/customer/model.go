package customer

import "time"

type Customer struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	Email     string    `json:"email"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}
