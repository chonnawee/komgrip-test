package repositories

import "komgrip-test/entities"

type BeersRepositoryInterface interface {
	Create(data entities.Beers) error
}
