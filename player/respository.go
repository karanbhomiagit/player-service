package player

import "github.com/karanbhomiagit/player-service/models"

// Repository represents the storage/retrieval as an interface
type Repository interface {
	FetchByID(id int) (*models.Team, error)
}
