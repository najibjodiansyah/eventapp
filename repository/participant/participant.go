package participant

import (
	"database/sql"
	"errors"
	"eventapp/entities"
	"log"
)

type ParticipantRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

func (r *ParticipantRepository) GetParticipantsByEventId(eventId int) ([]entities.Participant, error) {
	stmt, err := r.db.Prepare(`select u.name, u.avatar
								from users u left join participants p on p.participantid = u.id
								where p.deleted_at is NULL and p.eventid = ?`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := stmt.Query(eventId)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer result.Close()

	var participants []entities.Participant

	for result.Next() {
		var participant entities.Participant

		err := result.Scan(&participant.Name, &participant.Avatar)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		participants = append(participants, participant)
	}

	return participants, nil
}

func (r *ParticipantRepository) JoinEvent(eventId int, userId int) error {
	// check join event in the past
	err := r.checkAlreadyJoinEvent(eventId, userId)

	if err != nil {
		log.Println(err)
		return err
	}

	stmt, err := r.db.Prepare("insert into participants(eventid, participantid, joined_at) values(?, ?, CURRENT_TIMESTAMP)")

	if err != nil {
		log.Println(err)
		return err
	}

	result, err := stmt.Exec(eventId, userId)

	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("failed joining event")
	}

	return nil
}

func (r *ParticipantRepository) UnjoinEvent(eventId int, userId int) error {
	stmt, err := r.db.Prepare("update participants set deleted_at = CURRENT_TIMESTAMP where deleted_at is NULL and eventid = ? and participantid = ?")

	if err != nil {
		log.Println(err)
		return err
	}

	res, err := stmt.Exec(eventId, userId)

	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("failed joining event")
	}

	return nil
}

// ini seharusnya muncul di event repository, tetapi daripada conflict ya sudah ditaruh di participant repository
func (r *ParticipantRepository) GetEventsByParticipantId(userId int) ([]entities.Event, error) {
	stmt, err := r.db.Prepare(`select e.id, e.name, e.host, u.name, e.description, e.datetime, e.location, e.category, e.photo
								from events e left join users u on e.hostid = u.id join participants p on e.id = p.eventid
								where e.deleted_at is NULL and p.participantid = ?`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := stmt.Query(userId)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer result.Close()

	var events []entities.Event

	for result.Next() {
		var event entities.Event

		err := result.Scan(&event.Id, &event.Name, &event.Host, &event.UserName, &event.Description, &event.Datetime, &event.Location, &event.Category, &event.Photo)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *ParticipantRepository) UnjoinAllEvent(userId int) error {
	stmt, err := r.db.Prepare("update participants set deleted_at = CURRENT_TIMESTAMP where deleted_at is NULL and participantid = ?")

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(userId)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *ParticipantRepository) DeleteAllParticipantByEventId(eventId int) error {
	stmt, err := r.db.Prepare("update participants set deleted_at = CURRENT_TIMESTAMP where deleted_at is NULL and eventid = ?")

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(eventId)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *ParticipantRepository) checkAlreadyJoinEvent(eventId int, userId int) error {
	stmt, err := r.db.Prepare("select count(id) from participants where deleted_at is NULL and eventid = ? and participantid = ?")

	if err != nil {
		log.Println(err)
		return err
	}

	res, err := stmt.Query(eventId, userId)

	if err != nil {
		log.Println(err)
		return err
	}

	defer res.Close()

	count := 0

	if res.Next() {
		if err = res.Scan(&count); err != nil {
			return err
		}
	}

	// Detect already join event
	if count != 0 {
		return errors.New("already join event")
	}

	return nil
}
