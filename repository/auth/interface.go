package auth

import (
	"eventapp/entities/graph/model"
)

type Auth interface {
	Login(email string)(model.User,error)
}