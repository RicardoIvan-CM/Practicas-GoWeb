package product

import (
	"errors"
	"fmt"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"
)

type SliceRepository struct {
	data          []domain.Product
	lastProductID int
}

func NewSliceRepository() (repository *SliceRepository) {
	return &SliceRepository{
		data: []domain.Product{},
	}
}

func (repository *SliceRepository) Create(product *domain.Product) (*domain.Product, error) {
	repository.lastProductID++
	product.ID = repository.lastProductID
	for _, p := range repository.data {
		if product.CodeValue == p.CodeValue {
			return nil, ErrProductCodeValueExists
		}
	}
	repository.data = append(repository.data, *product)
	return product, nil
}

func (repository *SliceRepository) GetAll() (result []domain.Product, err error) {
	result = repository.data
	return result, nil
}

type boughtProduct struct {
	Product      *domain.Product
	SoldQuantity int
}

func (repository *SliceRepository) GetConsumerPrice(ids []int) (price float64, products []domain.Product, err error) {

	boughtProducts := make(map[int]*boughtProduct)

	var cuentaProductos int
	var precioTotal float64 = 0.0

	for _, id := range ids {
		product, err := repository.GetByID(id)
		if err != nil {
			return 0, nil, err
		}
		if !product.IsPublished {
			return 0, nil, errors.New(fmt.Sprintf("El producto %d no está publicado", id))
		}
		if bp, ok := boughtProducts[id]; ok {
			bp.SoldQuantity++
			if bp.SoldQuantity > bp.Product.Quantity {
				return 0, nil, errors.New(fmt.Sprintf("La cantidad solicitada del producto %d supera el stock", id))
			}
		} else {
			boughtProducts[id] = &boughtProduct{
				product,
				1,
			}
			cuentaProductos++
		}
		precioTotal += product.Price
	}

	boughtList := []domain.Product{}

	for _, bp := range boughtProducts {
		boughtList = append(boughtList, *bp.Product)
	}

	if cuentaProductos > 20 {
		return precioTotal * 1.15, boughtList, nil
	}

	if cuentaProductos > 10 {
		return precioTotal * 1.17, boughtList, nil
	}

	return precioTotal * 1.21, boughtList, nil
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
	for index, product := range repository.data {
		if product.ID == id {
			repository.data = append(repository.data[:index], repository.data[index+1:]...)
			return nil
		}
	}
	return ErrProductNotFound
}
