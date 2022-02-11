package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	_config "eventapp/config"
	_helpers "eventapp/delivery/helpers"
	_middlewares "eventapp/delivery/middlewares"
	_entities "eventapp/entities"
	_model "eventapp/entities/graph/model"
	_generated "eventapp/utils/graph/generated"
	"math"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input _model.NewUser) (*_model.CreateUserResponse, error) {
	// set default value
	id, phone, avatar := -1, "", ""
	_createdUser := _model.User{
		ID:     &id,
		Phone:  &phone,
		Avatar: &avatar,
	}

	// preprocessing input string
	name := strings.Title(strings.ToLower(strings.TrimSpace(input.Name)))
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	// check input string
	strings_to_check := []string{name, email, password}

	for _, s := range strings_to_check {
		// check empty string in mandatory input
		if s == "" {
			return &_model.CreateUserResponse{
				Code:    http.StatusBadRequest,
				Message: "mandatory input cannot be empty string",
				Data:    &_createdUser,
			}, nil
		}

		// check malicious character in input
		if err := _helpers.CheckStringInput(s); err != nil {
			return &_model.CreateUserResponse{
				Code:    http.StatusBadRequest,
				Message: s + ": " + err.Error(),
				Data:    &_createdUser,
			}, nil
		}
	}

	// check email pattern
	if err := _helpers.CheckEmailPattern(email); err != nil {
		return &_model.CreateUserResponse{
			Code:    http.StatusBadRequest,
			Message: email + ": " + err.Error(),
			Data:    &_createdUser,
		}, nil
	}

	// check password pattern
	if err := _helpers.CheckPasswordPattern(password); err != nil {
		return &_model.CreateUserResponse{
			Code:    http.StatusBadRequest,
			Message: password + ": " + err.Error(),
			Data:    &_createdUser,
		}, nil
	}

	// hashing password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	// detect failure in hashing password
	if err != nil {
		return &_model.CreateUserResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Data:    &_createdUser,
		}, nil
	}

	// prepare input to repository
	createUserData := _entities.User{
		Name:     name,
		Password: string(passwordHash),
		Email:    email,
	}

	// query via repository to create new user
	createdUser, code, err := r.userRepo.CreateUser(createUserData)

	// detect failure in repository
	if err != nil {
		return &_model.CreateUserResponse{
			Code:    code,
			Message: err.Error(),
			Data:    &_createdUser,
		}, nil
	}

	// prepare output to reponse
	id = createdUser.Id
	_createdUser.Name = createdUser.Name
	_createdUser.Email = createdUser.Email
	_createdUser.Password = "********"

	return &_model.CreateUserResponse{
		Code:    http.StatusOK,
		Message: "success create user",
		Data:    &_createdUser,
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id int, set _model.UpdateUser) (*_model.UpdateUserResponse, error) {
	// set default return value
	_id, phone, avatar := -1, "", ""
	_updatedUser := _model.User{
		ID:     &_id,
		Phone:  &phone,
		Avatar: &avatar,
	}

	// only registered user can update his/her own profile
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return &_model.UpdateUserResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
			Data:    &_updatedUser,
		}, nil
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	// detect unautorized update
	if id != convData.Id {
		return &_model.UpdateUserResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
			Data:    &_updatedUser,
		}, nil
	}

	// query via repository to get existing user profile
	updateUserData, err := r.userRepo.GetUserById(id)

	// detect failure in repository
	if err != nil {
		return &_model.UpdateUserResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    &_updatedUser,
		}, nil
	}

	// detect user unknown
	if updateUserData == (_entities.User{}) {
		return &_model.UpdateUserResponse{
			Code:    http.StatusBadRequest,
			Message: "user not found",
			Data:    &_updatedUser,
		}, nil
	}

	// detect change in user name
	if set.Name != nil && strings.TrimSpace(*set.Name) != "" {
		// preprocessing input string
		name := strings.Title(strings.ToLower(strings.TrimSpace(*set.Name)))

		// check malicious character in input
		if err := _helpers.CheckStringInput(name); err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusBadRequest,
				Message: name + ": " + err.Error(),
				Data:    &_updatedUser,
			}, nil
		}

		updateUserData.Name = name
	}

	// detect change in user email
	if set.Email != nil && strings.TrimSpace(*set.Email) != "" {
		// preprocessing input string
		email := strings.TrimSpace(*set.Email)

		// check malicious character in input
		if err := _helpers.CheckStringInput(email); err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusBadRequest,
				Message: email + ": " + err.Error(),
				Data:    &_updatedUser,
			}, nil
		}

		// check email pattern
		if err := _helpers.CheckEmailPattern(email); err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusBadRequest,
				Message: email + ": " + err.Error(),
				Data:    &_updatedUser,
			}, nil
		}

		updateUserData.Email = email
	}

	// detect change in user password
	if set.Password != nil && strings.TrimSpace(*set.Password) != "" {
		// preprocessing input string
		password := strings.TrimSpace(*set.Password)

		// check malicious character in input
		if err := _helpers.CheckStringInput(password); err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusBadRequest,
				Message: password + ": " + err.Error(),
				Data:    &_updatedUser,
			}, nil
		}

		// check password pattern
		if err := _helpers.CheckPasswordPattern(password); err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusBadRequest,
				Message: password + ": " + err.Error(),
				Data:    &_updatedUser,
			}, nil
		}

		// hashing password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

		// detect failure in hashing password
		if err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
				Data:    &_updatedUser,
			}, nil
		}

		updateUserData.Password = string(passwordHash)
	}

	// detect change in user phone
	if set.Phone != nil && strings.TrimSpace(*set.Phone) != "" {
		// preprocessing input string
		phone = strings.ReplaceAll(*set.Phone, " ", "")

		// check malicious character in input
		if err := _helpers.CheckStringInput(phone); err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusBadRequest,
				Message: phone + ": " + err.Error(),
				Data:    &_updatedUser,
			}, nil
		}

		// check phone pattern
		if err := _helpers.CheckPhonePattern(phone); err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusBadRequest,
				Message: phone + ": " + err.Error(),
				Data:    &_updatedUser,
			}, nil
		}

		updateUserData.Phone = phone
	}

	// detect change in user avatar
	if set.Avatar != nil && strings.TrimSpace(*set.Avatar) != "" {
		// preprocessing input string
		avatar := strings.TrimSpace(*set.Avatar)

		// check malicious character in input
		if err := _helpers.CheckStringInput(avatar); err != nil {
			return &_model.UpdateUserResponse{
				Code:    http.StatusBadRequest,
				Message: phone + ": " + err.Error(),
				Data:    &_updatedUser,
			}, nil
		}

		updateUserData.Avatar = avatar
	}

	updateUserData.Id = id

	// query via repository to update user
	updatedUser, code, err := r.userRepo.UpdateUser(updateUserData)

	// detect failure in repository
	if err != nil {
		return &_model.UpdateUserResponse{
			Code:    code,
			Message: err.Error(),
			Data:    &_updatedUser,
		}, nil
	}

	// prepare output to response
	_id, phone, avatar = updatedUser.Id, updatedUser.Phone, updatedUser.Avatar
	_updatedUser.Name = updatedUser.Name
	_updatedUser.Email = updatedUser.Email
	_updatedUser.Password = "********"

	return &_model.UpdateUserResponse{
		Code:    http.StatusOK,
		Message: "success update user",
		Data:    &_updatedUser,
	}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (*_model.DeleteUserResponse, error) {
	// only registered user can delete his/her own profile
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return &_model.DeleteUserResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, nil
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	// detect unautorized delete
	if id != convData.Id {
		return &_model.DeleteUserResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, nil
	}

	// // pertama, delete comment, jika ada
	// r.commentRepo.DeleteAllCommentByUserId(id)

	// // kedua, unjoin events, jika ada
	// r.participantRepo.UnjoinAllEvent(id)

	// // ketiga, delete events, jika ada
	// events, _ := r.eventRepo.GetEventByHostId(id)

	// for _, event := range events {
	// 	// untuk setiap event, delete comment, jika ada
	// 	r.commentRepo.DeleteAllCommentByEventId(event.Id)

	// 	// untuk setiap event, delete participants, jika ada
	// 	r.participantRepo.DeleteAllParticipantByEventId(event.Id)

	// 	r.eventRepo.DeleteEvent(event.Id)
	// }

	// terakhir, delete user
	code, err := r.userRepo.DeleteUser(id)

	// detect failure in repository
	if err != nil {
		return &_model.DeleteUserResponse{
			Code:    code,
			Message: err.Error(),
		}, nil
	}

	return &_model.DeleteUserResponse{
		Code:    http.StatusOK,
		Message: "success delete user",
	}, nil
}

