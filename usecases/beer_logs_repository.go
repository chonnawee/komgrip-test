package usecases

import "komgrip-test/entities"

type BeerLogsRepository interface {
	CreateLog(entities.BeerLogs) error
}
