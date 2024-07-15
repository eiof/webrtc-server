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

type OfferCandidateRequest struct {
	CallID    string        `json:"callId"`
	Candidate *IceCandidate `json:"candidate"`
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
	return fmt.Sprintf("%+v | %+v | %+v", c.Offer, c.OfferCandidates, c.AnswerCandidates)
}

func (c *Call) AddAnswerCandidate(candidate *IceCandidate) {
	c.AnswerCandidates = append(c.AnswerCandidates, candidate)
}

func (c *Call) AddOfferCandidate(candidate *IceCandidate) {
	c.OfferCandidates = append(c.OfferCandidates, candidate)
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
	if call, ok := c.calls[id]; !ok {
		call = &Call{}
		c.AddCall(call)
		return call
	} else {
		return call
	}
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
		fmt.Printf("Request offer: %+v\n", offer)
		call := callCollection.Call(offer.Id)
		call.Offer = offer

		fmt.Printf("CallCollection (new_id: %s): %+v\n", offer.Id, callCollection)
		// o := &OfferResponse{Id: callID}
		// res, _ := json.Marshal(o)
		res, _ := json.Marshal(call)
		w.Write(res)
	})

	r.Post("/offer-candidate", func(w http.ResponseWriter, r *http.Request) {
		var candidateRequest *OfferCandidateRequest
		if err := json.NewDecoder(r.Body).Decode(&candidateRequest); err != nil {
			fmt.Printf("Error unmarshaling request: %s\n", err)
		}
		fmt.Printf("Request offer-candidate: %+v\n", candidateRequest)
		call := callCollection.Call(candidateRequest.CallID)
		if call == nil {
			fmt.Printf("CALL IS NIL!!!!!!!\n")
		}
		fmt.Printf("******Offer candidate: %+v\n", candidateRequest.Candidate)
		call.AddOfferCandidate(candidateRequest.Candidate)
		fmt.Printf("CallCollection (new_id: %s): %+v\n", candidateRequest.CallID, callCollection)
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
