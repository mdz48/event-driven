package application

import "citasAPI/src/features/users/domain"

type LoginUserUseCase struct {
	userRepository domain.IUser
}

func NewLoginUserUseCase(db domain.IUser) *LoginUserUseCase {
	return &LoginUserUseCase{userRepository: db}
}

func (l *LoginUserUseCase) Execute(email string, password string) (domain.User, error) {
	return l.userRepository.Login(email, password)
}