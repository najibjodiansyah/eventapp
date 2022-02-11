package event

import (
	"database/sql"
	"errors"
	"eventapp/entities"
	"fmt"
	"log"
	"net/http"
)

type EventRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// sudah dicek
func (r *EventRepository) GetAllEvent(page int) ([]entities.Event, int, error) {
	var totalEvent int
	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id 
								where e.deleted_at IS NULL limit ? offset ?`)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	limit := 5
	offset := (page - 1) * limit

	res, err := stmt.Query(limit, offset)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	defer res.Close()

	var events []entities.Event

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Println(err)
			return nil, totalEvent, err
		}

		events = append(events, event)
	}

	stmt2, err := r.db.Prepare(`select count(e.id)
								from events e 
								where e.deleted_at is null
								`)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	res2, err := stmt2.Query()

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	defer res2.Close()

	for res2.Next() {

		err := res2.Scan(&totalEvent)

		if err != nil {
			log.Println(err)
			return nil, totalEvent, err
		}

	}

	return events, totalEvent, nil
}

// sudah dicek
func (r *EventRepository) GetEventByLocation(location string, page int) ([]entities.Event, int, error) {
	var totalEvent int
	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id
								where e.location = ? and e.deleted_at IS NULL 
								limit ? offset ?`)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	limit := 5
	offset := (page - 1) * limit

	res, errr := stmt.Query(location, limit, offset)

	if errr != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	stmt2, err := r.db.Prepare("select count(id) from events where location = ? and deleted_at IS NULL")
	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}
	res2, err2 := stmt2.Query(location)
	if err2 != nil {
		log.Println(err2)
		return nil, totalEvent, err2
	}
	if errr != nil {
		log.Println(err)
		return nil, totalEvent, err
	}
	fmt.Println("sampe sini jalan")
	defer res.Close()
	defer res2.Close()

	var events []entities.Event

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Println(err)
			return nil, totalEvent, err
		}
		// fmt.Println(totalEvent)

		events = append(events, event)
	}
	fmt.Println("sampe sini jalan2")
	for res2.Next() {

		err := res2.Scan(&totalEvent)

		if err != nil {
			log.Println(err)
			return nil, totalEvent, err
		}
		fmt.Println(totalEvent)
	}
	fmt.Println("sampe sini finish")

	return events, totalEvent, nil
}

func (r *EventRepository) CreateEvent(event entities.Event) (createdEvent entities.Event, code int, err error) {
	id, err := r.checkUserExistence(event.HostId)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdEvent, code, err
	}

	if id == 0 {
		log.Println("user no longer exist")
		code, err = http.StatusBadRequest, errors.New("user no longer exist")
		return createdEvent, code, err
	}

	stmt, err := r.db.Prepare("insert into events (name, host, description, datetime, location, category, photo, hostid, created_at) VALUES(?,?,?,?,?,?,?,?, CURRENT_TIMESTAMP)")

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdEvent, code, err
	}

	res, err := stmt.Exec(event.Name, event.Host, event.Description, event.Datetime, event.Location, event.Category, event.Photo, event.HostId)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdEvent, code, err
	}

	id, err = res.LastInsertId()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdEvent, code, err
	}

	event.Id = int(id)

	stmt, err = r.db.Prepare("select name from users where id = ?")

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdEvent, code, err
	}

	row, err := stmt.Query(event.HostId)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return createdEvent, code, err
	}

	defer row.Close()

	if row.Next() {
		var name string

		row.Scan(&name)

		event.UserName = name
	}

	return event, http.StatusOK, nil
}

func (r *EventRepository) DeleteEvent(eventid int) (code int, err error) {
	stmt, err := r.db.Prepare("update events set deleted_at = CURRENT_TIMESTAMP where id = ?")

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return code, err
	}

	res, err := stmt.Exec(eventid)

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
		log.Println("rows affected is 0 while delete event")
		code, err = http.StatusBadRequest, errors.New("event not deleted")
		return code, err
	}

	return http.StatusOK, nil
}

