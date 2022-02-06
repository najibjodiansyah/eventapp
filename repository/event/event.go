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

// sudah dicek
func (r *EventRepositeory) GetAllEvent(page int) ([]entities.Event, error) {
	var events []entities.Event

	stmt, err := r.db.Prepare(` select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id 
								where e.deleted_at IS NULL limit ? offset ?`)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	limit := 5
	offset := (page - 1) * limit

	res, err := stmt.Query(limit, offset)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// sudah dicek
func (r *EventRepositeory) GetEventByLocation(location string, page int) ([]entities.Event, error) {
	var events []entities.Event

	stmt, err := r.db.Prepare(`	select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id
								where e.location = ? and e.deleted_at IS NULL limit ? offset ?`)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	limit := 5
	offset := (page - 1) * limit

	res, errr := stmt.Query(location, limit, offset)

	if errr != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// edit by bagus, return ditambah event id dan user name
func (r *EventRepositeory) CreateEvent(hostId int, event entities.Event) (entities.Event, error) {
	stmt, err := r.db.Prepare("insert into events (name, category, hostid, host, description, datetime, location, photo) VALUES(?,?,?,?,?,?,?,?)")

	if err != nil {
		return event, err
	}

	res, err := stmt.Exec(event.Name, event.Category, hostId, event.Host, event.Description, event.Datetime, event.Location, event.Photo)

	if err != nil {
		return event, fmt.Errorf("event not created")
	}

	RowsAffected, _ := res.RowsAffected()

	if RowsAffected == 0 {
		return event, fmt.Errorf("event not created")
	}

	id, _ := res.LastInsertId()
	event.Id = int(id)

	stmt, _ = r.db.Prepare("select name from users where id = ?")

	row, _ := stmt.Query(hostId)

	defer row.Close()

	if row.Next() {
		var name string

		row.Scan(&name)

		event.UserName = name
	}

	return event, nil
}

func (r *EventRepositeory) DeleteEvent(id int) error {
	stmt, err := r.db.Prepare("UPDATE events SET deleted_at= ? where id = ?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

// sudah diceck, return ditambah user name
func (r *EventRepositeory) GetEventByKeyword(keyword string, page int) ([]entities.Event, error) {
	var events []entities.Event

	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id
								where e.deleted_at is null and e.name like ? limit ? offset ?`)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	like := "%" + keyword + "%"
	limit := 10
	offset := (page - 1) * limit

	res, err := stmt.Query(like, limit, offset)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// sudah dicek, return ditambah user name
func (r *EventRepositeory) GetEventByCategory(category string, page int) ([]entities.Event, error) {
	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id
								where e.deleted_at is null and e.category = ? limit ? offset ?`)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var events []entities.Event

	limit := 5
	offset := (page - 1) * limit

	res, err := stmt.Query(category, limit, offset)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// sudah dicek
func (r *EventRepositeory) GetEventByHostId(hostId int) ([]entities.Event, error) {
	stmt, err := r.db.Prepare("select events.id, events.name, events.userid, events.host, events.category, events.datetime, events.location, events.description, events.photo, users.name from events join users where events.hostid = users.id and events.deleted_at is null and events.hostid = ?")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var events []entities.Event

	res, err := stmt.Query(hostId)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.HostId, &event.Host, &event.Category, &event.Datetime, &event.Location, &event.Description, &event.Photo, &event.UserName)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepositeory) UpdateEvent(eventid int, event entities.Event) (entities.Event, error) {
	stmt, err := r.db.Prepare("UPDATE events SET name= ?, category= ?, host= ?, location= ?,description= ?, datetime=?,photo=? WHERE id = ? AND deleted_at is NULL")

	if err != nil {
		log.Fatal(err)
		return event, fmt.Errorf("update event failed")
	}

	result, error := stmt.Exec(event.Name, event.Category, event.Host, event.Location, event.Description, event.Datetime, event.Photo, eventid)

	if error != nil {
		return event, fmt.Errorf("update event failed")
	}

	RowsAffected, _ := result.RowsAffected()

	if RowsAffected == 0 {
		return event, fmt.Errorf("update event failed")
	}

	return event, nil
}

func (r *EventRepositeory) GetById(eventId int) (entities.Event, error) {
	var event entities.Event

	stmt, err := r.db.Prepare(` select name, hostid, category, host, description, datetime, location, photo 
								from events
								where deleted_at IS NULL and id = ?`)

	if err != nil {
		log.Fatal(err)
		return event, err
	}

	res, err := stmt.Query(eventId)

	if err != nil {
		return event, err
	}

	defer res.Close()

	if res.Next() {
		err := res.Scan(&event.Name, &event.HostId, &event.Category, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			return event, err
		}
	}

	return event, nil
}
