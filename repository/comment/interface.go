package comment

import (
	"eventapp/entities"
)

type Comment interface {
	GetCommentsByEventId(eventId int) ([]entities.Comment, error)
	CreateComment(eventId int, userId int, comment string) error
	DeleteComment(eventId int, userId int) error
}