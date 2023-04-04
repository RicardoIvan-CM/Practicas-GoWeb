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

func (repository *SliceRepository) Update(product *domain.Product) error {
	for i, p := range repository.data {
		if p.ID == product.ID {
			repository.data[i] = *product
			return nil
		}
	}
	return ErrProductNotFound
}

func (repository *SliceRepository) Delete(id int) error {
	var index int = -1
	for i, p := range repository.data {
		if p.ID == id {
			index = i
			break
		}
	}
	if index > 0 {
		repository.data = append(repository.data[:index], repository.data[index+1])
		return nil
	}
	return ErrProductNotFound
}
