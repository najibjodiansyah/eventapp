package user

import (
	"eventapp/entities"
)

type User interface {
	Get() ([]entities.User, error)
	GetUserById(id int) (entities.User, int, error)
	CreateUser(entities.User) (entities.User, int, error)
	UpdateUser(entities.User) (entities.User, int, error)
	DeleteUser(id int) (int, error)
}
