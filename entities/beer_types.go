package entities

import "time"

type BeerTypes struct {
	ID           uint64    `gorm:"primaryKey;comment:รหัสประเภทเบียร์"`
	BeerTypeName string    `gorm:"column:beer_type_name;not null;comment:ชื่อประเภทเบียร์"`
	CreatedAt    time.Time `gorm:"column:created_at;comment:วันที่สร้างข้อมูล"`
	UpdatedAt    time.Time `gorm:"column:updated_at;comment:วันที่แก้ไขข้อมูล"`
}
