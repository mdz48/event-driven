package application

import "event-driven/src/features/users/domain"

type GetAllUsersUseCase struct {
	db domain.IUser
}

func NewGetAllUsersUseCase(db domain.IUser) *GetAllUsersUseCase {
	return &GetAllUsersUseCase{db: db}
}

func (g *GetAllUsersUseCase) Execute() ([]domain.User, error) {
	return g.db.GetAll()
}