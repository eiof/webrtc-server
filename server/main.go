package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

// callDocument:
// - offer: sdp + type
// - offerCandidates[]: list of ice candidates
// - answerCandidates[]: list of ice candidates
type OfferResponse struct {
	Id string `json:"id"`
}

type Offer struct {
	Sdp  string `json:"sdp"`
	Type string `json:"type"`
}

type IceCandidate struct {
}

type Call struct {
	Offer            Offer
	OfferCandidates  []*IceCandidate
	AnswerCandidates []*IceCandidate
}

func (c *Call) String() string {
	return fmt.Sprint("%+v | %+v | %+v", c.Offer, c.OfferCandidates, c.AnswerCandidates)
}

type CallCollection struct {
	calls map[string]*Call
}

func NewCallCollection() *CallCollection {
	return &CallCollection{
		calls: make(map[string]*Call),
	}
}

func (c *CallCollection) Call(id string) *Call {
	// TODO: handle non-existence
	call, _ := c.calls[id]
	return call
}

func (c *CallCollection) AddCall(call *Call) string {
	id := uuid.NewString()
	c.calls[id] = call
	return id
}

var callCollection = NewCallCollection()

// Create an offer + call with a unique id
// Server sent event for answers
// Answer a call with the unique id
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/offer", func(w http.ResponseWriter, r *http.Request) {
		var offer Offer
		if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
			fmt.Printf("Error unmarshaling request: %s\n", err)
		}
		fmt.Printf("Request offer: %+v", offer)
		call := &Call{
			Offer: offer,
		}
		call_id := callCollection.AddCall(call)
		fmt.Printf("CallCollection (new_id: %s): %+v\n", call_id, callCollection)
		o := &OfferResponse{Id: "test"}
		res, _ := json.Marshal(o)
		w.Write(res)
	})
	http.ListenAndServe(":3000", r)
}
