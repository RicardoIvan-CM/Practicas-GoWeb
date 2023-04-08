package store

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"
)

type Store interface {
	GetAll() ([]domain.Product, error)
	GetOne(id int) (domain.Product, error)
	AddOne(product domain.Product) error
	UpdateOne(product domain.Product) error
	DeleteOne(id int) error
	saveProducts(products []domain.Product) error
	loadProducts() ([]domain.Product, error)
}

type JSONstore struct {
	filePath string
}

var (
	ErrFileNotOpened   = errors.New("The file could not be opened")
	ErrFileNotReadable = errors.New("The file could not be read")
	ErrFileNotWritable = errors.New("The file could not be written")
)

func NewJSONstore(path string) Store {
	return &JSONstore{
		filePath: path,
	}
}

func (store *JSONstore) loadProducts() (data []domain.Product, err error) {
	var products []domain.Product
	file, err := os.ReadFile(store.filePath)
	if err != nil {
		return nil, ErrFileNotOpened
	}
	err = json.Unmarshal(file, &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// saveProducts guarda los productos en un archivo json
func (store *JSONstore) saveProducts(products []domain.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return os.WriteFile(store.filePath, bytes, 0644)
}

// GetAll
func (store *JSONstore) GetAll() (result []domain.Product, err error) {
	products, err := store.loadProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetOne implements Store
func (store *JSONstore) GetOne(id int) (domain.Product, error) {
	products, err := store.loadProducts()
	if err != nil {
		return domain.Product{}, err
	}
	for _, product := range products {
		if product.ID == id {
			return product, nil
		}
	}
	return domain.Product{}, errors.New("The product was not found")
}

// AddOne implements Store
func (store *JSONstore) AddOne(product domain.Product) error {
	products, err := store.loadProducts()
	if err != nil {
		return err
	}
	product.ID = len(products) + 1
	products = append(products, product)
	return store.saveProducts(products)
}

// UpdateOne implements Store
func (store *JSONstore) UpdateOne(product domain.Product) error {
	products, err := store.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.ID == product.ID {
			products[i] = product
			return store.saveProducts(products)
		}
	}
	return errors.New("The product was not found")
}

// DeleteOne implements Store
func (store *JSONstore) DeleteOne(id int) error {
	products, err := store.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return store.saveProducts(products)
		}
	}
	return errors.New("The product was not found")
}
