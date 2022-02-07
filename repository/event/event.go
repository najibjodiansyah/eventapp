package event

import (
	"database/sql"
	"eventapp/entities"
	"fmt"
	"log"
)

type EventRepositeory struct {
	db *sql.DB
}

func New(db *sql.DB) *EventRepositeory {
	return &EventRepositeory{db: db}
}

// sudah dicek
func (r *EventRepositeory) GetAllEvent(page int) ([]entities.Event, int, error) {
	var totalEvent int
	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo 
								from events e join users u on e.hostid = u.id 
								where e.deleted_at IS NULL limit ? offset ?`)

	if err != nil {
		log.Println(err)
		return nil,totalEvent, err
	}

	limit := 5
	offset := (page - 1) * limit

	res, err := stmt.Query(limit, offset)

	if err != nil {
		log.Println(err)
		return nil,totalEvent, err
	}

	defer res.Close()

	var events []entities.Event

	for res.Next() {
		var event entities.Event

		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo)

		if err != nil {
			log.Println(err)
			return nil,totalEvent, err
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

	return events,totalEvent, nil
}

// sudah dicek
func (r *EventRepositeory) GetEventByLocation(location string, page int) ([]entities.Event, int, error) {
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

// edit by bagus, return ditambah event id dan user name
func (r *EventRepositeory) CreateEvent(hostId int, event entities.Event) (entities.Event, error) {
	stmt, err := r.db.Prepare("insert into events (name, category, hostid, host, description, datetime, location, photo, created_at) VALUES(?,?,?,?,?,?,?,?, CURRENT_TIMESTAMP)")

	if err != nil {
		return entities.Event{}, err
	}

	res, err := stmt.Exec(event.Name, event.Category, hostId, event.Host, event.Description, event.Datetime, event.Location, event.Photo)

	if err != nil {
		return entities.Event{}, err
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		return entities.Event{}, fmt.Errorf("event not created")
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

func (r *EventRepositeory) DeleteEvent(eventid int) error {
	stmt, err := r.db.Prepare("update events set deleted_at = CURRENT_TIMESTAMP where id = ?")
	if err != nil {
		log.Println(err)
	}

	_, err = stmt.Exec(eventid)
	if err != nil {
		return err
	}

	return nil
}

// sudah diceck, return ditambah user name
func (r *EventRepositeory) GetEventByKeyword(keyword string, page int) ([]entities.Event, int, error) {
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
func (r *EventRepositeory) GetEventByCategory(category string, page int) ([]entities.Event, int, error) {
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
func (r *EventRepositeory) GetEventByHostId(hostId int) ([]entities.Event, error) {
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

func (r *EventRepositeory) UpdateEvent(event entities.Event) (entities.Event, error) {
	stmt, err := r.db.Prepare("update events set name= ?, category= ?, host= ?, location= ?, description= ?, datetime= ?, photo= ? where id = ? and deleted_at is NULL")

	if err != nil {
		log.Println(err)
		return entities.Event{}, err
	}

	_, err = stmt.Exec(event.Name, event.Category, event.Host, event.Location, event.Description, event.Datetime, event.Photo, event.Id)

	if err != nil {
		log.Println(err)
		return entities.Event{}, err
	}
	fmt.Println("setelah exec",event)
	// rowsAffected, _ := result.RowsAffected()

	// if rowsAffected == 0 {
	// 	return entities.Event{}, fmt.Errorf("update event failed")
	// }
	return event, nil
}

func (r *EventRepositeory) GetEventByEventId(eventId int) (entities.Event, error) {
	stmt, err := r.db.Prepare(`select e.id, e.name, e.category, u.name, e.host, e.description, e.datetime, e.location, e.photo, e.hostid 
								from events e join users u on e.hostid = u.id
								where e.deleted_at IS NULL and e.id = ?`)

	if err != nil {
		log.Println(err)
		return entities.Event{}, err
	}

	res, err := stmt.Query(eventId)

	if err != nil {
		log.Println(err)
		return entities.Event{}, err
	}

	defer res.Close()

	var event entities.Event

	if res.Next() {
		err := res.Scan(&event.Id, &event.Name, &event.Category, &event.UserName, &event.Host, &event.Description, &event.Datetime, &event.Location, &event.Photo, &event.HostId)

		if err != nil {
			log.Println(err)
			return entities.Event{}, err
		}
	}

	return event, nil
}