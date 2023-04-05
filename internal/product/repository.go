package product

import (
	"errors"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"
)

type Repository interface {
	Create(*domain.Product) error
	GetAll() ([]domain.Product, error)
	GetConsumerPrice([]int) (float64, []domain.Product, error)
	GetByID(id int) (*domain.Product, error)
	GetBySearch(priceGt float64) ([]domain.Product, error)
	Update(*domain.Product) error
	Delete(id int) error
}

var ErrProductCodeValueExists = errors.New("The codevalue already exists")
var ErrProductNotFound = errors.New("The requested product was not found")
