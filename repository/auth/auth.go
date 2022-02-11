package auth

import (
	"database/sql"
	"errors"
	"eventapp/entities"
	"log"
)

type AuthRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// return repository berbentuk entity saja
func (r *AuthRepository) Login(email string) (entities.User, error) {
	stmt, err := r.db.Prepare(`select id, name, password from users where email = ?`)

	if err != nil {
		log.Println(err)
		return entities.User{}, errors.New("internal server error")
	}

	res, err := stmt.Query(email)

	if err != nil {
		log.Println(err)
		return entities.User{}, errors.New("internal server error")
	}

	defer res.Close()

	var user entities.User

	if res.Next() {
		err := res.Scan(&user.Id, &user.Name, &user.Password)

		if err != nil {
			log.Println(err)
			return entities.User{}, errors.New("internal server error")
		}
	}

	return user, nil
}
