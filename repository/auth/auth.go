package auth

import (
	"database/sql"
	"eventapp/entities/graph/model"
	"fmt"
)

type AuthRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Login(email string) (model.User, error) {
	var err error
	var user model.User

	res, err := r.db.Query("select id,name,email,password from users where email = ? ", email)
	if err != nil {
		return user, fmt.Errorf("query sql salah")
	}
	// defer res.Close()
	for res.Next() {
		err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			fmt.Println()
			return user, err
		}
	}
	fmt.Println(user)
	return user, nil
}
