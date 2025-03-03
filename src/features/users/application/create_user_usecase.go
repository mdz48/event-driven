package application

import "citasAPI/src/features/users/domain"

type CreateUserUseCase struct {
	userRepository domain.IUser
}

func NewCreateUserUseCase(userRepository domain.IUser) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: userRepository}
}

func (c *CreateUserUseCase) Execute(user domain.User) (domain.User, error) {
	return c.userRepository.Save(user)
}