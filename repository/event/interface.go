package event

import "eventapp/entities"

type Event interface {
	Get() ([]entities.Event, error)
	GetEventByLocation(location string, page int) ([]entities.Event, error)
	GetEventByKeyword(keyword string, page int) ([]entities.Event, error) //
	GetEventByCategory(category string, page int) ([]entities.Event, error) //
	GetbyId(id int) (entities.Event, error) //
	Create(entities.Event) (entities.Event, error)
	Update(id int, event entities.Event) (entities.Event, error)
	Delete(id int) error //
}
