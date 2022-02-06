package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"eventapp/delivery/helpers"
	"eventapp/delivery/middlewares"
	"eventapp/entities"
	"eventapp/entities/graph/model"
	"eventapp/utils/graph/generated"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	// check email pattern
	if err := helpers.CheckEmailPattern(input.Email); err != nil {
		return nil, err
	}

	// hashing password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	userData := entities.User{
		Name:     input.Name,
		Password: string(passwordHash),
		Email:    input.Email,
	}

	if input.PhoneNumber != nil {
		userData.PhoneNumber = *input.PhoneNumber
	}

	if input.Avatar != nil {
		userData.Avatar = *input.Avatar
	}

	res, err := r.userRepo.Create(userData)

	if err != nil {
		return nil, errors.New("failed create user")
	}

	id := res.Id
	phoneNumber := res.PhoneNumber
	avatar := res.Avatar

	responseMessage := model.User{
		ID:          &id,
		Name:        res.Name,
		Email:       res.Email,
		PhoneNumber: &phoneNumber,
		Avatar:      &avatar,
	}

	return &responseMessage, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id int, set model.UpdateUser) (*model.User, error) {
	dataLogin := ctx.Value("EchoContextKey") // auth jwt

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user", convData.Id)

	if id != convData.Id {
		return nil, errors.New("unauthorized")
	}

	user, err := r.userRepo.GetById(id)

	if err != nil {
		return nil, errors.New("not found")
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
		user.PhoneNumber = *set.PhoneNumber
	}

	if set.Avatar != nil {
		user.Avatar = *set.Avatar
	}

	res, errr := r.userRepo.Update(id, user)

	if errr != nil {
		return nil, errors.New("failed update user")
	}

	responseMessage := model.User{
		ID:          &id,
		Name:        res.Name,
		Email:       res.Email,
		PhoneNumber: &res.PhoneNumber,
		Avatar:      &res.Avatar,
	}

	return &responseMessage, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (*model.SuccessResponse, error) {
	dataLogin := ctx.Value("EchoContextKey") // auth jwt

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user", convData.Id)

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

	responseMessage := model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "succes delete user",
	}

	return &responseMessage, nil
}

