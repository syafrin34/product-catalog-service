package main

import (
	"product-catalog-service/internal/api"
	"product-catalog-service/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {

	// initialize product service

	productService := service.NewProductService()
	productHandler := api.NewProductHandler(*productService)

	// initialize echo
	e := echo.New()

	// routes
	e.GET("/products/:id/stock", productHandler.GetProductStock)
	e.GET("/products/:id/reserve", productHandler.ReserveProductStock)
	e.GET("/products/:id/release", productHandler.ReleaseProductStock)

	// start server
	e.Logger.Fatal(e.Start(":8081"))

	

}
