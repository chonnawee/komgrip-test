package entities

type Beers struct {
	ID          uint64    `gorm:"primaryKey;comment:รหัสชื่อเบียร์"`
	BeerTypeID  uint64    `gorm:"column:beer_type_id;not null;comment:รหัสประเภทเบียร์"`
	BeerName    string    `gorm:"column:beer_name;not null;comment:ชื่อเบียร์"`
	BeerDesc    string    `gorm:"column:beer_desc;not null;comment:รายละเอียดเบียร์"`
	BeerImgPath string    `gorm:"column:beer_img_path;not null;comment:path รูปเบียร์"`
	BeerType    BeerTypes `gorm:"foreignKey:BeerTypeID"`
}
