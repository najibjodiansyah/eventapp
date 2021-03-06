package main

import (
	_config "eventapp/config"
	_graph "eventapp/delivery/controllers/graph"
	_router "eventapp/delivery/routers"
	_authRepo "eventapp/repository/auth"
	_commentRepo "eventapp/repository/comment"
	_eventRepo "eventapp/repository/event"
	_participantRepo "eventapp/repository/participant"
	_userRepo "eventapp/repository/user"

	_util "eventapp/utils"

	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	//load config if available or set to default
	config := _config.GetConfig()

	//initialize database connection based on given config
	db := _util.MysqlDriver(config)

	//initiate user model
	userRepo := _userRepo.New(db)
	authRepo := _authRepo.New(db)
	eventRepo := _eventRepo.New(db)
	commentRepo := _commentRepo.New(db)
	participantRepo := _participantRepo.New(db)

	//create echo http
	e := echo.New()
	client := _graph.NewResolver(
		authRepo,
		commentRepo,
		eventRepo,
		participantRepo,
		userRepo,
	)
	srv := _router.NewGraphQLServer(client)

	//register API path and controller
	_router.RegisterPath(e, srv)

	// run server
	address := fmt.Sprintf(":%d", config.Port)

	if err := e.Start(address); err != nil {
		log.Println(err)
		log.Println("shutting down the server")
	}
}