func (r *mutationResolver) CreateEvent(ctx context.Context, input _model.NewEvent) (*_model.CreateEventResponse, error) {
	// set default return value
	id := -1
	_createdEvent := _model.Event{
		ID: &id,
	}

	// only registered user can create/host an event
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return &_model.CreateEventResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
			Data:    &_createdEvent,
		}, nil
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	// preprocessing input string
	name := strings.TrimSpace(input.Name)
	host := strings.TrimSpace(input.Host)
	datetime := strings.TrimSpace(input.Datetime)
	location := strings.TrimSpace(input.Location)
	category := strings.TrimSpace(input.Category)

	// check input string
	strings_to_check := []string{name, host, datetime, location, category}

	for _, s := range strings_to_check {
		// check empty string in mandatory input
		if s == "" {
			return &_model.CreateEventResponse{
				Code:    http.StatusBadRequest,
				Message: "mandatory input cannot be empty string",
				Data:    &_createdEvent,
			}, nil
		}

		// check malicious character in input
		if err := _helpers.CheckStringInput(s); err != nil {
			return &_model.CreateEventResponse{
				Code:    http.StatusBadRequest,
				Message: s + ": " + err.Error(),
				Data:    &_createdEvent,
			}, nil
		}
	}

	// check datetime pattern
	if err := _helpers.CheckDatetimePattern(datetime); err != nil {
		return &_model.CreateEventResponse{
			Code:    http.StatusBadRequest,
			Message: datetime + ": " + err.Error(),
			Data:    &_createdEvent,
		}, nil
	}

	// prepare input to repository
	createEventData := _entities.Event{
		Name:     name,
		Host:     host,
		Datetime: datetime,
		Location: location,
		Category: category,
		HostId:   convData.Id,
	}

	createdEvent, code, err := r.eventRepo.CreateEvent(createEventData)

	if err != nil {
		return &_model.CreateEventResponse{
			Code:    code,
			Message: err.Error(),
			Data:    &_createdEvent,
		}, nil
	}

	// prepare output to reponse
	id = createdEvent.Id
	_createdEvent.Name = createdEvent.Name
	_createdEvent.Username = createdEvent.UserName
	_createdEvent.Host = createdEvent.Host
	_createdEvent.Description = createdEvent.Description
	_createdEvent.Datetime = createdEvent.Datetime
	_createdEvent.Location = createdEvent.Location
	_createdEvent.Category = createdEvent.Category
	_createdEvent.Photo = createdEvent.Photo

	return &_model.CreateEventResponse{
		Code:    http.StatusOK,
		Message: "success create event",
		Data:    &_createdEvent,
	}, nil
}

