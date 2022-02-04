package event

import (
	"database/sql"
	"eventapp/entities"
	"fmt"
	"log"
	"time"
)

type EventRepositeory struct {
	db *sql.DB
}

func New(db *sql.DB) *EventRepositeory {
	return &EventRepositeory{db: db}
}

func (r *EventRepositeory) Get() ([]entities.Event, error) {
	var events []entities.Event
	result, err := r.db.Query("select events.id, events.name, events.category, events.location, events.description, events.photo from events join users on events.userid = users.id")
	if err != nil {
		return events, err
	}

	defer result.Close()

	for result.Next() {
		var event entities.Event
		err := result.Scan(&event.ID, &event.Name, &event.ID, &event.Name, &event.Category, &event.Location, &event.Description, &event.Photo)
		if err != nil {
			return events, fmt.Errorf("event not found")
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepositeory) GetbyId(id int) (entities.Event, error) {
	var event entities.Event
	// stmt, err := r.db.Prepare("select id, name, email,password, organization, phone, avatar from users where id = ?")
	// if err != nil {
	// 	//log.Fatal(err)
	// 	return user, fmt.Errorf("gagal prepare db")
	// }

	// result, err := stmt.Query(id)
	// if err != nil {
	// 	return user, fmt.Errorf("gagal query user")
	// }

	// defer result.Close()

	// for result.Next() {
	// 	err := result.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Organization, &user.PhoneNumber, &user.Avatar)
	// 	if err != nil {
	// 		return user, err
	// 	}
	// 	return user, nil
	// }

	return event, fmt.Errorf("user not found")
}

func (r *EventRepositeory) Create(user entities.Event) (entities.Event, error) {
	// stmt, err := r.db.Prepare("INSERT INTO users(name, email, password, organization, phone, avatar) VALUES(?,?,?,?,?,?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// result, err := stmt.Exec(user.Name, user.Email, user.Password, user.Organization, user.PhoneNumber, user.Avatar)
	// if err != nil {
	// 	return user, fmt.Errorf("gagal exec")
	// }

	// notAffected, _ := result.RowsAffected()
	// if notAffected == 0 {
	// 	return user, fmt.Errorf("user not created")
	// }

	return user, nil
}

func (r *EventRepositeory) Update(id int, user entities.Event) (entities.Event, error) {
	// stmt, err := r.db.Prepare("UPDATE users SET name= ?, email= ?, password= ?, organization= ?, phone= ?, avatar= ? WHERE id = ?")
	// if err != nil {
	// 	// log.Fatal(err)
	// 	return user, fmt.Errorf("gagal prepare update")
	// }

	// result, error := stmt.Exec(user.Name, user.Email, user.Password, user.Organization, user.PhoneNumber, user.Avatar, id)
	// if error != nil {
	// 	return user, fmt.Errorf("gagal exec update")
	// }

	// notAffected, _ := result.RowsAffected()
	// if notAffected == 0 {
	// 	return user, fmt.Errorf("row not affected")
	// }

	return user, nil
}

func (r *EventRepositeory) Delete(id int) error {
	stmt, err := r.db.Prepare("UPDATE events SET deleted_at= ? where id = ?")
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Exec(time.Now(), id)
	if err != nil {
		return err
	}

	notAffected, _ := result.RowsAffected()
	if notAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
