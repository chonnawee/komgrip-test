package entities

import "time"

type Beers struct {
	ID           uint64    `gorm:"primaryKey;comment:รหัสชื่อเบียร์"`
	BeerTypeName string    `gorm:"column:beer_type_name;not null;comment:ชื่อประเภทเบียร์"`
	BeerName     string    `gorm:"column:beer_name;not null;comment:ชื่อเบียร์"`
	BeerDesc     string    `gorm:"column:beer_desc;not null;comment:รายละเอียดเบียร์"`
	BeerImgPath  string    `gorm:"column:beer_img_path;not null;comment:path รูปเบียร์"`
	CreatedAt    time.Time `gorm:"column:created_at;comment:วันที่สร้างข้อมูล"`
	UpdatedAt    time.Time `gorm:"column:updated_at;comment:วันที่แก้ไขข้อมูล"`
}
