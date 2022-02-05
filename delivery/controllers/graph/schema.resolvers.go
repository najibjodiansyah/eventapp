package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"eventapp/delivery/middlewares"
	"eventapp/entities"
	"eventapp/entities/graph/model"
	"eventapp/utils/graph/generated"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	userData := model.User{
		Name: input.Name,
		// Password: input.Password,
		Password:    string(passwordHash),
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Avatar:      input.Avatar,
	}
	res, err := r.userRepo.Create(userData)
	if err != nil {
		return nil, errors.New("failed Create User")
	}
	responseMessage := model.User{
		Name:  res.Name,
		Email: res.Email,
		// Password: res.Password,
		PhoneNumber: res.PhoneNumber,
		Avatar:      res.Avatar,
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
	user, err := r.userRepo.GetbyId(id)
	if err != nil {
		return nil, errors.New("not found")
	}
	// user.Name = *set.Name //-----invalid memory address
	// user.Email = *set.Email
	// passwordHash, _ := bcrypt.GenerateFromPassword([]byte(*set.Password), bcrypt.MinCost)
	// user.Password = string(passwordHash)
	// user.Password = *set.Password

	if set.Name != nil {
		user.Name = *set.Name
	}
	fmt.Println("ga masuk")

	if set.Email != nil {
		user.Email = *set.Email
	}

	if set.Password != nil {
		// user.Password = *set.Password
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(*set.Password), bcrypt.MinCost)
		user.Password = string(passwordHash)
	}
	if set.PhoneNumber != nil {
		user.PhoneNumber = set.PhoneNumber
	}
	if set.Avatar != nil {
		user.Avatar = set.Avatar
	}
	fmt.Println(user)
	res, errr := r.userRepo.Update(id, user)
	if errr != nil {
		return nil, errors.New("fail create")
	}
	responseMessage := model.User{
		ID:    res.ID,
		Name:  res.Name,
		Email: res.Email,
		// Password: res.Password,
		PhoneNumber: res.PhoneNumber,
		Avatar:      res.Avatar,
	}
	return &responseMessage, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (*model.Message, error) {
	dataLogin := ctx.Value("EchoContextKey") // auth jwt
	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}
	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user", convData.Id)

	if id != convData.Id {
		return nil, errors.New("unauthorized")
	}

	err := r.userRepo.Delete(id)
	if err != nil {
		return nil, errors.New("failed Delete User")
	}
	responseMessage := model.Message{
		Message: "Succes Delete User",
	}
	return &responseMessage, nil
}

func (r *mutationResolver) CreateEvent(ctx context.Context, input model.NewEvent) (*model.Event, error) {
	// passwordHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	eventData := entities.Event{
		Name:        input.Name,
		Category:    input.Category,
		Host:        input.Host,
		Description: input.Description,
		Datetime:    input.Datetime,
		Location:    input.Location,
		Photo:       input.Photo,
	}
	res, err := r.eventRepo.CreateEvent(eventData)
	if err != nil {
		fmt.Println("cekk", err)
		return nil, errors.New("failed Create Event")
	}
	responseMessage := model.Event{
		Name:        res.Name,
		Category:    res.Category,
		Host:        res.Host,
		Description: res.Description,
		Datetime:    res.Datetime,
		Location:    res.Location,
		Photo:       res.Photo,
	}
	return &responseMessage, nil
}

func (r *mutationResolver) UpdateEvent(ctx context.Context, id int, set model.NewEvent) (*model.Event, error) {
	var event entities.Event
	event.Name = set.Name
	event.Category = set.Category
	event.Host = set.Host
	event.Description = set.Description
	event.Datetime = set.Datetime
	event.Location = set.Location
	event.Photo = set.Photo
	res, err := r.eventRepo.GetUpdateEvent(id, event)
	if err != nil {
		fmt.Println("cekk", err)
		return nil, errors.New("failed Create User")
	}
	responseMessage := model.Event{
		Name:        res.Name,
		Category:    res.Category,
		Host:        res.Host,
		Description: res.Description,
		Datetime:    res.Datetime,
		Location:    res.Location,
		Photo:       res.Photo,
	}
	return &responseMessage, nil
}

func (r *mutationResolver) DeleteEvent(ctx context.Context, id int) (*model.SuccessResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// bagus, hanya user yang login yang bisa create comment
func (r *mutationResolver) CreateComment(ctx context.Context, eventID int, input string) (*model.SuccessResponse, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	if err := r.commentRepo.CreateComment(eventID, convData.Id, input); err != nil {
		return nil, err
	}

	responseMessage := model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Succes create comment",
	}
	return &responseMessage, nil
}

// bagus, hanya user yang login yang bisa delete comment nya sendiri
func (r *mutationResolver) DeleteComment(ctx context.Context, eventID int) (*model.SuccessResponse, error) {
	dataLogin := ctx.Value("EchoContextKey")

	if dataLogin == nil {
		return nil, errors.New("unauthorized")
	}

	convData := ctx.Value("EchoContextKey").(*middlewares.User)
	fmt.Println("id user ", convData.Id)

	if err := r.commentRepo.DeleteComment(eventID, convData.Id); err != nil {
		return nil, err
	}

	responseMessage := model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Succes delete comment",
	}

	return &responseMessage, nil
}

