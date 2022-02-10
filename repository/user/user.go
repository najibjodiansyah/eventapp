package user

import (
	_entities "eventapp/entities"

	"database/sql"
	"errors"
	"log"
	"net/http"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAllUsers() ([]_entities.User, error) {
	stmt, err := r.db.Prepare("select id, name, email, phone, avatar from users where deleted_at is NULL")

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	var users []_entities.User

	result, err := stmt.Query()

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	defer result.Close()

	for result.Next() {
		var user _entities.User

		err := result.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Avatar)

		if err != nil {
			log.Println(err)
			return nil, errors.New("internal server error")
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUserById(id int) (user _entities.User, err error) {
	stmt, err := r.db.Prepare("select id, name, email, password, phone, avatar from users where id = ? and deleted_at is NULL")

	if err != nil {
		log.Println(err)
		return user, errors.New("internal server error")
	}

	res, err := stmt.Query(id)

	if err != nil {
		log.Println(err)
		return user, errors.New("internal server error")
	}

	defer res.Close()

	if res.Next() {
		err := res.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Avatar)

		if err != nil {
			log.Println(err)
			return user, errors.New("internal server error")
		}
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user _entities.User) (createdUser _entities.User, code int, err error) {
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

	stmt, err := r.db.Prepare("insert into users(name, email, password, phone, avatar, created_at) values(?,?,?,'','',CURRENT_TIMESTAMP)")

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

func (r *UserRepository) UpdateUser(user _entities.User) (updatedUser _entities.User, code int, err error) {
	id, err := r.checkEmailExistence(user.Email)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return updatedUser, code, err
	}

	if id != 0 && id != int64(user.Id) {
		log.Println("email already exist while create user")
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
		code, err = http.StatusBadRequest, errors.New("user not updated")
		return updatedUser, code, err
	}

	updatedUser = user

	return updatedUser, http.StatusOK, nil
}

func (r *UserRepository) DeleteUser(id int) (code int, err error) {
	stmt, err := r.db.Prepare("update users set deleted_at = CURRENT_TIMESTAMP where deleted_at is NULL and id = ?")

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return code, err
	}

	res, err := stmt.Exec(id)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return code, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return code, err
	}

	if rowsAffected == 0 {
		log.Println("rows affected is 0 while delete user")
		code, err = http.StatusBadRequest, errors.New("user not deleted")
		return code, err
	}

	return http.StatusOK, nil
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
