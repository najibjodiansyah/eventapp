package user

import (
	"database/sql"
	"eventapp/entities/graph/model"
	"fmt"
	"log"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Get() ([]model.User, error) {

	stmt, err := r.db.Prepare("select id, name, email, phone, avatar from users where deleted_at is NULL")
	if err != nil {
		log.Fatal(err)
	}

	var users []model.User

	result, err := stmt.Query()
	if err != nil {
		return users, err
	}

	defer result.Close()

	for result.Next() {
		var user model.User
		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Avatar)
		if err != nil {
			log.Fatal("error di scan getUser")
		}
		users = append(users, user)
	}
	return users, nil

}

func (r *UserRepository) GetbyId(id int) (model.User, error) {
	var user model.User
	stmt, err := r.db.Prepare("select id, name, email, phone, avatar from users where id = ? and deleted_at is NULL")
	if err != nil {
		//log.Fatal(err)
		return user, fmt.Errorf("gagal prepare db")
	}

	result, err := stmt.Query(id)
	if err != nil {
		return user, fmt.Errorf("gagal query user")
	}

	defer result.Close()

	for result.Next() {
		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Avatar)
		if err != nil {
			return user, err
		}
		return user, nil
	}

	return user, fmt.Errorf("user not found")
}

func (r *UserRepository) Create(user model.User) (model.User, error) {
	stmt, err := r.db.Prepare("INSERT INTO users(name, email, password, phone, avatar) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Exec(user.Name, user.Email, user.Password, user.PhoneNumber, user.Avatar)
	if err != nil {
		return user, fmt.Errorf("gagal exec")
	}

	notAffected, _ := result.RowsAffected()
	if notAffected == 0 {
		return user, fmt.Errorf("user not created")
	}

	return user, nil
}

func (r *UserRepository) Update(id int, user model.User) (model.User, error) {
	stmt, err := r.db.Prepare("UPDATE users SET name= ?, email= ?, password= ?, phone= ?, avatar= ? WHERE id = ? AND deleted_at is NULL")
	if err != nil {
		// log.Fatal(err)
		return user, fmt.Errorf("gagal prepare update")
	}

	result, error := stmt.Exec(user.Name, user.Email, user.Password, user.PhoneNumber, user.Avatar, id)
	if error != nil {
		return user, fmt.Errorf("gagal exec update")
	}

	notAffected, _ := result.RowsAffected()
	if notAffected == 0 {
		return user, fmt.Errorf("row not affected")
	}

	return user, nil
}

func (r *UserRepository) Delete(id int) error {
	// stmt, err := r.db.Prepare("DELETE from users where id = ?")
	stmt, err := r.db.Prepare("UPDATE users SET deleted_at= ? where id = ?")
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Exec(time.Now(), id)
	if err != nil {
		return fmt.Errorf("salah user")
	}

	notAffected, _ := result.RowsAffected()
	if notAffected == 0 {
		return fmt.Errorf("delete error")
	}

	return nil
}