// BAGUS, hanya user yang login yang bisa join event
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

// BAGUS, hanya user yang login yang bisa unjoin event
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
		convertID := int(*user.ID)
		userResponseData = append(userResponseData, &model.User{ID: &convertID, Name: user.Name, Email: user.Email, Password: user.Password, PhoneNumber: user.PhoneNumber, Avatar: user.Avatar})
	}

	return userResponseData, nil
}

func (r *queryResolver) UsersByID(ctx context.Context, id int) (*model.User, error) {
	responseData, err := r.userRepo.GetbyId(id)

	if err != nil {
		return nil, errors.New("not found")
	}

	responseUserData := model.User{}
	responseUserData.Email = responseData.Email
	responseUserData.ID = responseData.ID
	responseUserData.Name = responseData.Name
	responseUserData.Password = responseData.Password
	responseUserData.PhoneNumber = responseData.PhoneNumber
	responseUserData.Avatar = responseData.Avatar

	return &responseUserData, nil
}

func (r *queryResolver) AuthLogin(ctx context.Context, email string, password string) (*model.LoginResponse, error) {
	// password = hash
	user, err := r.authRepo.Login(email)
	if err != nil {
		return nil, err
	}

	errBcy := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if user.Email != email && errBcy == nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}

	authToken, err := middlewares.CreateToken((int(*user.ID)))
	if err != nil {
		return nil, fmt.Errorf("create token gagal")
	}

	response := model.LoginResponse{
		Message: "Success",
		ID:      *user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Token:   authToken,
	}

	return &response, nil
}

func (r *queryResolver) Events(ctx context.Context, page int) ([]*model.Event, error) {
	responseData, err := r.eventRepo.GetAllEvent(page)
	fmt.Println(responseData)

	if err != nil {
		return nil, err
	}

	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		var event model.Event
		event.ID = &v.ID
		event.Name = v.Name
		event.Category = v.Category
		event.Host = v.Host
		event.Description = v.Description
		event.Datetime = v.Datetime
		event.Location = v.Location
		event.Photo = v.Photo

		eventResponseData = append(eventResponseData, &event)
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByHostID(ctx context.Context, userID int) ([]*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) EventByLocation(ctx context.Context, location string, page int) ([]*model.Event, error) {
	responseData, err := r.eventRepo.GetEventByLocation(location, page)
	fmt.Println(responseData)

	if err != nil {
		return nil, errors.New("not found")
	}
	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		var event model.Event
		event.ID = &v.ID
		event.Name = v.Name
		event.Category = v.Category
		event.Datetime = v.Datetime
		event.Location = v.Location
		event.Description = v.Description
		event.Photo = v.Photo
		eventResponseData = append(eventResponseData, &event)
	}

	return eventResponseData, nil
}

func (r *queryResolver) EventByKeyword(ctx context.Context, keyword string, page int) ([]*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) EventByCategory(ctx context.Context, category string, page int) ([]*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

// BAGUS, semua orang bisa lihat daftar event by participant id
func (r *queryResolver) EventByParticipantID(ctx context.Context, userID int) ([]*model.Event, error) {
	responseData, err := r.participantRepo.GetEventsByParticipantId(userID)
	fmt.Println(responseData)

	if err != nil {
		return nil, errors.New("not found")
	}

	eventResponseData := []*model.Event{}

	for _, v := range responseData {
		var event model.Event
		event.ID = &v.ID
		event.Name = v.Name
		event.Username = v.UserName
		event.Category = v.Category
		event.Host = v.Host
		event.Description = v.Description
		event.Datetime = v.Datetime
		event.Location = v.Location
		event.Photo = v.Photo

		eventResponseData = append(eventResponseData, &event)
	}

	return eventResponseData, nil
}

// BAGUS, semua orang bisa lihat daftar peserta tanpa harus login
func (r *queryResolver) Participants(ctx context.Context, eventID int) ([]*model.Participant, error) {
	responseData, err := r.participantRepo.GetParticipantsByEventId(eventID)
	fmt.Println(responseData)

	if err != nil {
		return nil, errors.New("not found")
	}

	participantResponseData := []*model.Participant{}

	for _, v := range responseData {
		var participant model.Participant
		participant.Name = v.Name
		participant.Avatar = v.Avatar
		participantResponseData = append(participantResponseData, &participant)
	}

	return participantResponseData, nil
}

// BAGUS, semua orang bisa lihat comment tanpa harus login
func (r *queryResolver) Comments(ctx context.Context, eventID int) ([]*model.Comment, error) {
	responseData, err := r.commentRepo.GetCommentsByEventId(eventID)
	fmt.Println(responseData)

	if err != nil {
		return nil, errors.New("not found")
	}

	commentResponseData := []*model.Comment{}

	for _, v := range responseData {
		var comment model.Comment
		comment.Name = v.UserName
		comment.Avatar = v.Avatar
		comment.Content = v.Content
		comment.CreatedAt = v.CreatedAt
		commentResponseData = append(commentResponseData, &comment)
	}

	return commentResponseData, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