func (r *mutationResolver) UpdateEvent(ctx context.Context, id int, set _model.UpdateEvent) (*_model.UpdateEventResponse, error) {
	// set default return value
	_id := -1
	_updatedEvent := _model.Event{
		ID: &_id,
	}

	// only registered user can update his/her own hosted event
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return &_model.UpdateEventResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
			Data:    &_updatedEvent,
		}, nil
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	// query via repository to get existing event detail
	updateEventData, err := r.eventRepo.GetEventByEventId(id)

	// detect failure in repository
	if err != nil {
		return &_model.UpdateEventResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Data:    &_updatedEvent,
		}, nil
	}

	// detect unknown event
	if updateEventData.HostId == 0 {
		return &_model.UpdateEventResponse{
			Code:    http.StatusBadRequest,
			Message: "event not found",
			Data:    &_updatedEvent,
		}, nil
	}

	// detect unauthorized update
	if updateEventData.HostId != convData.Id {
		return &_model.UpdateEventResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
			Data:    &_updatedEvent,
		}, nil
	}

	// detect change in event name
	if set.Name != nil && strings.TrimSpace(*set.Name) != "" {
		// preprocessing input string
		name := strings.TrimSpace(*set.Name)

		// check malicious character in input
		if err := _helpers.CheckStringInput(name); err != nil {
			return &_model.UpdateEventResponse{
				Code:    http.StatusBadRequest,
				Message: name + ": " + err.Error(),
				Data:    &_updatedEvent,
			}, nil
		}

		updateEventData.Name = name
	}

	// detect change in event host name
	if set.Host != nil && strings.TrimSpace(*set.Host) != "" {
		// preprocessing input string
		host := strings.TrimSpace(*set.Host)

		// check malicious character in input
		if err := _helpers.CheckStringInput(host); err != nil {
			return &_model.UpdateEventResponse{
				Code:    http.StatusBadRequest,
				Message: host + ": " + err.Error(),
				Data:    &_updatedEvent,
			}, nil
		}

		updateEventData.Host = host
	}

	// detect change in event category
	if set.Category != nil && strings.TrimSpace(*set.Category) != "" {
		// preprocessing input string
		category := strings.TrimSpace(*set.Category)

		// check malicious character in input
		if err := _helpers.CheckStringInput(category); err != nil {
			return &_model.UpdateEventResponse{
				Code:    http.StatusBadRequest,
				Message: category + ": " + err.Error(),
				Data:    &_updatedEvent,
			}, nil
		}

		updateEventData.Category = category
	}

	// detect change in event date and time
	if set.Datetime != nil && strings.TrimSpace(*set.Datetime) != "" {
		// preprocessing input string
		datetime := strings.TrimSpace(*set.Datetime)

		// check malicious character in input
		if err := _helpers.CheckStringInput(datetime); err != nil {
			return &_model.UpdateEventResponse{
				Code:    http.StatusBadRequest,
				Message: datetime + ": " + err.Error(),
				Data:    &_updatedEvent,
			}, nil
		}

		// check datetime pattern
		if err := _helpers.CheckDatetimePattern(datetime); err != nil {
			return &_model.UpdateEventResponse{
				Code:    http.StatusBadRequest,
				Message: datetime + ": " + err.Error(),
				Data:    &_updatedEvent,
			}, nil
		}

		updateEventData.Datetime = datetime
	} else {
		// in case no update in event date and time
		updateEventData.Datetime = strings.ReplaceAll(updateEventData.Datetime, "T", " ")
		updateEventData.Datetime = strings.ReplaceAll(updateEventData.Datetime, "Z", "")
	}

	// detect change in event location
	if set.Location != nil && strings.TrimSpace(*set.Location) != "" {
		// preprocessing input string
		location := strings.TrimSpace(*set.Location)

		// check malicious character in input
		if err := _helpers.CheckStringInput(location); err != nil {
			return &_model.UpdateEventResponse{
				Code:    http.StatusBadRequest,
				Message: location + ": " + err.Error(),
				Data:    &_updatedEvent,
			}, nil
		}

		updateEventData.Location = location
	}

	// detect change in event description
	if set.Description != nil && strings.TrimSpace(*set.Description) != "" {
		// preprocessing input string
		description := strings.TrimSpace(*set.Description)

		// check malicious character in input
		if err := _helpers.CheckStringInput(description); err != nil {
			return &_model.UpdateEventResponse{
				Code:    http.StatusBadRequest,
				Message: description + ": " + err.Error(),
				Data:    &_updatedEvent,
			}, nil
		}

		updateEventData.Description = description
	}

	// detect change in event photo
	if set.Photo != nil && strings.TrimSpace(*set.Photo) != "" {
		// preprocessing input string
		photo := strings.TrimSpace(*set.Photo)

		// check malicious character in input
		if err := _helpers.CheckStringInput(photo); err != nil {
			return &_model.UpdateEventResponse{
				Code:    http.StatusBadRequest,
				Message: photo + ": " + err.Error(),
				Data:    &_updatedEvent,
			}, nil
		}

		updateEventData.Photo = photo
	}

	updateEventData.Id = id

	// query via repository to update event
	updatedEvent, code, err := r.eventRepo.UpdateEvent(updateEventData)

	// detect failure in repository
	if err != nil {
		return &_model.UpdateEventResponse{
			Code:    code,
			Message: err.Error(),
			Data:    &_updatedEvent,
		}, nil
	}

	// prepare output to response
	_id = updatedEvent.Id
	_updatedEvent.Name = updatedEvent.Name
	_updatedEvent.Host = updatedEvent.Host
	_updatedEvent.Description = updatedEvent.Description
	_updatedEvent.Datetime = updatedEvent.Datetime
	_updatedEvent.Location = updatedEvent.Location
	_updatedEvent.Category = updatedEvent.Category
	_updatedEvent.Photo = updatedEvent.Photo

	return &_model.UpdateEventResponse{
		Code:    http.StatusOK,
		Message: "success update event",
		Data:    &_updatedEvent,
	}, nil
}

