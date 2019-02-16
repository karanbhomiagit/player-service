package http

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/karanbhomiagit/player-service/player"
)

type HttpPlayerHandler struct {
	PUsecase player.Usecase
}

func NewHttpPlayerHandler(pu player.Usecase) {
	handler := &HttpPlayerHandler{
		PUsecase: pu,
	}

	http.HandleFunc("/players", handler.playersHandler)
}

//Move this to config file
var validTeams = []string{"Germany", "England", "France", "Spain", "Manchester United", "Arsenal", "Chelsea", "Barcelona", "Real Madrid", "Bayern Munich"}

func (h *HttpPlayerHandler) playersHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	//Check the http request method
	switch method {
	case http.MethodGet:
		fmt.Println("Request GET /players")
		//Make call to usecase layer to fetch the players by Teams
		res, err := h.PUsecase.FetchPlayersByTeams(validTeams)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}
		//Convert the response to expected string format
		b := ""
		for i, player := range res {
			num := strconv.Itoa(i + 1)
			b += num + ". " + player.Name + "; " + player.Age + "; " + strings.Join(player.Teams, ", ") + "\n"
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, b)
	default:
		//Return 405 in case of other methods
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Unsupported Request Method")
	}
}
