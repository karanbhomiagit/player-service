package usecase

import (
	"errors"
	"os"
	"testing"

	"github.com/karanbhomiagit/player-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedPlayerRepository struct {
	mock.Mock
}

func (pr *MockedPlayerRepository) FetchByID(id int) (*models.Team, error) {
	args := pr.Called(id)
	return args.Get(0).(*models.Team), args.Error(1)
}

/*
	Actual test functions
*/

func TestFetchPlayersByTeams(t *testing.T) {

	t.Run("Successfully fetch sorted players if repository returns teams with no players", func(t *testing.T) {
		testObj := new(MockedPlayerRepository)
		testTeam := models.Team{
			ID:      123,
			Name:    "A",
			Players: []models.Player{},
		}
		testObj.On("FetchByID", 1).Return(&testTeam, nil)

		testTeam1 := models.Team{
			ID:      124,
			Name:    "B",
			Players: []models.Player{},
		}
		testObj.On("FetchByID", 2).Return(&testTeam1, nil)

		testTeam2 := models.Team{
			ID:      125,
			Name:    "C",
			Players: []models.Player{},
		}
		testObj.On("FetchByID", 3).Return(&testTeam2, nil)

		playerUsecase := NewPlayerUsecase(testObj)
		validTeams := []string{"A", "B"}
		os.Setenv("MAX_CONCURRENT_REQUESTS", "3")
		response, err := playerUsecase.FetchPlayersByTeams(validTeams)
		assert := assert.New(t)
		assert.Nil(err)
		assert.Equal([]*models.Player{}, response)
		testObj.AssertExpectations(t)
	})

	t.Run("Successfully fetch sorted players if repository returns teams with unique players", func(t *testing.T) {
		testObj := new(MockedPlayerRepository)
		testTeam := models.Team{
			ID:   123,
			Name: "A",
			Players: []models.Player{
				{"25207", "Dead Pool", "24", []string{}},
				{"25208", "Charles Xavier", "58", []string{}},
			},
		}
		testObj.On("FetchByID", 1).Return(&testTeam, nil)

		testTeam1 := models.Team{
			ID:   124,
			Name: "B",
			Players: []models.Player{
				{"25209", "Bruce wayne", "24", []string{}},
				{"25210", "Ant Man", "38", []string{}},
			},
		}
		testObj.On("FetchByID", 2).Return(&testTeam1, nil)

		testTeam2 := models.Team{
			ID:      125,
			Name:    "C",
			Players: []models.Player{},
		}
		testObj.On("FetchByID", 3).Return(&testTeam2, nil)

		playerUsecase := NewPlayerUsecase(testObj)
		validTeams := []string{"A", "B"}
		os.Setenv("MAX_CONCURRENT_REQUESTS", "3")
		response, err := playerUsecase.FetchPlayersByTeams(validTeams)
		assert := assert.New(t)
		assert.Nil(err)
		expectedResponse := []models.Player{
			{"25210", "Ant Man", "38", []string{"B"}},
			{"25209", "Bruce wayne", "24", []string{"B"}},
			{"25208", "Charles Xavier", "58", []string{"A"}},
			{"25207", "Dead Pool", "24", []string{"A"}},
		}
		assert.Equal(expectedResponse[0], *(response[0]))
		assert.Equal(expectedResponse[1], *(response[1]))
		assert.Equal(expectedResponse[2], *(response[2]))
		assert.Equal(expectedResponse[3], *(response[3]))
		testObj.AssertExpectations(t)
	})

	t.Run("Successfully fetch sorted players if repository returns teams with common players", func(t *testing.T) {
		testObj := new(MockedPlayerRepository)
		testTeam := models.Team{
			ID:   123,
			Name: "A",
			Players: []models.Player{
				{"25207", "Dead Pool", "24", []string{}},
				{"25208", "Charles Xavier", "58", []string{}},
			},
		}
		testObj.On("FetchByID", 1).Return(&testTeam, nil)

		testTeam1 := models.Team{
			ID:   124,
			Name: "B",
			Players: []models.Player{
				{"25207", "Dead Pool", "24", []string{}},
				{"25210", "Ant Man", "38", []string{}},
			},
		}
		testObj.On("FetchByID", 2).Return(&testTeam1, nil)

		testTeam2 := models.Team{
			ID:      125,
			Name:    "C",
			Players: []models.Player{},
		}
		testObj.On("FetchByID", 3).Return(&testTeam2, nil)

		playerUsecase := NewPlayerUsecase(testObj)
		validTeams := []string{"A", "B"}
		os.Setenv("MAX_CONCURRENT_REQUESTS", "3")
		response, err := playerUsecase.FetchPlayersByTeams(validTeams)
		assert := assert.New(t)
		assert.Nil(err)
		expectedResponse := []models.Player{
			{"25210", "Ant Man", "38", []string{"B"}},
			{"25208", "Charles Xavier", "58", []string{"A"}},
			{"25207", "Dead Pool", "24", []string{"B", "A"}},
		}
		assert.Equal(expectedResponse[0], *(response[0]))
		assert.Equal(expectedResponse[1], *(response[1]))
		assert.Equal(expectedResponse[2].Name, (*response[2]).Name)
		assert.Equal(2, len((*response[2]).Teams))
		testObj.AssertExpectations(t)
	})

	t.Run("Successfully fetch sorted players even if repository returns error for 1 call", func(t *testing.T) {
		testObj := new(MockedPlayerRepository)
		testTeam := models.Team{
			ID:   123,
			Name: "A",
			Players: []models.Player{
				{"25207", "Dead Pool", "24", []string{}},
				{"25208", "Charles Xavier", "58", []string{}},
			},
		}
		testObj.On("FetchByID", 1).Return(&testTeam, nil)

		testTeam1 := models.Team{
			ID:   124,
			Name: "B",
			Players: []models.Player{
				{"25207", "Dead Pool", "24", []string{}},
				{"25210", "Ant Man", "38", []string{}},
			},
		}
		testObj.On("FetchByID", 2).Return(&testTeam1, nil)

		testObj.On("FetchByID", 3).Return(&models.Team{}, errors.New("Not found"))

		playerUsecase := NewPlayerUsecase(testObj)
		validTeams := []string{"A", "B"}
		os.Setenv("MAX_CONCURRENT_REQUESTS", "3")
		response, err := playerUsecase.FetchPlayersByTeams(validTeams)
		assert := assert.New(t)
		assert.Nil(err)
		expectedResponse := []models.Player{
			{"25210", "Ant Man", "38", []string{"B"}},
			{"25208", "Charles Xavier", "58", []string{"A"}},
			{"25207", "Dead Pool", "24", []string{"B", "A"}},
		}
		assert.Equal(expectedResponse[0], *(response[0]))
		assert.Equal(expectedResponse[1], *(response[1]))
		assert.Equal(expectedResponse[2].Name, (*response[2]).Name)
		assert.Equal(2, len((*response[2]).Teams))
		testObj.AssertExpectations(t)
	})

	t.Run("Successfully fetch sorted players with concurrency 1", func(t *testing.T) {
		testObj := new(MockedPlayerRepository)
		testTeam := models.Team{
			ID:      123,
			Name:    "A",
			Players: []models.Player{},
		}
		testObj.On("FetchByID", 1).Return(&testTeam, nil)

		testTeam1 := models.Team{
			ID:      124,
			Name:    "B",
			Players: []models.Player{},
		}
		testObj.On("FetchByID", 2).Return(&testTeam1, nil)

		testTeam2 := models.Team{
			ID:      125,
			Name:    "C",
			Players: []models.Player{},
		}
		testObj.On("FetchByID", 3).Return(&testTeam2, nil)

		playerUsecase := NewPlayerUsecase(testObj)
		validTeams := []string{"A", "C"}
		os.Setenv("MAX_CONCURRENT_REQUESTS", "1")
		response, err := playerUsecase.FetchPlayersByTeams(validTeams)
		assert := assert.New(t)
		assert.Nil(err)
		assert.Equal([]*models.Player{}, response)
		testObj.AssertExpectations(t)
	})
}
