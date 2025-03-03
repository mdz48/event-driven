package domain

type IUser interface {
	Save(user User) (User, error)
	GetAll() ([]User, error)
	GetByID(id int32) (User, error)
	Delete(id int32) error
	Update(id int32, user User) (User, error)
	Login(email string, password string) (User, error)
}