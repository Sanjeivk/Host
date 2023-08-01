package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func newAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc(("/account"), makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc(("/listing/{id}"), makeHTTPHandleFunc(s.handleListing))
	router.HandleFunc(("/listings"), makeHTTPHandleFunc(s.handleGetListings))
	router.HandleFunc(("/listing/{id}/reservation"), makeHTTPHandleFunc(s.handleReservation))
	// router.HandleFunc(("/account"), makeHTTPHandleFunc(s.handleAccount))
	// router.HandleFunc(("/account/{id}"), withJWTAuth(makeHTTPHandleFunc(s.handleAccountByID), s.store))
	// router.HandleFunc(("/transfer"), makeHTTPHandleFunc(s.handleTransfer))

	log.Println("JSON Server is running on port", s.listenAddr)

	// Create a new CORS handler
	corsHandler := cors.Default()

	// Wrap the router with the CORS handler
	handler := corsHandler.Handler(router)

	http.ListenAndServe(s.listenAddr, handler)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	return s.handleCreateAccount(w, r)
}

func (s *APIServer) handleListing(w http.ResponseWriter, r *http.Request) error {
	return s.handleCreateListing(w, r)
}
func (s *APIServer) handleGetListings(w http.ResponseWriter, r *http.Request) error {
	listings, err := s.store.GetListings()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, listings)
}

func (s *APIServer) handleReservation(w http.ResponseWriter, r *http.Request) error {
	return s.handleCreateReservation(w, r)
}

func (s *APIServer) handleCreateListing(w http.ResponseWriter, r *http.Request) error {
	//retrieving hostID
	id, err := getID(r)
	idStr := strconv.Itoa(id)

	if err != nil {
		return err
	}

	req := new(CreateListingRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	listing, err := newListing(idStr, req.Street, req.City, req.State, req.PostalCode, req.Country, req.Occasion, req.Pg, req.Byod, req.Notes, req.EventDate, req.EventType)

	if err != nil {
		return err
	}

	err = s.store.CreateListing(listing)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, listing)

}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	acc, err := newAccount(req.FirstName, req.LastName, req.Email, req.PhoneNumber, req.Password)

	if err != nil {
		return err
	}

	err = s.store.CreateAccount(acc)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleCreateReservation(w http.ResponseWriter, r *http.Request) error {
	listingId, err := getID(r)
	listingIdStr := strconv.Itoa(listingId)

	if err != nil {
		return err
	}

	req := new(CreateReservationRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	listing, err := s.store.GetListingByID(listingId)

	if err != nil {
		return err
	}

	reservation, err := newReservation(listingIdStr, req.AccountID, req.FirstName, req.LastName, req.PhoneNumber, req.Email, req.PartySize, req.Notes, listing.EventDate)

	if err != nil {
		return err
	}

	err = s.store.CreateReservation(reservation)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, reservation)
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
