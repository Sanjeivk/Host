package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountbyID(int) (*Account, error)
	CreateListing(*Listing) error
	GetListings() ([]*Listing, error)
	GetListingByID(int) (*Listing, error)
	CreateReservation(*Reservation) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=hostess sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.CreateAccountTable(); err != nil {
		return err
	}
	if err := s.CreateListingTable(); err != nil {
		return err
	}
	if err := s.CreateReservationTable(); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateAccountTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS account (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			email VARCHAR(100),
			phone_number VARCHAR(100),
			encrypted_password VARCHAR(100),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			deleted_at TIMESTAMP DEFAULT NULL
		)
	`)
	return err
}

func (s *PostgresStore) CreateListingTable() error {
	//		CREATE TYPE event_type AS ENUM ('MMA', 'Formula 1', 'Tennis', 'Golf', 'Football', 'Soccer', 'Basketball', 'Hockey', 'Baseball');
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS listing (
			id SERIAL PRIMARY KEY,
			host_id SERIAL,
			street VARCHAR(255),
			city VARCHAR(255),
			state VARCHAR(255),
			postal_code VARCHAR(255),
			country VARCHAR(255),
			number_of_guests INT DEFAULT 0,
			occasion VARCHAR(100),
			pg BOOL,
			byod BOOL,
			notes VARCHAR(1000),
			review INT,
			event_date TIMESTAMP NOT NULL,
			event_type event_type,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			deleted_at TIMESTAMP DEFAULT NULL,
			FOREIGN KEY (host_id) REFERENCES account(id)
		)
	`)
	return err
}

func (s *PostgresStore) CreateReservationTable() error {
	//CREATE TYPE reservation_status AS ENUM ('pending', 'cancelled', 'confirmed');
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS reservation (
		id VARCHAR(100) PRIMARY KEY,
		listing_id SERIAL,
		event_date TIMESTAMP NOT NULL,
		account_id SERIAL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		phone_number VARCHAR(100),
		email VARCHAR(100),
		party_size SERIAL,
		status reservation_status DEFAULT 'pending',
		notes VARCHAR(1000),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		deleted_at TIMESTAMP DEFAULT NULL,
		FOREIGN KEY (account_id) REFERENCES account(id),
		FOREIGN KEY (listing_id) REFERENCES listing(id)
	)
`)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	Query := `INSERT INTO account (
		first_name,
		last_name,
		email,
		phone_number,
		encrypted_password,
		created_at,
		updated_at
	)
	VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.Exec(Query, acc.FirstName, acc.LastName, acc.Email, acc.PhoneNumber, acc.EncryptedPassword, acc.CreatedAt, acc.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) GetAccountbyID(id int) (*Account, error) {
	return nil, nil
}

//listing

func (s *PostgresStore) CreateListing(listing *Listing) error {
	Query := `INSERT INTO listing (
		host_id,
		street,
		city,
		state,
		postal_code,
		country,
		number_of_guests,
		occasion,
		pg,
		byod,
		notes,
		review,
		event_date,
		event_type,
		created_at,
		updated_at
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	fmt.Println(listing.ID, listing.HostID, listing.Street, listing.City, listing.State, listing.PostalCode, listing.Country, listing.NumberOfGuests, listing.Occasion, listing.Pg, listing.Byod, listing.Notes, listing.Review, listing.EventDate, listing.CreatedAt, listing.UpdatedAt)
	_, err := s.db.Exec(Query, listing.HostID, listing.Street, listing.City, listing.State, listing.PostalCode, listing.Country, listing.NumberOfGuests, listing.Occasion, listing.Pg, listing.Byod, listing.Notes, listing.Review, listing.EventDate, listing.EventType, listing.CreatedAt, listing.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetListingByID(id int) (*Listing, error) {
	rows, err := s.db.Query(`SELECT * FROM listing WHERE id=$1`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntolisting(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetListings() ([]*Listing, error) {
	rows, err := s.db.Query(`SELECT * FROM listing`)

	if err != nil {
		return nil, err
	}

	listings := []*Listing{}

	for rows.Next() {
		listing, err := scanIntolisting(rows)

		if err != nil {
			return nil, err
		}

		listings = append(listings, listing)
	}
	return listings, nil
}

//reservations

func (s *PostgresStore) CreateReservation(reservation *Reservation) error {
	Query := `INSERT INTO reservation (
		id,
		listing_id,
		event_date,
		account_id,
		first_name,
		last_name,
		phone_number,
		email,
		party_size,
		status,
		notes,
		created_at,
		updated_at
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	fmt.Println(reservation.ID, reservation.ListingID, reservation.EventDate, reservation.AccountID, reservation.FirstName, reservation.LastName, reservation.PhoneNumber, reservation.Email, reservation.PartySize, reservation.Status, reservation.Notes, reservation.CreatedAt, reservation.UpdatedAt)
	_, err := s.db.Exec(Query, reservation.ID, reservation.ListingID, reservation.EventDate, reservation.AccountID, reservation.FirstName, reservation.LastName, reservation.PhoneNumber, reservation.Email, reservation.PartySize, reservation.Status, reservation.Notes, reservation.CreatedAt, reservation.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func scanIntolisting(rows *sql.Rows) (*Listing, error) {
	listing := new(Listing)
	err := rows.Scan(&listing.ID, &listing.HostID, &listing.Street, &listing.City, &listing.State, &listing.PostalCode, &listing.Country, &listing.NumberOfGuests, &listing.Occasion, &listing.Pg, &listing.Byod, &listing.Notes, &listing.Review, &listing.EventDate, &listing.EventType, &listing.CreatedAt, &listing.UpdatedAt,
		&listing.DeletedAt)
	return listing, err
}
