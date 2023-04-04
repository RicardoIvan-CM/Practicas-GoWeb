package product

import "github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"

type DefaultService struct {
	Storage Repository
}

func NewDefaultService(storage Repository) (defaultService *DefaultService) {
	return &DefaultService{
		Storage: storage,
	}
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

func (s DefaultService) Update(product *domain.Product) error {
	return s.Storage.Update(product)
}

func (s DefaultService) Delete(id int) error {
	return s.Storage.Delete(id)
}
