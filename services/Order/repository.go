package Order

import (
	"OrderServer/dbutil"
	"OrderServer/repobase"
	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

type Repository struct {
	repobase.Base
}

func NewRepository(sqlDB *sqlx.DB) *Repository {
	return &Repository{
		repobase.NewBase(sqlDB),
	}
}

func (r *Repository) txx(fn func(txRepo *Repository) error) error {
	return dbutil.WithTransaction(r.DB(), func(tx *sqlx.Tx) error {
		repoCopy := *r
		repoCopy.Base = r.Base.CopyWithTX(tx)
		return fn(&repoCopy)
	})
}
func (r *Repository) AddOrder(order AddOrderModel) (error, string) {
	psql := sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)
	insertBuilder := psql.Insert("orders").Columns("product_id", "user_id", "address_id", "cost", "payment_by_card", "payment_by_cod", "date_of_delivery")

	for i, cart := range order.CartDetails.Products {
		insertBuilder.Values(cart.ProductID, order.UserID, order.AddressID, (1.2*order.CartDetails.Products[i].Price)+ShippingCharges, order.PaymentByCard, order.PaymentByCod, order.DateOfDelivery)
	}
	sql, args, err := insertBuilder.ToSql()
	sql = sql + ` Returning id`
	if err != nil {
		log.Printf("AddCartProductsHelper Error: %v", err)
		return err, ""
	}
	var orderID string
	err = r.Get(&orderID, sql, args...)
	if err != nil {
		log.Printf("AddCartProductsHelper Error: %v", err)
		return err, ""
	}
	return nil, orderID
}

func (r *Repository) UpdateProductQuantity(cart CartModel) error {
	//language=SQL
	SQL := `SELECT quantity 
			from products p 
			where p.id in (?) and
          	archived_at IS NULL`
	quantity := make([]int, 0)
	productIDs := make([]string, 0)
	for _, product := range cart.Products {
		productIDs = append(productIDs, product.ProductID)
	}
	sqlStatement, args, err := sqlx.In(SQL, productIDs)
	sqlStatement = r.Rebind(sqlStatement)
	err = r.Sel(&quantity, sqlStatement, args...)
	if err != nil {
		log.Printf("UpdateProductQuantity: %v", err)
		return err
	}
	psql := sqrl.StatementBuilder
	insertBuilder := psql.Update("products")
	insertCase := sqrl.Case("id")
	for i, productID := range productIDs {
		insertCase.When(`'`+productID+`'`, strconv.Itoa(quantity[i]-cart.Products[i].Quantity))
	}
	insertBuilder.Set("quantity", insertCase)
	insertBuilder.Where("id IN (?)")
	sql, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("UpdateProductQuantity Error: %v", err)
		return err
	}
	sql, args, err = sqlx.In(sql, productIDs)
	sql = r.Rebind(sql)
	_, err = r.Exec(sql, args...)
	if err != nil {
		log.Printf("UpdateProductQuantity Error: %v", err)
		return err
	}
	return nil
}
func (r *Repository) RemoveCartHelper(cartID string) error {
	//language=SQL
	SQL := `UPDATE carts
		   SET archived_at=now()
		   WHERE user_id=$1`
	_, err := r.Exec(SQL, cartID)
	if err != nil {
		log.Printf("RemoveCartHelper Error: %v", err)
		return err
	}
	return nil
}

func (r *Repository) FetchCartDetailsHelper(userID string) (CartModel, error) {
	//language=SQL
	var cartDetails CartModel
	cartProducts := make([]AddCartProductModel, 0)
	go func() {
		SQL := `SELECT total_cost 
		   from carts 
		   where user_id=$1 and
		   archived_at IS NULL`
		var cost float64
		err := r.Get(&cost, SQL, userID)
		if err != nil {
			log.Printf("FetchCartDetailsHelper Error: %v ", err)
		}
		cartDetails.TotalCostWithShipping = cost
	}()
	SQL := `SELECT cp.product_id,cp.quantity,p.price
			FROM cart_products cp 
			JOIN products p
			on cp.product_id=p.id
			JOIN carts c on cp.cart_id = c.id
			WHERE c.user_id=$1 AND      
            c.archived_at IS NULL`

	err := r.Sel(&cartProducts, SQL, userID)
	if err != nil {
		log.Printf("FetchCartDetailsHelper Error: %v ", err)
		return cartDetails, err
	}
	if err != nil {
		return cartDetails, err
	}
	cartDetails.Products = cartProducts
	return cartDetails, nil
}
func (r *Repository) FetchOrderDetails(orderID string) (error, OrderModel) {
	SQL := `SELECT id,product_id,user_id,address_id,cost,payment_by_card,payment_by_cod,date_of_delivery 
		  from orders 
		  where id=$1 AND
          archived_at IS NULL`
	var orderDetails OrderModel
	err := r.Get(&orderDetails, SQL, orderID)
	if err != nil {
		log.Printf("FetchOrderDetails Error: %v ", err)
		return err, orderDetails
	}
	return nil, orderDetails
}
