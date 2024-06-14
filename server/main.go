package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// callDocument:
// - offer: sdp + type
// - offerCandidates[]: list of ice candidates
// - answerCandidates[]: list of ice candidates
type OfferResponse struct {
	Id string `json:"id"`
}

// Create an offer + call with a unique id
// Server sent event for answers
// Answer a call with the unique id
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/offer", func(w http.ResponseWriter, r *http.Request) {
		o := &OfferResponse{Id: "test"}
		res, _ := json.Marshal(o)
		w.Write(res)
	})
	http.ListenAndServe(":3000", r)
}
