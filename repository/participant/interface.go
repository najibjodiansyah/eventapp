package participant

import (
	"eventapp/entities"
)

type Participant interface {
	GetParticipantsByEventId(eventId int, page int) ([]entities.Participant, error)
	JoinEvent(eventId int, userId int) error
	UnjoinEvent(eventId int, userId int) error
}
