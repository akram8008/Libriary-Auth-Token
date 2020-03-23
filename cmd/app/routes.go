package app

import (
	"github.com/akram8008/mux/pkg/mux"
	"github.com/akram8008/mux/pkg/mux/middleware/logger"

)

func (receiver *server) InitRoutes() {
	mux := receiver.router.(*mux.ExactMux)
	mux.POST(
		"/newUser",
		receiver.handleNewUser(),
		logger.Logger("Registration"),
	)

	mux.POST(
		"/login",
		receiver.handleLogin(),
		logger.Logger("Autorization"),
	)
}
