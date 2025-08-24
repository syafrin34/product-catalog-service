package main

import (
	"product-catalog-service/internal/api"
	"product-catalog-service/internal/repository"
	"product-catalog-service/internal/service"

	"github.com/go-redis/redis/v8"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// initialize product service
	productRepo := repository.NewProductRepository()
	productService := service.NewProductService(*productRepo, rdb)
	productHandler := api.NewProductHandler(*productService)

	// initialize echo
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echojwt.JWT([]byte("secret")))

	// routes
	e.GET("/products/:id/stock", productHandler.GetProductStock)
	e.GET("/products/:id/reserve", productHandler.ReserveProductStock)
	e.GET("/products/:id/release", productHandler.ReleaseProductStock)

	// start server
	e.Logger.Fatal(e.Start(":8081"))

	

}
