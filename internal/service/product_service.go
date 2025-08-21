package service

import "fmt"

type ProductService struct {
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (p *ProductService) GetProductStock(productID int) (int, error) {
	stock := 50
	return stock, nil
}
func (p *ProductService) ReserveProductStock(productID int, quantity int) error {
	fmt.Printf("Reserved %d units of product %d\n", quantity, productID)
	return nil
}
func (p *ProductService) ReleaseProductStock(productID int, quantity int) error {
	fmt.Printf("Release %d units of product %d\n", quantity, productID)
	return nil

}
