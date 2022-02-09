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
	id, phoneNumber, avatar := -1, "", ""
	_createdUser := _model.User{
		ID:          &id,
		Name:        "",
		Email:       "",
		Password:    "",
		PhoneNumber: &phoneNumber,
		Avatar:      &avatar,
	}

	// check mallicious character in input
	strings_to_check := []string{input.Name, input.Email, input.Password}

	for _, s := range strings_to_check {
		if err := _helpers.CheckStringInput(s); err != nil {
			return &_model.CreateUserResponse{
				Code:    http.StatusBadRequest,
				Message: s + ": " + err.Error(),
				Data:    &_createdUser,
			}, nil
		}
	}

	// check email pattern
	if err := _helpers.CheckEmailPattern(input.Email); err != nil {
		return &_model.CreateUserResponse{
			Code:    http.StatusBadRequest,
			Message: input.Email + ": " + err.Error(),
			Data:    &_createdUser,
		}, nil
	}

	// hashing password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	// prepare input to repository
	userData := _entities.User{
		Name:     input.Name,
		Password: string(passwordHash),
		Email:    input.Email,
	}

	// input to repository
	createdUser, code, err := r.userRepo.Create(userData)

	// detect failure in repository
	if err != nil {
		return &_model.CreateUserResponse{
			Code:    code,
			Message: err.Error(),
			Data:    &_createdUser,
		}, nil
	}

	id = createdUser.Id
	createdUser.Password = "********"

	// prepare output to reponse
	_createdUser.Name = createdUser.Name
	_createdUser.Email = createdUser.Email
	_createdUser.Password = createdUser.Password
	_createdUser.PhoneNumber = &phoneNumber
	_createdUser.Avatar = &avatar

	return &_model.CreateUserResponse{
		Code:    code,
		Message: "success create user",
		Data:    &_createdUser,
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id int, set _model.UpdateUser) (*_model.User, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey) // auth jwt

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	if id != convData.Id {
		return nil, errors.New("unauthorized")
	}

	user, err := r.userRepo.GetById(id)

	if err != nil {
		return nil, err
	}

	if set.Name != nil {
		user.Name = *set.Name
	}

	if set.Email != nil {
		user.Email = *set.Email
	}

	if set.Password != nil {
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(*set.Password), bcrypt.MinCost)
		user.Password = string(passwordHash)
	}

	if set.PhoneNumber != nil {
		phone := *set.PhoneNumber

		if err := _helpers.CheckPhonePattern(phone); err != nil {
			return nil, err
		}

		user.PhoneNumber = phone
	}

	if set.Avatar != nil {
		user.Avatar = *set.Avatar
	}

	res, err := r.userRepo.Update(id, user)

	if err != nil {
		return nil, errors.New("failed update user")
	}

	responseMessage := _model.User{
		ID:          &id,
		Name:        res.Name,
		Email:       res.Email,
		PhoneNumber: &res.PhoneNumber,
		Avatar:      &res.Avatar,
	}

	return &responseMessage, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (*_model.SuccessResponse, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey) // auth jwt

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	if id != convData.Id {
		return nil, errors.New("unauthorized")
	}

	// pertama, delete comment, jika ada
	r.commentRepo.DeleteAllCommentByUserId(id)

	// kedua, unjoin events, jika ada
	r.participantRepo.UnjoinAllEvent(id)

	// ketiga, delete events, jika ada
	events, _ := r.eventRepo.GetEventByHostId(id)

	for _, event := range events {
		// untuk setiap event, delete comment, jika ada
		r.commentRepo.DeleteAllCommentByEventId(event.Id)

		// untuk setiap event, delete participants, jika ada
		r.participantRepo.DeleteAllParticipantByEventId(event.Id)

		r.eventRepo.DeleteEvent(event.Id)
	}

	// terakhir, delete user
	if err := r.userRepo.Delete(id); err != nil {
		return nil, errors.New("failed delete user")
	}

	responseMessage := _model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "succes delete user",
	}

	return &responseMessage, nil
}

func (r *mutationResolver) CreateEvent(ctx context.Context, input _model.NewEvent) (*_model.Event, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	// untuk saat ini, semua field dibuat required
	eventData := _entities.Event{
		Name:        input.Name,
		Category:    input.Category,
		Host:        input.Host,
		Description: input.Description,
		Datetime:    input.Datetime,
		Location:    input.Location,
		Photo:       input.Photo,
	}

	res, err := r.eventRepo.CreateEvent(convData.Id, eventData)

	if err != nil {
		return nil, errors.New("failed create event")
	}

	id := res.Id

	responseMessage := _model.Event{
		ID:          &id,
		Name:        res.Name,
		Username:    res.UserName,
		Host:        res.Host,
		Description: res.Description,
		Datetime:    res.Datetime,
		Location:    res.Location,
		Category:    res.Category,
		Photo:       res.Photo,
	}

	return &responseMessage, nil
}

