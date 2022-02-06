package user

import (
	"database/sql"
	"errors"
	"eventapp/entities"
	"log"
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
		log.Fatal(err)
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
			log.Fatal(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// edit by Bagus, return repository memakai entity saja
func (r *UserRepository) GetById(id int) (entities.User, error) {
	stmt, err := r.db.Prepare("select id, name, email, phone, avatar from users where id = ? and deleted_at is NULL")

	if err != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	res, err := stmt.Query(id)

	if err != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	defer res.Close()

	var user entities.User

	if res.Next() {
		err := res.Scan(&user.Id, &user.Name, &user.Email, &user.PhoneNumber, &user.Avatar)

		if err != nil {
			log.Fatal(err)
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
func (r *UserRepository) Create(user entities.User) (entities.User, error) {
	err := r.checkEmailExistence(user.Email)

	if err != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	stmt, err := r.db.Prepare("insert into users(name, email, password, phone, avatar) values(?,?,?,?,?)")

	if err != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	res, err := stmt.Exec(user.Name, user.Email, user.Password, user.PhoneNumber, user.Avatar)

	if err != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return entities.User{}, errors.New("user not created")
	}

	id, _ := res.LastInsertId()
	user.Id = int(id)

	return user, nil
}

// edit by Bagus, parameter dan return repository memakai entity saja
// email harus unik, dicek dengan checkEmailExistence
func (r *UserRepository) Update(id int, user entities.User) (entities.User, error) {
	err := r.checkEmailExistence(user.Email)

	if err != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	stmt, err := r.db.Prepare("update users set name= ?, email= ?, password= ?, phone= ?, avatar= ? where id = ? and deleted_at is null")

	if err != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	res, error := stmt.Exec(user.Name, user.Email, user.Password, user.PhoneNumber, user.Avatar, id)

	if error != nil {
		log.Fatal(err)
		return entities.User{}, err
	}

	notAffected, _ := res.RowsAffected()

	if notAffected == 0 {
		return entities.User{}, errors.New("user not updated")
	}

	return user, nil
}

func (r *UserRepository) Delete(id int) error {
	// stmt, err := r.db.Prepare("DELETE from users where id = ?")
	stmt, err := r.db.Prepare("update users set deleted_at = CURRENT_TIMESTAMP where id = ?")

	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) checkEmailExistence(email string) error {
	stmt, err := r.db.Prepare("select count(id) from users where email = ?")

	if err != nil {
		log.Fatal(err)
		return err
	}

	res, err := stmt.Query(email)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer res.Close()

	count := 0

	if res.Next() {
		if err = res.Scan(&count); err != nil {
			return err
		}
	}

	// Detect email duplicate
	if count != 0 {
		return errors.New("email already exists")
	}

	return nil
}
