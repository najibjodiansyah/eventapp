package auth

import (
	"database/sql"
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
	stmt, err := r.db.Prepare(`select id, name, email, password from users where email = ?`)

	if err != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	res, err := stmt.Query(email)

	if err != nil {
		return entities.User{}, err
	}

	defer res.Close()

	var user entities.User

	if res.Next() {
		err := res.Scan(&user.Id, &user.Name, &user.Email, &user.Password)

		if err != nil {
			return entities.User{}, err
		}
	}

	return user, nil
}
