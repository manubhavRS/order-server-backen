package Order

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}
func (s *Service) AddOrderHelper(order AddOrderModel) (error, string) {
	//	var cart CartModel
	var orderID string
	var err error
	txErr := s.repository.txx(func(txRepo *Repository) error {
		err, orderID = s.repository.AddOrder(order)
		if err != nil {
			return err
		}
		err = s.repository.UpdateProductQuantity(order.CartDetails) //Bulk update using sqrl and sqlx.IN()
		if err != nil {
			return err
		}
		err = s.repository.RemoveCartHelper(order.UserID)
		return err
	})
	if txErr != nil {
		return txErr, ""
	}
	return nil, orderID
}
func (s *Service) FetchOrderHelper(order string) (error, OrderModel) {
	err, orderDetails := s.repository.FetchOrderDetails(order)
	if err != nil {
		return err, OrderModel{}
	}
	return nil, orderDetails
}
func (s *Service) FetchCartDetails(userID string) (CartModel, error) {
	var cart CartModel
	cart, err := s.repository.FetchCartDetailsHelper(userID)
	if err != nil {
		return cart, err
	}
	var num int
	for _, cartProduct := range cart.Products {
		num = num + cartProduct.Quantity
	}
	cart.TotalCostWithShipping = cart.TotalCostWithShipping + float64(num*ShippingCharges)
	return cart, err
}
