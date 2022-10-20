package Order

import "time"

type AddOrderModel struct {
	UserID         string    `db:"user_id" json:"userID"`
	AddressID      string    `db:"address_id" json:"addressID"`
	CartID         string    `db:"cart_id" json:"cartID"`
	CartDetails    CartModel `db:"-" json:"cartDetails"`
	PaymentByCard  string    `db:"payment_by_card" json:"paymentByCard"`
	PaymentByCod   bool      `db:"payment_by_cod" json:"paymentByCod"`
	DateOfDelivery time.Time `db:"date_of_delivery" json:"dateOfDelivery"`
}
type CartModel struct {
	Products              []AddCartProductModel `db:"products" json:"products"`
	TotalCostWithShipping float64               `db:"total_cost"`
}
type AddCartProductModel struct {
	ProductID string  `db:"product_id" json:"productID"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Price     float64 `db:"price" json:"price"`
}
type OrderModel struct {
	OrderID        string    `db:"id" json:"orderID"`
	ProductID      string    `db:"product_id" json:"productID"`
	UserID         string    `db:"user_id" json:"userID"`
	AddressID      string    `db:"address_id" json:"addressID"`
	Cost           string    `db:"cost" json:"cost"`
	PaymentByCard  string    `db:"payment_by_card" json:"paymentByCard"`
	PaymentByCod   bool      `db:"payment_by_cod" json:"paymentByCod"`
	DateOfDelivery time.Time `db:"date_of_delivery" json:"dateOfDelivery"`
}