func (r *mutationResolver) DeleteEvent(ctx context.Context, id int) (*_model.DeleteEventResponse, error) {
	// only registered user can delete his/her own hosted event
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return &_model.DeleteEventResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, nil
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	// query via repository to get existing event detail
	deleteEventData, err := r.eventRepo.GetEventByEventId(id)

	// detect failure in repository
	if err != nil {
		return &_model.DeleteEventResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}, nil
	}

	// detect unknown event
	if deleteEventData.HostId == 0 {
		return &_model.DeleteEventResponse{
			Code:    http.StatusBadRequest,
			Message: "event not found",
		}, nil
	}

	// detect unauthorized delete
	if deleteEventData.HostId != convData.Id {
		return &_model.DeleteEventResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, nil
	}

	// // pertama, delete all comments, jika ada
	// r.commentRepo.DeleteAllCommentByEventId(id)

	// // kedua, delete all participants, jika ada
	// r.participantRepo.DeleteAllParticipantByEventId(id)

	// // terakhir, delete event
	code, err := r.eventRepo.DeleteEvent(id)

	// detect failure in repository
	if err != nil {
		return &_model.DeleteEventResponse{
			Code:    code,
			Message: err.Error(),
		}, nil
	}

	return &_model.DeleteEventResponse{
		Code:    http.StatusOK,
		Message: "success delete event",
	}, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, eventID int, input string) (*_model.Comment, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	comment, err := r.commentRepo.CreateComment(eventID, convData.Id, input)

	if err != nil {
		return nil, err
	}

	id := comment.Id

	responseMessage := _model.Comment{
		ID:        &id,
		UserID:    convData.Id,
		Name:      comment.UserName,
		Avatar:    comment.Avatar,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}

	return &responseMessage, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, commentID int) (*_model.SuccessResponse, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	if err := r.commentRepo.DeleteComment(commentID, convData.Id); err != nil {
		return nil, err
	}

	responseMessage := _model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Succes delete comment",
	}

	return &responseMessage, nil
}

