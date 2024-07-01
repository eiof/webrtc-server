package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// callDocument:
// - offer: sdp + type
// - offerCandidates[]: list of ice candidates
// - answerCandidates[]: list of ice candidates
type OfferResponse struct {
	Id string `json:"id"`
}

type AnswerResponse struct {
	Id string `json:"id"`
}

type Offer struct {
	Id        string `json:"id"`
	Sdp       string `json:"sdp"`
	Type      string `json:"type"`
	Candidate string `json:"candidate"`
}

type Answer struct {
	CallID    string `json:"call_id"`
	Sdp       string `json:"sdp"`
	Type      string `json:"type"`
	Candidate string `json:"candidate"`
}

type IceCandidate struct {
	Candidate        string `json:"candidate"`
	SdpMLineIndex    int    `json:"sdpMLineIndex"`
	SdpMid           string `json:"sdpMid"`
	UsernameFragment string `json:"usernameFragment"`
}

type Call struct {
	Offer            Offer
	OfferCandidates  []*IceCandidate
	AnswerCandidates []*IceCandidate
}

func (c *Call) String() string {
	return fmt.Sprint("%+v | %+v | %+v", c.Offer, c.OfferCandidates, c.AnswerCandidates)
}

func (c *Call) AddAnswerCandidate(candidate *IceCandidate) {
	c.AnswerCandidates = append(c.AnswerCandidates, candidate)
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
	// id := uuid.NewString()
	c.calls[call.Offer.Id] = call
	return call.Offer.Id
}

var callCollection = NewCallCollection()

// Create an offer + call with a unique id
// Server sent event for answers
// Answer a call with the unique id
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))
	r.Post("/offer", func(w http.ResponseWriter, r *http.Request) {
		var offer Offer
		if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
			fmt.Printf("Error unmarshaling request: %s\n", err)
		}
		fmt.Printf("Request offer: %+v", offer)
		call := &Call{
			Offer: offer,
		}
		callID := callCollection.AddCall(call)
		fmt.Printf("CallCollection (new_id: %s): %+v\n", callID, callCollection)
		// o := &OfferResponse{Id: callID}
		// res, _ := json.Marshal(o)
		res, _ := json.Marshal(call)
		w.Write(res)
	})

	r.Post("/answer", func(w http.ResponseWriter, r *http.Request) {
		var answer Answer
		if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
			fmt.Printf("Error unmarshaling request: %s\n", err)
		}
		fmt.Printf("Request answer: %+v", answer)

		call := callCollection.Call(answer.CallID)
		call.AddAnswerCandidate(&IceCandidate{Candidate: answer.Candidate})
	})
	http.ListenAndServe(":3000", r)
}
