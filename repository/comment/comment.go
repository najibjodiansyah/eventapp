package comment

import (
	"database/sql"
	"errors"
	"eventapp/entities"
	"log"
	"time"
)

type CommentRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// ditambah comment id
func (r *CommentRepository) GetCommentsByEventId(eventId int) ([]entities.Comment, error) {
	stmt, err := r.db.Prepare(`select c.id, u.id, u.name, u.avatar, c.comment, c.created_at
								from users u left join comments c on c.userid = u.id
								where c.deleted_at is NULL and c.eventid = ?`)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	res, err := stmt.Query(eventId)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Close()

	var comments []entities.Comment

	for res.Next() {
		var comment entities.Comment

		err := res.Scan(&comment.Id, &comment.UserId, &comment.UserName, &comment.Avatar, &comment.Content, &comment.CreatedAt)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) CreateComment(eventId int, userId int, content string) (entities.Comment, error) {
	stmt, err := r.db.Prepare("insert into comments(userid, eventid, comment, created_at) values(?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
		return entities.Comment{}, err
	}

	created_at := time.Now()

	res, err := stmt.Exec(eventId, userId, content, created_at)

	if err != nil {
		log.Fatal(err)
		return entities.Comment{}, err
	}

	affected, _ := res.RowsAffected()

	if affected == 0 {
		return entities.Comment{}, errors.New("failed create comment")
	}

	id, _ := res.LastInsertId()

	stmt, _ = r.db.Prepare("select name, avatar from users where id = ?")

	row, _ := stmt.Query(userId)

	defer row.Close()

	comment := entities.Comment{
		Id:        int(id),
		UserId:    userId,
		Content:   content,
		CreatedAt: created_at.String(),
	}

	if row.Next() {
		var name, avatar string

		row.Scan(&name, &avatar)

		comment.UserName = name
		comment.Avatar = avatar
	}

	return comment, nil
}

func (r *CommentRepository) DeleteComment(commentId int, userId int) error {
	stmt, err := r.db.Prepare("update comments set deleted_at = CURRENT_TIMESTAMP where id = ? and userid = ?")

	if err != nil {
		log.Fatal(err)
		return errors.New("failed delete comment")
	}

	_, err = stmt.Exec(commentId, userId)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *CommentRepository) DeleteAllCommentByUserId(userId int) error {
	stmt, err := r.db.Prepare("update comments set deleted_at = CURRENT_TIMESTAMP where deleted_at is NULL and userid = ?")

	if err != nil {
		log.Fatal(err)
		return errors.New("failed delete comment")
	}

	_, err = stmt.Exec(userId)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *CommentRepository) DeleteAllCommentByEventId(eventId int) error {
	stmt, err := r.db.Prepare("update comments set deleted_at = CURRENT_TIMESTAMP where deleted_at is NULL and eventid = ?")

	if err != nil {
		log.Fatal(err)
		return errors.New("failed delete comment")
	}

	_, err = stmt.Exec(eventId)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
