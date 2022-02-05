package participant

import (
	"eventapp/entities"
)

type Participant interface {
	GetParticipantsByEventId(eventId int) ([]entities.Participant, error)
	GetEventsByParticipantId(userId int) ([]entities.Event, error)
	JoinEvent(eventId int, userId int) error
	UnjoinEvent(eventId int, userId int) error
}