func (r *mutationResolver) UpdateEvent(ctx context.Context, id int, set _model.UpdateEvent) (*_model.Event, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	event, _ := r.eventRepo.GetEventByEventId(id)

	if event.HostId == 0 {
		return nil, errors.New("not found")
	}

	if event.HostId != convData.Id {
		return nil, errors.New("unauthorized")
	}

	if set.Name != nil {
		event.Name = *set.Name
	}

	if set.Host != nil {
		event.Host = *set.Host
	}

	if set.Category != nil {
		event.Category = *set.Category
	}

	if set.Datetime != nil {
		event.Datetime = *set.Datetime
	}

	if set.Location != nil {
		event.Location = *set.Location
	}

	if set.Description != nil {
		event.Description = *set.Description
	}

	if set.Photo != nil {
		event.Photo = *set.Photo
	}

	event.Datetime = strings.ReplaceAll(event.Datetime, "T", " ")
	event.Datetime = strings.ReplaceAll(event.Datetime, "Z", "")

	event.Id = id
	res, err := r.eventRepo.UpdateEvent(event)

	if err != nil {
		return nil, errors.New("failed update event")
	}

	responseMessage := _model.Event{
		Name:        res.Name,
		Username:    res.UserName,
		Host:        res.Host,
		Description: res.Description,
		Datetime:    res.Datetime,
		Location:    res.Location,
		Category:    res.Category,
		Photo:       res.Photo,
	}

	return &responseMessage, nil
}

func (r *mutationResolver) DeleteEvent(ctx context.Context, id int) (*_model.SuccessResponse, error) {
	dataLogin := ctx.Value(_config.GetConfig().ContextKey)

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value(_config.GetConfig().ContextKey).(*_middlewares.User)

	event, _ := r.eventRepo.GetEventByEventId(id)

	if event.HostId == 0 {
		return nil, errors.New("not found")
	}

	if event.HostId != convData.Id {
		return nil, errors.New("unauthorized")
	}

	// pertama, delete all comments, jika ada
	r.commentRepo.DeleteAllCommentByEventId(id)

	// kedua, delete all participants, jika ada
	r.participantRepo.DeleteAllParticipantByEventId(id)

	// terakhir, delete event
	if err := r.eventRepo.DeleteEvent(id); err != nil {
		return nil, errors.New("failed delete event")
	}

	responseMessage := _model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "succes delete event",
	}

	return &responseMessage, nil
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

func (r *queryResolver) Users(ctx context.Context) ([]*_model.User, error) {
	responseData, err := r.userRepo.Get()

	if err != nil {
		return nil, errors.New("user not found")
	}

	userResponseData := []*_model.User{}

	for _, user := range responseData {
		id := user.Id
		phoneNumber := user.PhoneNumber
		avatar := user.Avatar
		userResponseData = append(userResponseData, &_model.User{ID: &id, Name: user.Name, Email: user.Email, Password: user.Password, PhoneNumber: &phoneNumber, Avatar: &avatar})
	}

	return userResponseData, nil
}

func (r *queryResolver) UserByID(ctx context.Context, id int) (*_model.User, error) {
	responseData, err := r.userRepo.GetById(id)

	if err != nil {
		return nil, errors.New("not found")
	}

	phoneNumber := responseData.PhoneNumber
	avatar := responseData.Avatar

	responseUserData := _model.User{
		ID:          &id,
		Name:        responseData.Name,
		Email:       responseData.Email,
		PhoneNumber: &phoneNumber,
		Avatar:      &avatar,
	}

	return &responseUserData, nil
}

func (r *queryResolver) AuthLogin(ctx context.Context, email string, password string) (*_model.LoginResponse, error) {
	user, err := r.authRepo.Login(email)

	if err != nil {
		return nil, err
	}

	if user == (_entities.User{}) {
		return nil, errors.New("email is wrong")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("password does not match")
	}

	authToken, err := _middlewares.CreateToken((user.Id))

	if err != nil {
		return nil, errors.New("failed create token")
	}

	response := _model.LoginResponse{
		Message: "Login success",
		ID:      user.Id,
		Name:    user.Name,
		Email:   user.Email,
		Token:   authToken,
	}

	return &response, nil
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
