package auth

import (
	"database/sql"
	"eventapp/delivery/middlewares"
	"eventapp/entities/graph/model"
	"fmt"
)

type AuthRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *AuthRepository{
	return &AuthRepository{db:db}
}

func (r *AuthRepository) Login(email string, password string)(model.User,  string,  error){
	var err error
	var authToken string
	var user model.User

	res, err := r.db.Query("select id,name,email,password from users where email = ? and password = ?",email,password)
	if err != nil {
		return  user,authToken,fmt.Errorf("query sql salah")
	}
	fmt.Println(res)
	// defer res.Close()
	for res.Next(){
		err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return user," ", err
		}
	}

	//bandingin hash password

	if user.Email != email && user.Password != password {
		return user,"",fmt.Errorf("user tidak ditemukan")
	}

	authToken, err = middlewares.CreateToken((int(*user.ID)))
		if err != nil {
			return  user,"",fmt.Errorf("create token gagal")
		}
		return user, authToken, nil
}