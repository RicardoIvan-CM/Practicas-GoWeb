package product

import "github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"

type DefaultService struct {
	Repository Repository
}

func NewDefaultService(repository Repository) Service {
	return &DefaultService{
		Repository: repository,
	}
}

func (s *DefaultService) Create(product domain.Product) (domain.Product, error) {
	return s.Repository.Create(product)
}

func (s *DefaultService) GetAll() ([]domain.Product, error) {
	return s.Repository.GetAll()
}
func (s *DefaultService) GetConsumerPrice(ids []int) (float64, []domain.Product, error) {
	return s.Repository.GetConsumerPrice(ids)
}
func (s *DefaultService) GetByID(id int) (domain.Product, error) {
	return s.Repository.GetByID(id)
}
func (s *DefaultService) GetBySearch(priceGt float64) ([]domain.Product, error) {
	return s.Repository.GetBySearch(priceGt)
}

func (s *DefaultService) Update(product domain.Product) (domain.Product, error) {
	return s.Repository.Update(product)
}

func (s *DefaultService) Delete(id int) error {
	return s.Repository.Delete(id)
}