// sudah diceck, return ditambah user name
func (r *EventRepository) GetEventByKeyword(keyword string, page int) ([]entities.Event, int, error) {
	var totalEvent int
	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id
								where e.deleted_at is null and e.name like ? limit ? offset ?`)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	like := "%" + keyword + "%"
	limit := 5
	offset := (page - 1) * limit

	res, err := stmt.Query(like, limit, offset)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	defer res.Close()

	var events []entities.Event

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Println(err)
			return nil, totalEvent, err
		}

		events = append(events, event)
	}

	stmt2, err := r.db.Prepare(`select count(e.id)
								from events e 
								where e.deleted_at is null and e.name like ? 
								`)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	like2 := "%" + keyword + "%"

	res2, err := stmt2.Query(like2)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	defer res2.Close()

	for res2.Next() {

		err := res2.Scan(&totalEvent)

		if err != nil {
			log.Println(err)
			return nil, totalEvent, err
		}

	}
	return events, totalEvent, nil
}

// sudah dicek, return ditambah user name
func (r *EventRepository) GetEventByCategory(category string, page int) ([]entities.Event, int, error) {
	var totalEvent int
	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id
								where e.deleted_at is null and e.category = ? limit ? offset ?`)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	limit := 5
	offset := (page - 1) * limit

	res, err := stmt.Query(category, limit, offset)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	defer res.Close()

	var events []entities.Event

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Println(err)
			return nil, totalEvent, err
		}

		events = append(events, event)
	}

	stmt2, err := r.db.Prepare(`select count(e.id)
								from events e 
								where e.deleted_at is null and e.category = ? 
								`)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	res2, err := stmt2.Query(category)

	if err != nil {
		log.Println(err)
		return nil, totalEvent, err
	}

	defer res2.Close()

	for res2.Next() {

		err := res2.Scan(&totalEvent)

		if err != nil {
			log.Println(err)
			return nil, totalEvent, err
		}

	}
	return events, totalEvent, nil
}

// sudah dicek
func (r *EventRepository) GetEventByHostId(hostId int) ([]entities.Event, error) {
	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id
								where e.deleted_at is null and e.hostid = ?`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := stmt.Query(hostId)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer res.Close()

	var events []entities.Event

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) UpdateEvent(event entities.Event) (updatedEvent entities.Event, code int, err error) {
	stmt, err := r.db.Prepare("update events set name=?, host=?, description=?, datetime=?, location=?, category=?, photo=? where id=? and deleted_at is NULL")

	if err != nil {
		log.Println(err)
		code, err := http.StatusInternalServerError, errors.New("internal server error")
		return updatedEvent, code, err
	}

	res, err := stmt.Exec(event.Name, event.Host, event.Description, event.Datetime, event.Location, event.Category, event.Photo, event.Id)

	if err != nil {
		log.Println(err)
		code, err := http.StatusInternalServerError, errors.New("internal server error")
		return updatedEvent, code, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return updatedEvent, code, err
	}

	if rowsAffected == 0 {
		log.Println("rows affected is 0 while update event")
		code, err = http.StatusBadRequest, errors.New("event not updated")
		return updatedEvent, code, err
	}

	updatedEvent = event

	return updatedEvent, http.StatusOK, nil
}

func (r *EventRepository) GetEventByEventId(eventId int) (event entities.Event, err error) {
	stmt, err := r.db.Prepare(` select e.id, e.name, e.host, e.description, e.datetime, e.location, e.category, e.photo, e.hostid, u.name 
								from events e join users u on e.hostid = u.id
								where e.deleted_at IS NULL and e.id = ?`)

	if err != nil {
		log.Println(err)
		return event, errors.New("internal server error")
	}

	res, err := stmt.Query(eventId)

	if err != nil {
		log.Println(err)
		return event, errors.New("internal server error")
	}

	defer res.Close()

	if res.Next() {
		err := res.Scan(&event.Id, &event.Name, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Category, &event.Photo, &event.HostId, &event.UserName)

		if err != nil {
			log.Println(err)
			return event, errors.New("internal server error")
		}
	}

	return event, nil
}

func (r *EventRepository) checkUserExistence(userId int) (id int64, err error) {
	stmt, err := r.db.Prepare("select id from users where deleted_at is null and id = ?")

	if err != nil {
		return 0, err
	}

	res, err := stmt.Query(userId)

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
