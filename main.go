package main

import (
	"log"
	"net/http"
	"os"

	httpDeliver "github.com/karanbhomiagit/player-service/player/delivery/http"
	playerRepo "github.com/karanbhomiagit/player-service/player/repository"
	playerUsecase "github.com/karanbhomiagit/player-service/player/usecase"
)

func main() {
	//Initializing the repository
	pr := playerRepo.NewHttpTeamRepository()
	//Initializing the usecase
	pu := playerUsecase.NewPlayerUsecase(pr)

	//Initializing the delivery
	httpDeliver.NewHttpPlayerHandler(pu)
	log.Fatal(http.ListenAndServe(port(), nil))
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}
