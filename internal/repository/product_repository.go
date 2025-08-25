package repository

import (
	"database/sql"
	"product-catalog-service/internal/entity"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository (db *sql.DB)*ProductRepository{
	return &ProductRepository{
		db: db,
	}
}


func(p *ProductRepository)CreateProduct(product *entity.Product)(*entity.Product, error){
	query := `INSERT INTO products(name, description, price, stock)VALUES(?,?,?,?)`
	res, err := p.db.Exec(query, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	product.ID = int64(id)
	return  product, nil
}
func(p *ProductRepository)GetProductByID(id int)(*entity.Product, error){
	product := &entity.Product{}
	query := `SELECT id, name, description, price, stock FROM products WHERE id = ?`
	err := p.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		return  nil, err
	} 
	return  product, nil
}

func(p *ProductRepository)UpdateProduct(product *entity.Product)(*entity.Product, error){
	query := `UPDATE products SET nama = ?, descripton = ?, price = ?, stock = ? WHERE id = ?`
	_, err := p.db.Exec(query, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func(p *ProductRepository)DeleteProduct(id int)error{
	query := `DELETE FROM products WHERE id = ?`
	_, err := p.db.Exec(query,id)
	if err != nil {
		return  err
	}
	return nil
}

func(p *ProductRepository)GetProducts()([]*entity.Product, error){
	var products []*entity.Product
	query := `SELECT id, name, description, price, stock FROM products`
	rows, err := p.db.Query(query)
	if err != nil {
		return  nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}