package auth

import (
	"eventapp/entities"
)

type Auth interface {
	Login(email string) (entities.User, error)
}
