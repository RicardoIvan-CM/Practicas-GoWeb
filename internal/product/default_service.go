package product

import "github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"

type DefaultService struct {
	Storage Repository
}

func (s DefaultService) Create(product *domain.Product) (err error) {
	return s.Storage.Create(product)
}

func (s DefaultService) GetAll() ([]domain.Product, error) {
	return s.Storage.GetAll()
}
func (s DefaultService) GetByID(id int) (*domain.Product, error) {
	return s.Storage.GetByID(id)
}
func (s DefaultService) GetBySearch(priceGt float64) ([]domain.Product, error) {
	return s.Storage.GetBySearch(priceGt)
}
