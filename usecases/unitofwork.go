package usecases

type UnitOfWork interface {
	Begin() error
	Commit() error
	Rollback() error
	BeersRepo() BeersRepository
}
