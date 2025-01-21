package adapters

import (
	"komgrip-test/entities"
	"komgrip-test/usecases"

	"gorm.io/gorm"
)

type beersRepositoryDB struct {
	db *gorm.DB
}

func NewBeersRepositoryDB(db *gorm.DB) usecases.BeersRepository {
	return &beersRepositoryDB{db: db}
}

func (r *beersRepositoryDB) Create(beer entities.Beers) error {
	err := r.db.Create(&beer).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *beersRepositoryDB) GetDatas(params usecases.GetDatasParams) (beers []entities.Beers, err error) {
	query := r.db.Model(&entities.Beers{})
	if params.BeerName != "" {
		query = query.Where("beer_name like ?", "%"+params.BeerName+"%")
	}
	if params.Limit != 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}
	err = query.Find(&beers).Error
	if err != nil {
		return nil, err
	}
	return beers, nil
}

func (r *beersRepositoryDB) GetData(id uint64) (beer *entities.Beers, err error) {
	err = r.db.Where("id = ?", id).First(&beer).Error
	if err != nil {
		return nil, err
	}
	return beer, nil
}

func (r *beersRepositoryDB) Update(id uint64, beer entities.Beers) error {
	err := r.db.Where("id = ?", id).Updates(&beer).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *beersRepositoryDB) Delete(id uint64) error {
	err := r.db.Where("id = ?", id).Delete(&entities.Beers{}).Error
	if err != nil {
		return err
	}
	return nil
}
