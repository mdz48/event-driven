package domain

type IOrder interface {
	GetOrder(orderID int) (Order, error)
	CreateOrder(order Order) (Order, error)
	UpdateStatus(orderID int, status string) (Order, error)
	GetAll() ([]Order, error)
	DeleteOrder(orderID int) error
}