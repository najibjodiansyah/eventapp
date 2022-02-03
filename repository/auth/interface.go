package auth

import "eventapp/entities/graph/model"

type Auth interface {
	Login(email string, password string)(model.User,string,error)
}