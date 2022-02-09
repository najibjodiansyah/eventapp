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

		err := result.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Avatar)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUserById(id int) (user entities.User, code int, err error) {
	stmt, err := r.db.Prepare("select name, email, password, phone, avatar from users where id = ? and deleted_at is NULL")

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return user, code, err
	}

	res, err := stmt.Query(id)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return user, code, err
	}

	defer res.Close()

	if res.Next() {
		err := res.Scan(&user.Name, &user.Email, &user.Password, &user.Phone, &user.Avatar)

		if err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return user, code, err
		}
	}

	if user == (entities.User{}) {
		log.Println("user unknown while get user by id")
		code, err := http.StatusBadRequest, errors.New("user unknown")
		return user, code, err
	}

	return user, http.StatusOK, nil
}

func (r *UserRepository) CreateUser(user entities.User) (createdUser entities.User, code int, err error) {
	id, err := r.checkEmailExistence(user.Email)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdUser, code, err
	}

	if id != 0 {
		log.Println("email already exist while create user")
		code, err = http.StatusBadRequest, errors.New("email already exist")
		return createdUser, code, err
	}

	stmt, err := r.db.Prepare("insert into users(name, email, password) values(?,?,?)")

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdUser, code, err
	}

	res, err := stmt.Exec(user.Name, user.Email, user.Password)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdUser, code, err
	}

	id, err = res.LastInsertId()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdUser, code, err
	}

	createdUser = user
	createdUser.Id = int(id)

	return createdUser, http.StatusOK, nil
}

func (r *UserRepository) Update(user entities.User) (updatedUser entities.User, code int, err error) {
	id, err := r.checkEmailExistence(user.Email)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return updatedUser, code, err
	}

	if id != int64(user.Id) {
		log.Println(id, "email already exist while create user")
		code, err = http.StatusBadRequest, errors.New("email already exist")
		return updatedUser, code, err
	}

	stmt, err := r.db.Prepare("update users set name= ?, email= ?, password= ?, phone= ?, avatar= ? where id = ? and deleted_at is null")

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return updatedUser, code, err
	}

	res, err := stmt.Exec(user.Name, user.Email, user.Password, user.Phone, user.Avatar, user.Id)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return updatedUser, code, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return updatedUser, code, err
	}

	if rowsAffected == 0 {
		log.Println("rows affected is 0 while update user")
		code, err = http.StatusInternalServerError, errors.New("user not updated")
		return updatedUser, code, err
	}

	updatedUser = user

	return updatedUser, code, err
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

func (r *UserRepository) checkEmailExistence(email string) (id int64, err error) {
	stmt, err := r.db.Prepare("select id from users where email = ?")

	if err != nil {
		return 0, err
	}

	res, err := stmt.Query(email)

	if err != nil {
		return 0, err
	}

	defer res.Close()

	if res.Next() {
		if err = res.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}
