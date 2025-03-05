package application

import "event-driven/src/features/orders/domain"

type GetAllUseCase struct {
	db domain.IOrder
}

func NewGetAllUseCase(db domain.IOrder) *GetAllUseCase {
	return &GetAllUseCase{db: db}
}

func (u *GetAllUseCase) Execute() ([]domain.Order, error) {
	orders, err := u.db.GetAll()
	if err != nil {
		return nil, err
	}
	return orders, nil
}