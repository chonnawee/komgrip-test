package adapters

import (
	"komgrip-test/usecases"

	"gorm.io/gorm"
)

type uowDB struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewUnitOfWorkDB(db *gorm.DB) usecases.UnitOfWork {
	return &uowDB{db: db}
}

func (u *uowDB) checkTX() *gorm.DB {
	tx := u.db
	if u.tx != nil {
		tx = u.tx
	}

	return tx
}

func (u *uowDB) Begin() error {
	u.tx = u.db.Begin()
	if u.tx.Error != nil {
		return u.tx.Error
	}

	return nil
}

func (u *uowDB) Commit() error {
	if u.tx == nil {
		return nil
	}
	if err := u.tx.Commit().Error; err != nil {
		u.Rollback()
		return err
	}

	u.tx = nil
	return nil
}

func (u *uowDB) Rollback() error {
	if u.tx == nil {
		return nil
	}

	err := u.tx.Rollback().Error
	u.tx = nil
	return err
}

func (u *uowDB) BeersRepo() usecases.BeersRepository {
	tx := u.checkTX()
	return NewBeersRepositoryDB(tx)
}
