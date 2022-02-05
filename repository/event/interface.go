package event

import "eventapp/entities"

type Event interface {
	GetAllEvent(page int) ([]entities.Event, error)
	GetEventByLocation(location string, page int) ([]entities.Event, error)
	GetEventByKeyword(keyword string, page int) ([]entities.Event, error) //done
	GetEventByCategory(category string, page int) ([]entities.Event, error) //done
	GetbyHostId(userid int) ([]entities.Event, error) //done
	CreateEvent(entities.Event) (entities.Event, error)
	Update(id int, event entities.Event) (entities.Event, error)
	Delete(id int) error //done
}
