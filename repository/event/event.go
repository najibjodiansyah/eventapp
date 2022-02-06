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
	stmt, err := r.db.Prepare(` select events.id, events.name, events.category, events.host, 
								events.description, events.datetime,events.location, events.photo 
								from events 
								where deleted_at IS NULL limit ? offset ?`)
	if err != nil {
		log.Fatal(err)
	}
	limit := 5

	offset := (page - 1) * limit

	result, errr := stmt.Query(limit, offset)
	if errr != nil {
		return events, err
	}
	defer result.Close()

	for result.Next() {
		var event entities.Event
		err := result.Scan(&event.ID, &event.Name, &event.Category, &event.Host,
			&event.Description, &event.Datetime, &event.Location, &event.Photo)
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepositeory) GetEventByLocation(location string, page int) ([]entities.Event, error) {
	var eventsLocation []entities.Event
	stmt, err := r.db.Prepare(`	select events.id, events.name, events.category, events.host, 
								events.description, events.datetime,events.location,events.photo 
								from events 
								where location = ? and deleted_at IS NULL limit ? offset ?`)

	if err != nil {
		log.Fatal(err)
	}
	limit := 5

	offset := (page - 1) * limit

	result, errr := stmt.Query(location, limit, offset)
	if errr != nil {
		return eventsLocation, err
	}
	defer result.Close()

	for result.Next() {
		var eventLocation entities.Event

		err := result.Scan(&eventLocation.ID, &eventLocation.Name,
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

func (r *EventRepositeory) CreateEvent(event entities.Event) (entities.Event, error) {
	stmt, err := r.db.Prepare("insert into events (name,category,host,description,datetime,location,photo) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return event, err
	}
	result, err := stmt.Exec(event.Name, event.Category, event.Host,
		event.Description, event.Datetime, event.Location, event.Photo)
	if err != nil {
		fmt.Println("eror ya", err)
		return event, fmt.Errorf("gagal exec")

	}
	fmt.Println(result)
	RowsAffected, _ := result.RowsAffected()
	if RowsAffected == 0 {
		return event, fmt.Errorf("event not created")
	}
	return event, nil
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
	stmt, err := r.db.Prepare("select events.id, events.name, events.host, events.category, events.date, events.location, events.description, events.photo, users.name from events join users where events.userid = users.id and events.deleted_at is null and events.name like ? limit ? offset ?")
	if err != nil {
		log.Fatal(err)
	}

	var events []entities.Event
	like := "%"+keyword+"%"
	limit := 10
	offset := (page - 1) * limit
	result, errr := stmt.Query(like, limit, offset)
	if errr != nil {
		return events, err
	}

	defer result.Close()

	for result.Next() {
		var event entities.Event
		var user entities.User
		err := result.Scan(&event.ID, &event.Name, &event.Host, &event.Category, &event.Datetime, &event.Location, &event.Description, &event.Photo, &user.Name)
		if err != nil {
			log.Fatal("error di scan getEvent")
		}

		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepositeory) GetEventByCategory(category string, page int) ([]entities.Event, error) {
	stmt, err := r.db.Prepare("select events.id, events.name, events.host, events.category, events.date, events.location, events.description, events.photo, users.name from events join users where events.userid = users.id and events.deleted_at is null and category = ? limit ? offset ?")
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
		var user entities.User
		err := result.Scan(&event.ID, &event.Name, &event.Host, &event.Category, &event.Datetime, &event.Location, &event.Description, &event.Photo, &user.Name)
		if err != nil {
			log.Fatal("error di scan getEvent")
		}
		
		events = append(events, event)
	}
	fmt.Println(events)
	return events, nil
}


func (r *EventRepositeory) GetbyHostId(userid int) ([]entities.Event, error)  {
	stmt, err := r.db.Prepare("select events.id, events.name, events.host, events.category, events.date, events.location, events.description, events.photo, users.name from events join users where events.userid = users.id and events.deleted_at is null and events.userid = ?")
	if err != nil {
		log.Fatal(err)
	}

	var events []entities.Event

	result, errr := stmt.Query(userid)
	if errr != nil {
		return events, err
	}

	defer result.Close()

	for result.Next() {
		var event entities.Event
		var user entities.User
		err := result.Scan(&event.ID, &event.Name, &event.Host, &event.Category, &event.Datetime, &event.Location, &event.Description, &event.Photo, &user.Name)
		if err != nil {
			log.Fatal("error di scan getEvent")
		}
		
		events = append(events, event)
	}
	fmt.Println(events)
	return events, nil
}

// func (r *EventRepositeory) GetEventByCategory(category string, page int) ([]entities.Event, error) {
// 	var event []entities.Event

// 	return event, nil
// }
func (r *EventRepositeory) GetbyId(id int) (entities.Event, error) {
	var event entities.Event

	return event, nil
}
func (r *EventRepositeory) GetUpdateEvent(id int, event entities.Event) (entities.Event, error) {
	stmt, err := r.db.Prepare("UPDATE events SET name= ?, category= ?, host= ?, location= ?,description= ?, datetime=?,photo=? WHERE id = ? AND deleted_at is NULL")
	if err != nil {
		// log.Fatal(err)
		return event, fmt.Errorf("gagal prepare update")
	}

	result, error := stmt.Exec(event.Name, event.Category, event.Host, event.Location,
		event.Description, event.Datetime, event.Photo, id)
	if error != nil {
		return event, fmt.Errorf("gagal exec update")
	}

	notAffected, _ := result.RowsAffected()
	if notAffected == 0 {
		return event, fmt.Errorf("row not affected")
	}
	return event, nil
}

// GetAllEvent(page int) ([]entities.Event, error)
// 	GetEventByLocation(location string, page int) ([]entities.Event, error)
// 	GetEventByKeyword(keyword string, page int) (entities.Event, error)
// 	GetEventByCategory(category string, page int) (entities.Event, error)
// 	GetbyId(id int) (entities.Event, error)
// 	CreateEvent(entities.Event) (entities.Event, error)
// 	Update(id int, event entities.Event) (entities.Event, error)
// 	Delete(id int) error
// }