func (r *mutationResolver) JoinEvent(ctx context.Context, eventID int) (*_model.SuccessResponse, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	// get event buat dapet tanggal
	// kondisi tanggal sekarang di bandingkan dengan format rfc
	event, err := r.eventRepo.GetEventByEventId(eventID)
	if err != nil {
		return nil, err
	}

	current_time := time.Now()
	eventdate, _ := time.Parse(time.RFC3339, event.Datetime)
	current_time.Before(eventdate)

	if !current_time.Before(eventdate) {
		return nil, errors.New(" Cant join, Event already past ")
	}

	if err := r.participantRepo.JoinEvent(eventID, convData.Id); err != nil {
		return nil, err
	}

	responseMessage := _model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Succes join event",
	}
	return &responseMessage, nil
}

func (r *mutationResolver) UnjoinEvent(ctx context.Context, eventID int) (*_model.SuccessResponse, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	if err := r.participantRepo.UnjoinEvent(eventID, convData.Id); err != nil {
		return nil, err
	}

	responseMessage := _model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Succes unjoin event",
	}
	return &responseMessage, nil
}

func (r *queryResolver) GetUsers(ctx context.Context) (*_model.GetUsersResponse, error) {
	// set default return value
	_allUsers := []*_model.User{}

	// query via repository to get all users
	allUsers, err := r.userRepo.GetAllUsers()

	// detect failure in repository
	if err != nil {
		return &_model.GetUsersResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    _allUsers,
		}, nil
	}

	// detect empty user directory
	if len(allUsers) == 0 {
		return &_model.GetUsersResponse{
			Code:    http.StatusBadRequest,
			Message: "users directory is empty",
			Data:    _allUsers,
		}, nil
	}

	for _, user := range allUsers {
		id, phone, avatar := user.Id, user.Phone, user.Avatar
		_allUsers = append(_allUsers,
			&_model.User{
				ID:       &id,
				Name:     user.Name,
				Email:    user.Email,
				Password: "********",
				Phone:    &phone,
				Avatar:   &avatar,
			})
	}

	return &_model.GetUsersResponse{
		Code:    http.StatusBadRequest,
		Message: "success get all users",
		Data:    _allUsers,
	}, nil
}