func (r *mutationResolver) CreateEvent(ctx context.Context, input model.NewEvent) (*model.Event, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	// untuk saat ini, semua field dibuat required
	eventData := entities.Event{
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

	responseMessage := model.Event{
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

func (r *mutationResolver) UpdateEvent(ctx context.Context, id int, set model.UpdateEvent) (*model.Event, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	event, _ := r.eventRepo.GetById(id)

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

	res, err := r.eventRepo.UpdateEvent(id, event)

	if err != nil {
		return nil, errors.New("failed update event")
	}

	responseMessage := model.Event{
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

func (r *mutationResolver) DeleteEvent(ctx context.Context, id int) (*model.SuccessResponse, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	event, _ := r.eventRepo.GetById(id)

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

	responseMessage := model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "succes delete event",
	}

	return &responseMessage, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, eventID int, input string) (*model.Comment, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	comment, err := r.commentRepo.CreateComment(eventID, convData.Id, input)

	if err != nil {
		return nil, err
	}

	id := comment.Id

	responseMessage := model.Comment{
		ID:        &id,
		UserID:    convData.Id,
		Name:      comment.UserName,
		Avatar:    comment.Avatar,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}

	return &responseMessage, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, commentID int) (*model.SuccessResponse, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	if err := r.commentRepo.DeleteComment(commentID, convData.Id); err != nil {
		return nil, err
	}

	responseMessage := model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Succes delete comment",
	}

	return &responseMessage, nil
}

func (r *mutationResolver) JoinEvent(ctx context.Context, eventID int) (*model.SuccessResponse, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	if err := r.participantRepo.JoinEvent(eventID, convData.Id); err != nil {
		return nil, err
	}

	responseMessage := model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Succes join event",
	}
	return &responseMessage, nil
}

func (r *mutationResolver) UnjoinEvent(ctx context.Context, eventID int) (*model.SuccessResponse, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	if err := r.participantRepo.UnjoinEvent(eventID, convData.Id); err != nil {
		return nil, err
	}

	responseMessage := model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Succes unjoin event",
	}
	return &responseMessage, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	responseData, err := r.userRepo.Get()

	if err != nil {
		return nil, errors.New("user not found")
	}

	userResponseData := []*model.User{}

	for _, user := range responseData {
		id := user.Id
		phoneNumber := user.PhoneNumber
		avatar := user.Avatar
		userResponseData = append(userResponseData, &model.User{ID: &id, Name: user.Name, Email: user.Email, Password: user.Password, PhoneNumber: &phoneNumber, Avatar: &avatar})
	}

	return userResponseData, nil
}

func (r *queryResolver) UsersByID(ctx context.Context, id int) (*model.User, error) {
	responseData, err := r.userRepo.GetById(id)

	if err != nil {
		return nil, errors.New("not found")
	}

	phoneNumber := responseData.PhoneNumber
	avatar := responseData.Avatar

	responseUserData := model.User{
		ID:          &id,
		Name:        responseData.Name,
		Email:       responseData.Email,
		PhoneNumber: &phoneNumber,
		Avatar:      &avatar,
	}

	return &responseUserData, nil
}

func (r *queryResolver) AuthLogin(ctx context.Context, email string, password string) (*model.LoginResponse, error) {
	user, err := r.authRepo.Login(email)

	if err != nil {
		return nil, err
	}

	if user == (entities.User{}) {
		return nil, errors.New("email is wrong")
	}

	match := helpers.CheckPasswordHash(password, user.Password)

	if !match {
		return nil, errors.New("password does not match")
	}

	authToken, err := middlewares.CreateToken((user.Id))

	if err != nil {
		return nil, errors.New("failed create token")
	}

	response := model.LoginResponse{
		Message: "Login success",
		ID:      user.Id,
		Name:    user.Name,
		Email:   user.Email,
		Token:   authToken,
	}

	return &response, nil
}

func (r *queryResolver) Events(ctx context.Context, page int) ([]*model.Event, error) {
	responseData, err := r.eventRepo.GetAllEvent(page)

	if err != nil {
		return nil, err
	}

	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByHostID(ctx context.Context, userID int) ([]*model.Event, error) {
	responseData, err := r.eventRepo.GetEventByHostId(userID)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByLocation(ctx context.Context, location string, page int) ([]*model.Event, error) {
	responseData, err := r.eventRepo.GetEventByLocation(location, page)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByKeyword(ctx context.Context, keyword string, page int) ([]*model.Event, error) {
	responseData, err := r.eventRepo.GetEventByKeyword(keyword, page)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByCategory(ctx context.Context, category string, page int) ([]*model.Event, error) {
	responseData, err := r.eventRepo.GetEventByCategory(category, page)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByParticipantID(ctx context.Context, userID int) ([]*model.Event, error) {
	responseData, err := r.participantRepo.GetEventsByParticipantId(userID)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		id := v.Id

		eventResponseData = append(eventResponseData, &model.Event{ID: &id, Name: v.Name, Host: v.Host, Category: v.Category, Datetime: v.Datetime, Location: v.Location, Description: v.Description, Photo: v.Photo, Username: v.UserName})
	}

	return eventResponseData, nil
}

func (r *queryResolver) Participants(ctx context.Context, eventID int) ([]*model.Participant, error) {
	responseData, err := r.participantRepo.GetParticipantsByEventId(eventID)

	if err != nil {
		return nil, errors.New("not found")
	}

	participantResponseData := []*model.Participant{}

	for _, v := range responseData {
		participantResponseData = append(participantResponseData, &model.Participant{Name: v.Name, Avatar: v.Avatar})
	}

	return participantResponseData, nil
}

func (r *queryResolver) Comments(ctx context.Context, eventID int) ([]*model.Comment, error) {
	responseData, err := r.commentRepo.GetCommentsByEventId(eventID)

	if err != nil {
		return nil, errors.New("not found")
	}

	commentResponseData := []*model.Comment{}

	for _, v := range responseData {
		id := v.Id

		commentResponseData = append(commentResponseData, &model.Comment{ID: &id, UserID: v.UserId, Name: v.UserName, Avatar: v.Avatar, Content: v.Content, CreatedAt: v.CreatedAt})
	}

	return commentResponseData, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
