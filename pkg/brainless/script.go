package brainless

import (
	"bytes"
	"text/template"
)

func getScript(puzzle, host string) string {
	t, err := template.New("script").Parse(script)
	if err != nil {
		panic(err)
		return ""
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, struct {
		Puzzle string
		Host   string
	}{
		Puzzle: puzzle,
		Host:   host,
	})
	if err != nil {
		panic(err)
	}
	script := buf.String()
	return script
}

const script = `// copy and run this script on {{ .Puzzle }}

function solve(solver, stay) {
    document.getElementById('robot').value = 1
    var xhr = new XMLHttpRequest()
    xhr.open("POST", solver, true)
    xhr.setRequestHeader("Content-Type", "application/json")
    xhr.onreadystatechange = () => {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            resp = JSON.parse(xhr.responseText)
            update(resp)
            submit(stay)
        }
    }
    xhr.send(JSON.stringify(Game.task))
}

function update(response) {
    for (var key in resp) {
        if (resp.hasOwnProperty(key)) {
            Game.currentState[key] = resp[key]
        }
    }
}

function submit(stay) {
    if (!stay) {
        setTimeout(() => {
            document.getElementById('btnReady').click()
        }, 0)
    }
    Game.drawCurrentState()
}

solve("{{ .Host }}")
`
