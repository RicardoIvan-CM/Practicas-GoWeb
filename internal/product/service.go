package product

import "github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"

type Service interface {
	Create(*domain.Product) error
	GetAll() ([]domain.Product, error)
	GetByID(id int) (*domain.Product, error)
	GetBySearch(priceGt float64) ([]domain.Product, error)
	Update(*domain.Product) error
	Delete(id int) error
}
