package comment

import (
	"eventapp/entities"
)

type Comment interface {
	GetCommentsByEventId(eventId int) ([]entities.Comment, error)
	CreateComment(eventId int, userId int, content string) (entities.Comment, error)
	DeleteComment(commentId int, userId int) error
	DeleteAllCommentByUserId(userId int) error
	DeleteAllCommentByEventId(eventId int) error
}
