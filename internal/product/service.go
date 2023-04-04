package product

import "github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"

type Service interface {
	Create(product *domain.Product) error
	GetAll() ([]domain.Product, error)
	GetByID(id int) (*domain.Product, error)
	GetBySearch(priceGt float64) ([]domain.Product, error)
}
