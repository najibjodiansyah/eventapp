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
		stmt, err := r.db.Prepare("select id, name, userid, host, category, date, location, description, photo from events where id = ? deleted_at is null ")
		if err != nil {
			//log.Fatal(err)
			return event, fmt.Errorf("gagal prepare db")
		}
	
		result, err := stmt.Query(id)
		if err != nil {
			return event, fmt.Errorf("gagal query event")
		}
	
		defer result.Close()
	
		for result.Next() {
			err := result.Scan(&event.ID, &event.Name, &event.UserId, &event.Host, &event.Category, &event.Datetime, &event.Location, &event.Description, &event.Photo)
			if err != nil {
				return event, err
			}
			return event, nil
		}
	
		return event, fmt.Errorf("event not found")
	}

func (r *EventRepositeory) Create(user entities.Event) (entities.Event, error) {
	
	return user, nil
}

func (r *EventRepositeory) Update(id int, user entities.Event) (entities.Event, error) {
	

	return user, nil
}

func (r *EventRepositeory) Delete(id int) error {
	stmt, err := r.db.Prepare("UPDATE events SET deleted_at = ? where id = ?")
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Exec(time.Now(),id)
	if err != nil {
		return err
	}

	notAffected, _ := result.RowsAffected()
	if notAffected == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

func (r *EventRepositeory) GetEventByKeyword(keyword string, page int) ([]entities.Event, error) {
	stmt, err := r.db.Prepare("select id, name, userid, host, category, date, location, description, photo from events where deleted_at is null and keyword like %?% limit ? offset ? ")
	if err != nil {
		log.Fatal(err)
	}

	var events []entities.Event

	result, errr := stmt.Query()
	if errr != nil {
		return events, err
	}

	defer result.Close()

	for result.Next() {
		var event entities.Event
		limit := 10
		offset := (page - 1) * limit
		err := result.Scan(&event.ID, &event.Name, &event.UserId, &event.Host, &event.Category, &event.Datetime, &event.Location, &event.Description, &event.Photo, keyword, limit, offset)
		if err != nil {
			log.Fatal("error di scan getEvent")
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepositeory) GetEventByCategory(category string, page int) ([]entities.Event, error) {
	stmt, err := r.db.Prepare("select id, name, userid, host, category, date, location, description, photo from events where deleted_at is null and category = ? limit ? offset ?")
	if err != nil {
		log.Fatal(err)
	}

	var events []entities.Event
	
	limit := 5

	offset := (page - 1) * limit

	result, errr := stmt.Query(category,limit,offset)
	if errr != nil {
		return events, err
	}

	defer result.Close()

	for result.Next() {
		var event entities.Event
		err := result.Scan(&event.ID, &event.Name, &event.UserId, &event.Host, &event.Category, &event.Datetime, &event.Location, &event.Description, &event.Photo)
		if err != nil {
			log.Fatal("error di scan getEvent")
		}
		
		events = append(events, event)
	}
	fmt.Println(events)
	return events, nil
}

func (r *EventRepositeory) GetEventByLocation(location string, page int) ([]entities.Event, error) {
	
	return nil, nil
}
