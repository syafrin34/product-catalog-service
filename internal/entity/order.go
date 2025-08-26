package entity

type Order struct {
	ID              int              `json:"id"`
	UserID          int              `json:"user_id"`
	OrderID         int              `json:"order_id"`
	Productrequests []ProductRequest `json:"product_requests"`
	Quantity        int              `json:"quantity"`
	Total           int              `json:"total"`
	TotalMarkUp     int              `json:"total_mark_up"`
	TotalDiscount   int              `json:"total_discount"`
	Status          int              `json:"id"`
}

type ProductRequest struct {
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	MarkUp     float64 `json:"mark_up"`
	Discount   float64 `json:"discount"`
	FinalPrice float64 `json:"final_price"`
}
