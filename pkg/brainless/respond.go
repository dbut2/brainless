package brainless

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responder struct {
	url string
	w   http.ResponseWriter
	r   *http.Request
}

func newResponder(w http.ResponseWriter, r *http.Request, url string) *responder {
	return &responder{
		w:   w,
		r:   r,
		url: url,
	}
}

func (r *responder) PreFlight() {
	headers := r.w.Header()
	headers.Add("Access-Control-Allow-Origin", r.url)
	headers.Add("Access-Control-Allow-Headers", "Content-Type")
	headers.Add("Access-Control-Allow-Methods", "POST")
	r.w.WriteHeader(http.StatusOK)
}

func (r *responder) RespondJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		r.RespondError(err)
		return
	}
	err = r.Respond(http.StatusOK, data)
	if err != nil {
		r.RespondError(err)
		return
	}
}

func (r *responder) RespondJS(script string) {
	r.w.Header().Set("Content-Type", "application/javascript")
	err := r.Respond(http.StatusOK, []byte(script))
	if err != nil {
		r.RespondError(err)
		return
	}
}

func (r *responder) RespondError(err error) {
	err = r.Respond(http.StatusBadRequest, []byte(err.Error()))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (r *responder) Respond(status int, data []byte) error {
	r.w.Header().Set("Access-Control-Allow-Origin", r.url)
	r.w.WriteHeader(status)
	_, err := r.w.Write(data)
	if err != nil {
		return err
	}
	return nil
}
