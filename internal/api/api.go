package api

import (
	"product-catalog-service/internal/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService)*ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// get product stock
func (p *ProductHandler)GetProductStock(c echo.Context)error{
	productID := c.Param("id")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid product id"})
	}
	ctx := c.Request().Context()
	stock, err := p.productService.GetProductStock(ctx,productIDInt)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]int{"stock": stock})
}


// reserve product stock for a product

func (p *ProductHandler)ReserveProductStock(c echo.Context)error{
	reservation := struct {
		ProductID int `json:"product_id"`
		Quantity int `json:"quantity"`
	}{}

	if err := c.Bind(&reservation); err != nil {
			return c.JSON(400, map[string]string{"error":"Invalid request payload"})
	}

	ctx := c.Request().Context()
	err := p.productService.ReserveProductStock(ctx, reservation.ProductID, reservation.Quantity)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]string{"message":"stock reserved"})
}

// release product stock 
func (p *ProductHandler)ReleaseProductStock (c echo.Context)error{
	release := struct {
		ProductID int `json:"product_id"`
		Quantity int `json:"quantity"`
	}{}

	if err := c.Bind(&release); err != nil {
			return c.JSON(400, map[string]string{"error":"Invalid request payload"})
	}

	ctx := c.Request().Context()
	err := p.productService.ReleaseProductStock(ctx, release.ProductID, release.ProductID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]string{"message":"stock released"})
}