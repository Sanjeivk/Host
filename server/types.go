package main

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ReservationStatus string

const (
	StatusPending   ReservationStatus = "pending"
	StatusCancelled ReservationStatus = "cancelled"
	StatusConfirmed ReservationStatus = "confirmed"
)

type CreateReservationRequest struct {
	AccountID   string `json:"account_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	PartySize   string `json:"party_size"`
	Notes       string `json:"notes"`
}

type Reservation struct {
	ID          string            `json:"id"`
	ListingID   string            `json:"listing_id"`
	EventDate   time.Time         `json:"event_date"`
	AccountID   string            `json:"account_id"`
	FirstName   string            `json:"firstName"`
	LastName    string            `json:"lastName"`
	PhoneNumber string            `json:"phoneNumber"`
	Email       string            `json:"email"`
	PartySize   string            `json:"party_size"`
	Status      ReservationStatus `json:"status"`
	Notes       string            `json:"notes"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	DeletedAt   *time.Time        `json:"deletedAt"`
}

type CreateListingRequest struct {
	Host      string    `json:"host"`
	Address   string    `json:"address"`
	Occasion  string    `json:"occasion"`
	Pg        bool      `json:"pg"`
	Byod      bool      `json:"byod"`
	Notes     string    `json:"notes"`
	EventDate time.Time `json:"event_date"`
}

type Listing struct {
	ID             string     `json:"id"`
	HostID         string     `json:"host_id"`
	Address        string     `json:"address"`
	NumberOfGuests int64      `json:"number_of_guests"`
	Occasion       string     `json:"occasion"`
	Pg             bool       `json:"pg"`
	Byod           bool       `json:"byod"`
	Notes          string     `json:"notes"`
	Review         int64      `json:"review"`
	EventDate      time.Time  `json:"event_date"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
}

type CreateAccountRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type Account struct {
	ID                string     `json:"id"`
	FirstName         string     `json:"firstName"`
	LastName          string     `json:"lastName"`
	Email             string     `json:"email"`
	EncryptedPassword string     `json:"-"`
	PhoneNumber       string     `json:"phoneNumber"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt"`
}

func (a *Account) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

func newAccount(FirstName string, LastName string, Email string, PhoneNumber string, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	id := uuid.New()
	idStr := id.String()

	if err != nil {
		return nil, err
	}
	return &Account{
		ID:                idStr,
		FirstName:         FirstName,
		LastName:          LastName,
		Email:             Email,
		PhoneNumber:       PhoneNumber,
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}, nil
}

func newListing(hostId string, address string, occasion string, pg bool, byod bool, notes string, eventDate time.Time) (*Listing, error) {
	id := uuid.New()
	idStr := id.String()

	return &Listing{
		ID:             idStr,
		HostID:         hostId,
		Address:        address,
		NumberOfGuests: 0,
		Occasion:       occasion,
		Pg:             pg,
		Byod:           byod,
		Notes:          notes,
		Review:         5,
		EventDate:      eventDate,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}, nil
}

func newReservation(listingID, accountID, firstName, lastName, phoneNumber, email, partySize, notes string,
	eventDate time.Time) (*Reservation, error) {
	id := uuid.New()
	idStr := id.String()

	return &Reservation{
		ID:          idStr,
		ListingID:   listingID,
		EventDate:   eventDate,
		AccountID:   accountID,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
		Email:       email,
		PartySize:   partySize,
		Status:      StatusPending,
		Notes:       notes,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}
