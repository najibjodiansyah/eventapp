package user

import (
	"eventapp/entities/graph/model"
)

type User interface {
	Get() ([]model.User, error)
	GetbyId(id int) (model.User, error)
	Create(model.User) (model.User, error)
	Update(id int, user model.User) (model.User, error)
	Delete(id int) error
}