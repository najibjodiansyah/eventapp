package user

import (
	"database/sql"
	"errors"
	"eventapp/entities"
	"log"
	"net/http"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// edit by Bagus, return repository memakai entity saja
func (r *UserRepository) Get() ([]entities.User, error) {
	stmt, err := r.db.Prepare("select id, name, email, phone, avatar from users where deleted_at is NULL")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var users []entities.User

	result, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		var user entities.User

		err := result.Scan(&user.Id, &user.Name, &user.Email, &user.PhoneNumber, &user.Avatar)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// edit by Bagus, return repository memakai entity saja
func (r *UserRepository) GetById(id int) (entities.User, error) {
	stmt, err := r.db.Prepare("select id, name, email, password, phone, avatar from users where id = ? and deleted_at is NULL")

	if err != nil {
		log.Println(err)
		return entities.User{}, err
	}

	res, err := stmt.Query(id)

	if err != nil {
		log.Println(err)
		return entities.User{}, err
	}

	defer res.Close()

	var user entities.User

	if res.Next() {
		err := res.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.PhoneNumber, &user.Avatar)

		if err != nil {
			log.Println(err)
			return entities.User{}, err
		}
	}

	if user == (entities.User{}) {
		return entities.User{}, errors.New("user not found")
	}

	return user, nil
}

// edit by Bagus, return repository memakai entity saja
// ditambah user id
// email harus unik, dicek dengan checkEmailExistence
func (r *UserRepository) Create(user entities.User) (createdUser entities.User, code int, err error) {
	code, err = r.checkEmailExistence(user.Email)

	if err != nil {
		return entities.User{}, code, err
	}

	stmt, err := r.db.Prepare("insert into users(name, email, password) values(?,?,?)")

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return entities.User{}, code, err
	}

	res, err := stmt.Exec(user.Name, user.Email, user.Password)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return entities.User{}, code, err
	}

	id, _ := res.LastInsertId()
	user.Id = int(id)

	return user, http.StatusOK, nil
}

// edit by Bagus, parameter dan return repository memakai entity saja
func (r *UserRepository) Update(id int, user entities.User) (entities.User, error) {

	stmt, err := r.db.Prepare("update users set name= ?, email= ?, password= ?, phone= ?, avatar= ? where id = ? and deleted_at is null")

	if err != nil {
		log.Println(err)
		return entities.User{}, err
	}

	_, error := stmt.Exec(user.Name, user.Email, user.Password, user.PhoneNumber, user.Avatar, id)

	if error != nil {
		log.Println(err)
		return entities.User{}, err
	}

	// rowsAffected, _ := res.RowsAffected()

	// if rowsAffected == 0 {
	// 	return entities.User{}, errors.New("user not updated")
	// }

	return user, nil
}

func (r *UserRepository) Delete(id int) error {
	// stmt, err := r.db.Prepare("DELETE from users where id = ?")
	stmt, err := r.db.Prepare("update users set deleted_at = CURRENT_TIMESTAMP where id = ?")

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) checkEmailExistence(email string) (code int, err error) {
	stmt, err := r.db.Prepare("select count(id) from users where email = ?")

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, errors.New("internal server error")
	}

	res, err := stmt.Query(email)

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, errors.New("internal server error")
	}

	defer res.Close()

	count := 0

	if res.Next() {
		if err = res.Scan(&count); err != nil {
			log.Println(err)
			return http.StatusInternalServerError, errors.New("internal server error")
		}
	}

	// Detect email duplicate
	if count != 0 {
		return http.StatusBadRequest, errors.New("email already exists")
	}

	return http.StatusOK, nil
}
