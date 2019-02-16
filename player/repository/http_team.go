package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/karanbhomiagit/player-service/models"
	"github.com/karanbhomiagit/player-service/player"
)

type httpTeamRepository struct {
}

func NewHttpTeamRepository() player.Repository {
	return &httpTeamRepository{}
}

//FetchByID contacts oneFootball via http call to get the team for a particular id
func (*httpTeamRepository) FetchByID(id int) (*models.Team, error) {
	idStr := strconv.Itoa(id)
	if idStr == "" {
		return nil, errors.New("ID should not be null")
	}

	url := fmt.Sprintf("%s%s%s", getBaseURLForTeams(), idStr, ".json")

	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error : ", err)
		return nil, err
	}

	//If the response status code is not 200, return appropriate error
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("Not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to fetch record")
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error : ", err)
		return nil, err
	}
	r, err := teamFromJSON(body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func getBaseURLForTeams() string {
	//Move this string to config file
	return "https://vintagemonster.onefootball.com/api/teams/en/"
}

func teamFromJSON(data []byte) (*models.Team, error) {
	t := models.TeamAPIResponse{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return t.Data.Team, nil
}
