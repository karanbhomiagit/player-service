package player

import "github.com/karanbhomiagit/player-service/models"

// Usecase represents the player service's business logic as an interface
type Usecase interface {
	FetchPlayersByTeams(teamNames []string) ([]*models.Player, error)
}