func (r *queryResolver) GetUserByID(ctx context.Context, id int) (*_model.GetUserResponse, error) {
	// set default value
	_id, phone, avatar := -1, "", ""
	_getUser := _model.User{
		ID:       &_id,
		Name:     "",
		Email:    "",
		Password: "",
		Phone:    &phone,
		Avatar:   &avatar,
	}

	// query via repository to get user by id
	getUser, err := r.userRepo.GetUserById(id)

	// detect failure in repository
	if err != nil {
		return &_model.GetUserResponse{
			Code:    http.StatusOK,
			Message: err.Error(),
			Data:    &_getUser,
		}, nil
	}

	// detect unknown user
	if getUser == (_entities.User{}) {
		return &_model.GetUserResponse{
			Code:    http.StatusOK,
			Message: "user not found",
			Data:    &_getUser,
		}, nil
	}

	// prepare output to response
	_id, phone, avatar = getUser.Id, getUser.Phone, getUser.Avatar

	// prepare output to response
	_getUser.Name = getUser.Name
	_getUser.Email = getUser.Email
	_getUser.Password = "********"

	return &_model.GetUserResponse{
		Code:    http.StatusOK,
		Message: "success get user",
		Data:    &_getUser,
	}, nil
}

func (r *queryResolver) AuthLogin(ctx context.Context, email string, password string) (*_model.AuthLoginResponse, error) {
	// set default value
	id, name, token := -1, "", ""
	_loginUser := _model.Login{
		ID:    &id,
		Name:  &name,
		Token: token,
	}

	// check malicious character in input
	strings_to_check := []string{email, password}

	for _, s := range strings_to_check {
		if err := _helpers.CheckStringInput(s); err != nil {
			return &_model.AuthLoginResponse{
				Code:    http.StatusBadRequest,
				Message: s + ": " + err.Error(),
				Data:    &_loginUser,
			}, nil
		}
	}

	// check email pattern
	if err := _helpers.CheckEmailPattern(email); err != nil {
		return &_model.AuthLoginResponse{
			Code:    http.StatusBadRequest,
			Message: email + ": " + err.Error(),
			Data:    &_loginUser,
		}, nil
	}

	// query via repository
	loginData, err := r.authRepo.Login(email)

	// detect failure in repository
	if err != nil {
		return &_model.AuthLoginResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    &_loginUser,
		}, nil
	}

	// detect unauthorized login (email unknown)
	if loginData == (_entities.User{}) {
		return &_model.AuthLoginResponse{
			Code:    http.StatusUnauthorized,
			Message: "email is unknown",
			Data:    &_loginUser,
		}, nil
	}

	// detect unauhorized login (password mismatch)
	if err = bcrypt.CompareHashAndPassword([]byte(loginData.Password), []byte(password)); err != nil {
		return &_model.AuthLoginResponse{
			Code:    http.StatusUnauthorized,
			Message: "password does not match",
			Data:    &_loginUser,
		}, nil
	}

	token, err = _middlewares.CreateToken(loginData.Id)

	// detect failure in creating token
	if err != nil {
		return &_model.AuthLoginResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    &_loginUser,
		}, nil
	}

	id = loginData.Id

	// prepare output to reponse
	_loginUser.ID = &id
	_loginUser.Name = &loginData.Name
	_loginUser.Token = token

	return &_model.AuthLoginResponse{
		Code:    http.StatusOK,
		Message: "success login",
		Data:    &_loginUser,
	}, nil
}

func (r *queryResolver) Events(ctx context.Context, page int) (*_model.EventResponse, error) {
	responseData, totalEvent, err := r.eventRepo.GetAllEvent(page)

	if err != nil {
		return nil, err
	}

	eventResponseData := []*_model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &_model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	limit := 5
	totalPage := int(math.Ceil(float64(totalEvent) / float64(limit)))

	eventResponse := _model.EventResponse{
		Event:     eventResponseData,
		TotalPage: totalPage,
	}

	return &eventResponse, nil
}

