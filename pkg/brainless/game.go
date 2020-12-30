package brainless

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Game struct {
	config    Config
	responder *responder
	Brain     Brain
	response  interface{}
}

func NewGame(w http.ResponseWriter, r *http.Request, config Config) *Game {
	g := &Game{
		config:    config,
		responder: newResponder(w, r, config.PuzzleURL),
	}
	return g
}

func (g *Game) HandleRequest() {
	g.handleRequest()
}

func (g *Game) getTask() {
	err := g.responder.r.ParseForm()
	if err != nil {
		g.respondError(err)
		return
	}
	var data [][]interface{}
	err = json.NewDecoder(g.responder.r.Body).Decode(&data)
	if err != nil {
		g.respondError(err)
		return
	}
	task := make([][]int, len(data))
	for i, row := range data {
		task[i] = make([]int, len(row))
		for j, v := range row {
			switch cell := v.(type) {
			case float64:
				task[i][j] = int(cell)
			case string:
				task[i][j], err = strconv.Atoi(cell)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	g.Brain.GetTask(task)
}

func (g *Game) ConnectBrain(b Brain) {
	g.Brain = b
}

func (g *Game) handleRequest() {
	switch g.responder.r.Method {
	case "OPTIONS":
		g.responder.PreFlight()
	case "GET":
		g.responder.RespondJS(getScript(g.config.PuzzleURL, g.config.FunctionURL))
	case "POST":
		g.solveAndRespond()
	default:
		err := g.responder.Respond(http.StatusMethodNotAllowed, nil)
		if err != nil {
			g.respondError(err)
		}
	}
}

func (g *Game) solve() {
	g.getTask()
	g.Brain.Setup()
	for i := 0; i < 1000 && !g.Brain.CheckDone(); i++ {
		g.Brain.Step()
	}
	g.response = g.Brain.ToResponse()
}

func (g *Game) respond(v interface{}) {
	g.responder.RespondJSON(v)
}

func (g *Game) solveAndRespond() {
	g.solve()
	g.respond(g.response)
}

func (g *Game) respondError(err error) {
	g.responder.RespondError(err)
}
