package usecases

import (
	"komgrip-test/entities"
	"os"
	"time"
)

type GetBeersResponse struct {
	ID           uint64 `json:"id"`
	BeerTypeName string `json:"beer_type_name"`
	BeerName     string `json:"beer_name"`
	BeerDesc     string `json:"beer_desc"`
	BeerImgPath  string `json:"beer_img_path"`
}

type GetBeersRequest struct {
	BeerName string
	Page     int
	PageSize int
}

type BeersRequest struct {
	BeerTypeName string `form:"beer_type_name" validate:"required"`
	BeerName     string `form:"beer_name" validate:"required"`
	BeerDesc     string `form:"beer_desc"`
	BeerImgPath  string
}

type BeersUseCase interface {
	CreateBeer(request BeersRequest) error
	GetBeers(request GetBeersRequest) (responses []GetBeersResponse, err error)
	UpdateBeer(id uint64, request BeersRequest) error
	DeleteBeer(id uint64) error
}

type BeersService struct {
	uow UnitOfWork
}

func NewBeersService(uow UnitOfWork) BeersUseCase {
	return &BeersService{uow: uow}
}

func (s *BeersService) setBeerData(request BeersRequest, flag string) entities.Beers {
	data := entities.Beers{
		BeerTypeName: request.BeerTypeName,
		BeerName:     request.BeerName,
		BeerDesc:     request.BeerDesc,
		BeerImgPath:  request.BeerImgPath,
	}
	if flag == "add" {
		data.CreatedAt = time.Now()
		data.UpdatedAt = time.Now()
	} else {
		data.UpdatedAt = time.Now()
	}
	return data
}

func (s *BeersService) CreateBeer(request BeersRequest) error {
	err := s.uow.Begin()
	if err != nil {
		return err
	}
	err = s.uow.BeersRepo().Create(s.setBeerData(request, "add"))
	if err != nil {
		s.uow.Rollback()
		return err
	}
	s.uow.Commit()
	return nil
}

func (s *BeersService) GetBeers(request GetBeersRequest) (responses []GetBeersResponse, err error) {
	page := 1
	if request.Page > 1 {
		page = request.Page
	}
	offset := (page - 1) * request.PageSize
	beers, err := s.uow.BeersRepo().GetDatas(GetDatasParams{
		BeerName: request.BeerName,
		Limit:    request.PageSize,
		Offset:   offset,
	})
	responses = make([]GetBeersResponse, 0, len(beers))
	for _, beer := range beers {
		responses = append(responses, GetBeersResponse{
			ID:           beer.ID,
			BeerTypeName: beer.BeerTypeName,
			BeerName:     beer.BeerName,
			BeerDesc:     beer.BeerDesc,
			BeerImgPath:  beer.BeerImgPath,
		})
	}
	return responses, nil
}

func (s *BeersService) UpdateBeer(id uint64, request BeersRequest) error {
	err := s.uow.Begin()
	if err != nil {
		return err
	}
	err = s.uow.BeersRepo().Update(id, s.setBeerData(request, "update"))
	if err != nil {
		s.uow.Rollback()
		return err
	}
	s.uow.Commit()
	return nil
}

func (s *BeersService) DeleteBeer(id uint64) error {
	err := s.uow.Begin()
	if err != nil {
		return err
	}
	beer, err := s.uow.BeersRepo().GetData(id)
	if err != nil {
		return err
	}
	err = s.uow.BeersRepo().Delete(id)
	if err != nil {
		s.uow.Rollback()
		return err
	}
	s.uow.Commit()
	os.Remove(beer.BeerImgPath)
	return nil
}
