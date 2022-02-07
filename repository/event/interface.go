package event

import "eventapp/entities"

type Event interface {
	GetEventByEventId(eventid int) (entities.Event, error)
	GetAllEvent(page int) ([]entities.Event, error)
	GetEventByLocation(location string, page int) ([]entities.Event, int, error)
	GetEventByKeyword(keyword string, page int) ([]entities.Event, int, error)
	GetEventByCategory(category string, page int) ([]entities.Event, int, error)
	GetEventByHostId(hostid int) ([]entities.Event, error)
	CreateEvent(hostId int, event entities.Event) (entities.Event, error)
	UpdateEvent(event entities.Event) (entities.Event, error)
	DeleteEvent(eventid int) error
}
