package user

import (
	"eventapp/entities"
)

type User interface {
	Get() ([]entities.User, error)
	GetById(id int) (entities.User, error)
	Create(entities.User) (entities.User, int, error)
	Update(id int, user entities.User) (entities.User, error)
	Delete(id int) error
}
