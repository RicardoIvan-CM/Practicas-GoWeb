package store

import (
	"encoding/json"
	"errors"
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

	archivo, err := os.Open(fileName)
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
	_, err := repository.file.Write(bytes)
	if err != nil {
		return ErrFileNotWritable
	}

	return nil
}

func (repository *JSONRepository) Create(producto *domain.Product) error {
	producto.ID = len(repository.data) + 1
	for _, p := range repository.data {
		if producto.CodeValue == p.CodeValue {
			return product.ErrProductCodeValueExists
		}
	}
	repository.data = append(repository.data, *producto)
	return repository.writeJSON()
}

func (repository *JSONRepository) GetAll() (result []domain.Product, err error) {
	result = repository.data
	return result, nil
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
