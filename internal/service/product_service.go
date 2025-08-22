package service

import (
	"fmt"
	"os"
	"product-catalog-service/internal/repository"

	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(pRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: pRepo,
	}
}

func (p *ProductService) GetProductStock(productID int) (int, error) {
	product, err := p.productRepo.GetProductByID(productID)
	if err != nil {
		logger.Error().Err(err).Msgf("Error getting product by ID", productID)
		return  0, err
	}
	return  product.Stock, nil
}

// reserve stock for an order
func (p *ProductService) ReserveProductStock(productID int, quantity int) error {
	
	product, err := p.productRepo.GetProductByID(productID)
	if err != nil {
		logger.Error().Err(err).Msgf("Error getting product by ID", productID)
		return  err
	}
	if product.Stock < quantity {
		logger.Warn().Msgf("Product %d out of stock", productID)
		return fmt.Errorf("product out of stock")
	}

	product.Stock -= quantity
	_, err = p.productRepo.UpdateProduct(product)
	if err != nil {
		logger.Error().Err(err).Msgf("Error updating product %d", productID)
		return  err
	}
	return  nil
}

// release product reserved stock when an order is cancelled
func (p *ProductService) ReleaseProductStock(productID int, quantity int) error {
	product, err := p.productRepo.GetProductByID(productID)
	if err != nil {
		logger.Error().Err(err).Msgf("Error getting product by ID", productID)
		return  err
	}
	product.Stock += quantity
	_, err = p.productRepo.UpdateProduct(product)
	if err != nil {
		logger.Error().Err(err).Msgf("Error getting product by ID", productID)
		return  err
	}
	return  nil

}
