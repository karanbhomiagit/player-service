package usecase

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/karanbhomiagit/player-service/models"
	"github.com/karanbhomiagit/player-service/player"
)

type playerUsecase struct {
	playerRepo player.Repository
}

func NewPlayerUsecase(p player.Repository) player.Usecase {
	return &playerUsecase{
		playerRepo: p,
	}
}

//FetchPlayersByTeams returns a list of players for the given teams
func (p *playerUsecase) FetchPlayersByTeams(teamNames []string) ([]*models.Player, error) {
	playerCollection := &models.PlayerCollection{
		TeamCount:     0,
		PlayersLookup: map[string]models.Player{},
		Mutex:         &sync.Mutex{},
	}
	validTeamCount := len(teamNames)
	//Assign maxConcurrency to determine the number of concurrent requests to repository layer
	maxConcurrencyEnv := os.Getenv("MAX_CONCURRENT_REQUESTS")
	maxConcurrency := 100
	if len(maxConcurrencyEnv) != 0 {
		maxConcurrency, _ = strconv.Atoi(maxConcurrencyEnv)
	}
	teamCounter := 1
	//Convert array of valid teams to map. This is because the source of teams contains duplicate names for different IDs.
	validTeamLookup := map[string]int{}
	for _, team := range teamNames {
		validTeamLookup[team] = 0
	}

	//fmt.Println(validTeamCount, maxConcurrency, validTeamLookup)
	//Run for loop till we have details of all valid teams
	for validTeamCount > playerCollection.TeamCount {
		var wg sync.WaitGroup
		//Start go routines to fetch team details by ID
		for i := 0; i < maxConcurrency; i++ {
			wg.Add(1)
			go func(teamNumber int) {
				defer wg.Done()
				res, err := p.playerRepo.FetchByID(teamNumber)
				if err != nil {
					fmt.Println("Response Error : ", err, teamNumber)
				} else if isValidTeam(res.Name, &validTeamLookup) {
					//Add response to playerCollection
					fmt.Println("Processing team ", res.Name, res.ID)
					playerCollection = p.addTeamInfo(playerCollection, res)
					playerCollection.TeamCount++
				}
			}(teamCounter + i)
		}
		wg.Wait()
		teamCounter += maxConcurrency
	}
	//Retrieve players slice from playerCollection and Sort alphabetically by name
	players := sortPlayersByName(playerCollection.PlayersLookup)
	return players, nil
}

func (p *playerUsecase) addTeamInfo(playerCollection *models.PlayerCollection, resp *models.Team) *models.PlayerCollection {
	players, teamName := resp.Players, resp.Name
	for _, player := range players {
		// sync.mutex to avoid concurrent access to PlayersLookup
		playerCollection.Mutex.Lock()
		//Check if the player exists in the playerCollection
		if val, ok := playerCollection.PlayersLookup[player.ID]; !ok {
			//If no, add an entry for the same
			player.Teams = []string{teamName}
			playerCollection.PlayersLookup[player.ID] = player
		} else {
			//If the player exists, append the team name to the corresponding entry
			val.Teams = append(val.Teams, teamName)
			playerCollection.PlayersLookup[val.ID] = val
		}
		playerCollection.Mutex.Unlock()
	}
	return playerCollection
}

func sortPlayersByName(playersMap map[string]models.Player) []*models.Player {
	players := make([]*models.Player, 0)
	for _, p := range playersMap {
		player := &models.Player{
			p.ID, p.Name, p.Age, p.Teams,
		}
		players = append(players, player)
	}
	sort.Slice(players, func(i, j int) bool {
		return (*players[i]).Name < (*players[j]).Name
	})
	return players
}

func isValidTeam(name string, validNamesLookup *map[string]int) bool {
	if val, ok := (*validNamesLookup)[name]; ok {
		if val == 0 {
			(*validNamesLookup)[name] = 1
			return true
		}
		return false

	}
	return false
}
