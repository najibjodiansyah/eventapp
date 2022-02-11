package user

import (
	_entities "eventapp/entities"
)

type User interface {
	GetAllUsers() ([]_entities.User, error)
	GetUserById(id int) (_entities.User, error)
	CreateUser(_entities.User) (_entities.User, int, error)
	UpdateUser(_entities.User) (_entities.User, int, error)
	DeleteUser(id int) (int, error)
}
