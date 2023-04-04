package handler

import (
	"errors"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"
)

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CodeValue   string  `json:"code_value" binding:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

var (
	ErrProductNameRequired       = errors.New("The name is required")
	ErrProductQuantityInvalid    = errors.New("The quantity is not valid")
	ErrProductCodeValueRequired  = errors.New("The codevalue is required")
	ErrProductExpirationRequired = errors.New("The expiration date is required")
	ErrProductExpirationInvalid  = errors.New("The expiration date is not valid")
	ErrProductPriceInvalid       = errors.New("The price is invalid")
)

func (req *ProductRequest) ToDomain() domain.Product {
	var newProduct = domain.Product{
		Name:        req.Name,
		Quantity:    req.Quantity,
		CodeValue:   req.CodeValue,
		IsPublished: req.IsPublished,
		Expiration:  req.Expiration,
		Price:       req.Price,
	}
	return newProduct
}
