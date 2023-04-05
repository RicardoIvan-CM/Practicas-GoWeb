package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/product"
)

type JSONRepository struct {
	file *os.File
	data []domain.Product
}

var (
	ErrFileNotOpened   = errors.New("The file could not be opened")
	ErrFileNotReadable = errors.New("The file could not be read")
	ErrFileNotWritable = errors.New("The file could not be written")
)

func NewJSONRepository(fileName string) (repository *JSONRepository, err error) {
	var store = &JSONRepository{}

	archivo, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, ErrFileNotOpened
	}

	store.file = archivo

	bytes, err := io.ReadAll(archivo)
	if err != nil {
		return nil, ErrFileNotReadable
	}

	json.Unmarshal(bytes, &store.data)
	return store, nil
}

func (repository *JSONRepository) writeJSON() error {
	//Escribir en el archivo
	bytes, _ := json.MarshalIndent(repository.data, "", "")
	repository.file.Truncate(0)
	_, err := repository.file.WriteAt(bytes, 0)
	if err != nil {
		return err
	}

	return nil
}

func (repository *JSONRepository) Create(producto *domain.Product) (*domain.Product, error) {
	producto.ID = len(repository.data) + 1
	for _, p := range repository.data {
		if producto.CodeValue == p.CodeValue {
			return nil, product.ErrProductCodeValueExists
		}
	}
	repository.data = append(repository.data, *producto)
	return producto, repository.writeJSON()
}

func (repository *JSONRepository) GetAll() (result []domain.Product, err error) {
	result = repository.data
	return result, nil
}

type boughtProduct struct {
	Product      *domain.Product
	SoldQuantity int
}

func (repository *JSONRepository) GetConsumerPrice(ids []int) (price float64, products []domain.Product, err error) {

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

func (repository *JSONRepository) GetByID(id int) (result *domain.Product, err error) {
	for _, product := range repository.data {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, product.ErrProductNotFound
}

func (repository *JSONRepository) GetBySearch(priceGt float64) (result []domain.Product, err error) {
	foundProducts := []domain.Product{}
	for _, product := range repository.data {
		if product.Price > priceGt {
			foundProducts = append(foundProducts, product)
		}
	}
	return foundProducts, nil
}

func (repository *JSONRepository) Update(producto *domain.Product) error {
	for i, p := range repository.data {
		if p.ID == producto.ID {
			repository.data[i] = *producto
			return repository.writeJSON()
		}
	}
	return product.ErrProductNotFound
}

func (repository *JSONRepository) Delete(id int) error {
	for index, product := range repository.data {
		if product.ID == id {
			repository.data = append(repository.data[:index], repository.data[index+1:]...)
			return repository.writeJSON()
		}
	}
	return product.ErrProductNotFound
}
