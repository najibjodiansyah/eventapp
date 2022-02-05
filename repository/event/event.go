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

func (r *EventRepositeory) GetAllEvent(page int) ([]entities.Event, error) {
	var events []entities.Event
	result, err := r.db.Query(`	select events.id, events.name, events.category, events.host, 
							   	events.description, events.datetime,events.location,events.photo 
							   	from events 
								where location = ? and deleted_at IS NULL limit ?`, page)
	if err != nil {
		return events, err
	}

	defer result.Close()

	for result.Next() {
		var event entities.Event
		err := result.Scan(&event.ID, &event.Name, &event.Category, &event.Host,
			&event.Description, &event.Datetime, &event.Location, &event.Photo)
		if err != nil {
			return events, fmt.Errorf("event not found")
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepositeory) GetEventByLocation(location string, page int) ([]entities.Event, error) {
	var eventsLocation []entities.Event
	eventByLocation, err := r.db.Query(`select events.id, events.name, events.category, events.host, 
									events.description, events.datetime,events.location,events.photo 
									from events 
									where location = ? and deleted_at IS NULL limit ?`, location, page)
	if err != nil {
		return eventsLocation, err
	}
	defer eventByLocation.Close()

	for eventByLocation.Next() {
		var eventLocation entities.Event

		err := eventByLocation.Scan(&eventLocation.ID, &eventLocation.Name,
			&eventLocation.Category, &eventLocation.Host,
			&eventLocation.Description, &eventLocation.Datetime,
			&eventLocation.Location, &eventLocation.Photo)
		if err != nil {
			return eventsLocation, err
		}
		eventsLocation = append(eventsLocation, eventLocation)
	}
	if location == "" {
		return eventsLocation, fmt.Errorf("lokasi tidak ada")
	}
	return eventsLocation, nil
}

func (r *EventRepositeory) GetEventByCategory(category string, page int) ([]entities.Event, error) {
	return nil, nil
}

func (r *EventRepositeory) CreateEvent(event entities.Event) (entities.Event, error) {
	stmt, err := r.db.Prepare(`insert into events 
							   (name,category,host,location,description,datetime,location,photo) 
							   VALUES(?,?,?,?,?,?,,?,?)`)
	if err != nil {
		return event, err
	}
	result, err := stmt.Exec(event.Name, event.Category, event.Host, event.Location,
		event.Description, event.Datetime, event.Location, event.Photo, event.CreatedAt)
	if err != nil {
		return event, fmt.Errorf("gagal exec")
	}
	RowsAffected, _ := result.RowsAffected()
	if RowsAffected == 0 {
		return event, fmt.Errorf("event not created")
	}
	return event, nil
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
