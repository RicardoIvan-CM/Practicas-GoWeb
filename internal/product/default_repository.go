package product

import (
	"errors"
	"fmt"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/pkg/store"
)

type DefaultRepository struct {
	storage store.Store
}

func NewDefaultRepository(storage store.Store) Repository {
	return &DefaultRepository{
		storage,
	}
}

func (repository *DefaultRepository) Create(product domain.Product) (domain.Product, error) {
	list, err := repository.storage.GetAll()
	if err != nil {
		return domain.Product{}, err
	}
	for _, p := range list {
		if product.CodeValue == p.CodeValue {
			return domain.Product{}, ErrProductCodeValueExists
		}
	}
	repository.storage.AddOne(product)
	return product, nil
}

func (repository *DefaultRepository) GetAll() (result []domain.Product, err error) {
	list, err := repository.storage.GetAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (repository *DefaultRepository) GetByID(id int) (result domain.Product, err error) {
	list, err := repository.storage.GetAll()
	if err != nil {
		return domain.Product{}, err
	}
	for _, product := range list {
		if product.ID == id {
			return product, nil
		}
	}
	return domain.Product{}, ErrProductNotFound
}

func (repository *DefaultRepository) GetBySearch(priceGt float64) (result []domain.Product, err error) {
	list, err := repository.storage.GetAll()
	if err != nil {
		return nil, err
	}
	foundProducts := []domain.Product{}
	for _, product := range list {
		if product.Price > priceGt {
			foundProducts = append(foundProducts, product)
		}
	}
	return foundProducts, nil
}

type boughtProduct struct {
	Product      *domain.Product
	SoldQuantity int
}

func (repository *DefaultRepository) GetConsumerPrice(ids []int) (price float64, products []domain.Product, err error) {
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
				&product,
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

func (repository *DefaultRepository) Update(product domain.Product) (domain.Product, error) {
	list, err := repository.storage.GetAll()
	if err != nil {
		return domain.Product{}, err
	}
	for _, p := range list {
		if product.CodeValue == p.CodeValue {
			return domain.Product{}, ErrProductCodeValueExists
		}
	}
	err = repository.storage.UpdateOne(product)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (repository *DefaultRepository) Delete(id int) error {
	err := repository.storage.DeleteOne(id)
	if err != nil {
		return err
	}
	return nil
}
