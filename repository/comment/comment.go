package comment

import (
	"database/sql"
	"errors"
	"eventapp/entities"
	"log"
)

type CommentRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) GetCommentsByEventId(eventId int) ([]entities.Comment, error) {
	stmt, err := r.db.Prepare("select u.id, u.name, u.avatar, c.comment, c.created_at from users u left join comments c on c.userid = u.id where c.deleted_at is NULL and c.eventid = ?")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result, err := stmt.Query(eventId)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer result.Close()

	var comments []entities.Comment

	for result.Next() {
		var comment entities.Comment

		err := result.Scan(&comment.UserId, &comment.UserName, &comment.Avatar, &comment.Content, &comment.CreatedAt) // ada baiknya dibuat istilah yang konsisten, photo atau avatar atau image

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) CreateComment(eventId int, userId int, comment string) error {
	stmt, err := r.db.Prepare("insert into comments(userid, eventid, comment, created_at) values(?, ?, ?, CURRENT_TIMESTAMP)")

	if err != nil {
		log.Fatal(err)
		return errors.New("failed create comment")
	}

	result, err := stmt.Exec(eventId, userId, comment)

	if err != nil {
		log.Fatal(err)
		return errors.New("failed create comment")
	}

	row, _ := result.RowsAffected()

	if row == 0 {
		return errors.New("failed create comment")
	}

	return nil
}

func (r *CommentRepository) DeleteComment(eventId int, userId int) error {
	stmt, err := r.db.Prepare("update comments set deleted_at = CURRENT_TIMESTAMP where eventid = ? and userid = ?")

	if err != nil {
		log.Fatal(err)
		return errors.New("failed delete comment")
	}

	result, err := stmt.Exec(eventId, userId)

	if err != nil {
		log.Fatal(err)
		return err
	}

	row, _ := result.RowsAffected()

	if row == 0 {
		return errors.New("failed delete comment")
	}

	return nil
}