func (r *queryResolver) EventByHostID(ctx context.Context, userID int) ([]*_model.Event, error) {
	responseData, err := r.eventRepo.GetEventByHostId(userID)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*_model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &_model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByLocation(ctx context.Context, location string, page int) (*_model.EventResponse, error) {
	responseData, totalEvent, err := r.eventRepo.GetEventByLocation(location, page)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*_model.Event{}

	for _, v := range responseData {
		id := v.Id
		eventResponseData = append(eventResponseData, &_model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}
	limit := 5

	totalPage := int(math.Ceil(float64(totalEvent) / float64(limit)))

	eventResponse := _model.EventResponse{
		Event:     eventResponseData,
		TotalPage: totalPage,
	}

	return &eventResponse, nil
}

func (r *queryResolver) EventByKeyword(ctx context.Context, keyword string, page int) (*_model.EventResponse, error) {
	responseData, totalEvent, err := r.eventRepo.GetEventByKeyword(keyword, page)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*_model.Event{}

	for _, v := range responseData {
		id := v.Id
		eventResponseData = append(eventResponseData, &_model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	limit := 5
	totalPage := int(math.Ceil(float64(totalEvent) / float64(limit)))

	eventResponse := _model.EventResponse{
		Event:     eventResponseData,
		TotalPage: totalPage,
	}

	return &eventResponse, nil
}

func (r *queryResolver) EventByCategory(ctx context.Context, category string, page int) (*_model.EventResponse, error) {
	responseData, totalEvent, err := r.eventRepo.GetEventByCategory(category, page)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*_model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &_model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	limit := 5
	totalPage := int(math.Ceil(float64(totalEvent) / float64(limit)))

	eventResponse := _model.EventResponse{
		Event:     eventResponseData,
		TotalPage: totalPage,
	}

	return &eventResponse, nil
}

func (r *queryResolver) EventByParticipantID(ctx context.Context, userID int) ([]*_model.Event, error) {
	responseData, err := r.participantRepo.GetEventsByParticipantId(userID)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*_model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &_model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByID(ctx context.Context, id int) (*_model.Event, error) {
	responseData, _ := r.eventRepo.GetEventByEventId(id)

	if responseData == (_entities.Event{}) {
		return nil, errors.New("not found")
	}

	eventid := responseData.Id

	responseEventData := _model.Event{
		ID:          &eventid,
		Name:        responseData.Name,
		Username:    responseData.UserName,
		Host:        responseData.Host,
		Description: responseData.Description,
		Datetime:    responseData.Datetime,
		Location:    responseData.Location,
		Category:    responseData.Category,
		Photo:       responseData.Photo,
	}

	return &responseEventData, nil
}

func (r *queryResolver) Participants(ctx context.Context, eventID int) ([]*_model.Participant, error) {
	responseData, err := r.participantRepo.GetParticipantsByEventId(eventID)

	if err != nil {
		return nil, errors.New("not found")
	}

	participantResponseData := []*_model.Participant{}

	for _, v := range responseData {
		participantResponseData = append(participantResponseData, &_model.Participant{Name: v.Name, Avatar: v.Avatar})
	}

	return participantResponseData, nil
}

func (r *queryResolver) Comments(ctx context.Context, eventID int) ([]*_model.Comment, error) {
	responseData, err := r.commentRepo.GetCommentsByEventId(eventID)

	if err != nil {
		return nil, errors.New("not found")
	}

	commentResponseData := []*_model.Comment{}

	for _, v := range responseData {
		id := v.Id

		commentResponseData = append(commentResponseData, &_model.Comment{ID: &id, UserID: v.UserId, Name: v.UserName, Avatar: v.Avatar, Content: v.Content, CreatedAt: v.CreatedAt})
	}

	return commentResponseData, nil
}

// Mutation returns _generated.MutationResolver implementation.
func (r *Resolver) Mutation() _generated.MutationResolver { return &mutationResolver{r} }

// Query returns _generated.QueryResolver implementation.
func (r *Resolver) Query() _generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
