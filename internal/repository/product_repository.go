package repository

import (
	"product-catalog-service/internal/entity"
)

type ProductRepository struct {

}

func NewProductRepository ()*ProductRepository{
	return &ProductRepository{}
}

var products = map[int]*entity.Product{
	1:{ID: 1, Name: "Product1", Description: "product 1 Description",Price: 100, Stock: 10},
	2:{ID: 2, Name: "Product2", Description: "product 2 Description",Price: 200, Stock: 20},
}
func(p *ProductRepository)CreateProduct(product *entity.Product)(*entity.Product, error){
	product.ID = 3
	products[int(product.ID)] = product
	return product, nil
}
func(p *ProductRepository)GetProductByID(id int)(*entity.Product, error){
	product, ok := products[id]
	if !ok {
		return  nil, nil
	}
	return  product, nil
}

func(p *ProductRepository)UpdateProduct(product *entity.Product)(*entity.Product, error){
	products[int(product.ID)] = product
	return product, nil
}

func(p *ProductRepository)DeleteProduct(id int)error{
	delete(products, id)
	return nil
}

func(p *ProductRepository)GetProducts()([]*entity.Product, error){
	result := make([]*entity.Product, 0, len(products))
	for _, product := range products{
		result = append(result, product)
	}
	return  result, nil
}