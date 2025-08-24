package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"product-catalog-service/internal/entity"
	"product-catalog-service/internal/repository"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
type ProductService struct {
	productRepo repository.ProductRepository
	rdb *redis.Client
}


func NewProductService(pRepo repository.ProductRepository, rdb *redis.Client) *ProductService {
	return &ProductService{
		productRepo: pRepo,
		rdb: rdb,
		
	}
}

func (p *ProductService) GetProductStock(ctx context.Context,  productID int) (int, error) {
	// get product from redis
	key := fmt.Sprintf("product:%d", productID)
	productCache, err := p.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil){
			logger.Warn().Msgf("stock for product %d not found in cache", productID)
		}else {
			
			logger.Error().Err(err).Msgf("error getting  product %d from cache", productID)
			return 0, err
		}
	}
	if productCache != "" {
		var product entity.Product
		err = json.Unmarshal([]byte(productCache), &product)
		if err != nil {
			logger.Error().Err(err).Msgf("Error un marshalling product %d", product)
			return 0, err
		}
		logger.Info().Msgf("retrieved stock for product %d: %d", productID, product.Stock)
	}



   // get get product to db
	product, err := p.productRepo.GetProductByID(productID)
	if err != nil {
		logger.Error().Err(err).Msgf("Error getting product by ID", productID)
		return  0, err
	}

	// save produtc to redis
	err = p.rdb.Set(ctx, key, product,0).Err()
	return  product.Stock, nil
}



// reserve stock for an order
func (p *ProductService) ReserveProductStock(ctx context.Context, productID int, quantity int) error {
	
	// get product from redis
	key := fmt.Sprintf("product:%d", productID)
	productCache, err := p.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil){
			logger.Warn().Msgf("stock for product %d not found in cache", productID)
		}else {
			
			logger.Error().Err(err).Msgf("error getting  product %d from cache", productID)
			return err
		}
	}

	var productData entity.Product
	err = json.Unmarshal([]byte(productCache), &productData)
	if err != nil {
			logger.Error().Err(err).Msgf("Error un marshalling product %d", productID)
			return  err
	}

	if productData.ID == 0 {
		product, err := p.productRepo.GetProductByID(productID)
		if err != nil {
		logger.Error().Err(err).Msgf("Error getting product by ID", productData.ID)
		return  err
		}
		productData = *product
	}

	if productData.Stock < quantity {
		logger.Warn().Msgf("Product %d out of stock", productID)
		return fmt.Errorf("product out of stock")
	}

	productData.Stock -= quantity
	_, err = p.productRepo.UpdateProduct(&productData)
	if err != nil {
		logger.Error().Err(err).Msgf("Error updating product %d", productID)
		return  err
	}

	// delete from cache redis
	err = p.rdb.Del(ctx, key).Err()
	if err != nil {
		logger.Error().Err(err).Msgf("Error deleting product %d from cache", productData.ID)
		return  err
	}
	return  nil
}

// release product reserved stock when an order is cancelled
func (p *ProductService) ReleaseProductStock(ctx context.Context, productID int, quantity int) error {
	// get product from redis
	key := fmt.Sprintf("product:%d", productID)
	productCache, err := p.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil){
			logger.Warn().Msgf("stock for product %d not found in cache", productID)
		}else {
			
			logger.Error().Err(err).Msgf("error getting  product %d from cache", productID)
			return err
		}
	}

	var productData entity.Product
	err = json.Unmarshal([]byte(productCache), &productData)
	if err != nil {
			logger.Error().Err(err).Msgf("Error un marshalling product %d", productID)
			return  err
	}

	if productData.ID == 0 {
		product, err := p.productRepo.GetProductByID(productID)
		if err != nil {
		logger.Error().Err(err).Msgf("Error getting product by ID", productData.ID)
		return  err
		}
		productData = *product
	}


	productData.Stock += quantity
	_, err = p.productRepo.UpdateProduct(&productData)
	if err != nil {
		logger.Error().Err(err).Msgf("Error getting product by ID", productID)
		return  err
	}

	// Delete from cache redis
	err = p.rdb.Del(ctx, key).Err()
	if err != nil {
		logger.Error().Err(err).Msgf("Error deleting product %d from cache", productData.ID)
		return  err
	}
	return  nil

}
