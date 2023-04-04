package product

import "github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"

type SliceRepository struct {
	data          []domain.Product
	lastProductID int
}

func (repository *SliceRepository) Create(product *domain.Product) error {
	repository.lastProductID++
	product.ID = repository.lastProductID
	for _, p := range repository.data {
		if product.CodeValue == p.CodeValue {
			return ErrProductCodeValueExists
		}
	}
	repository.data = append(repository.data, *product)
	return nil
}

func (repository *SliceRepository) GetAll() (result []domain.Product, err error) {
	result = repository.data
	return result, nil
}

func (repository *SliceRepository) GetByID(id int) (result *domain.Product, err error) {
	for _, product := range repository.data {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, ErrProductNotFound
}

func (repository *SliceRepository) GetBySearch(priceGt float64) (result []domain.Product, err error) {
	foundProducts := []domain.Product{}
	for _, product := range repository.data {
		if product.Price > priceGt {
			foundProducts = append(foundProducts, product)
		}
	}
	return foundProducts, nil
}
