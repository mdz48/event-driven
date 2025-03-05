package application

import "event-driven/src/features/orders/domain"

type CreateOrderUseCase struct {
	db domain.IOrder
}

func NewCreateOrderUseCase(db domain.IOrder) *CreateOrderUseCase {
	return &CreateOrderUseCase{db: db}
}

func (uc *CreateOrderUseCase) Execute(order domain.Order) (domain.Order, error) {
	return uc.db.CreateOrder(order)
}