package usecases

import "komgrip-test/entities"

type GetDatasParams struct {
	BeerName string
	Limit    int
	Offset   int
}

type BeersRepository interface {
	Create(beer entities.Beers) error
	GetDatas(params GetDatasParams) (beers []entities.Beers, err error)
	Update(id uint64, beer entities.Beers) error
	Delete(id uint64) error
}
