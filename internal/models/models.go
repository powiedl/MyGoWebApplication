package models

import "time"

// User model
type User struct {
	ID int
	FullName string
	Email string
	Password string
	Role int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Bungalow model
type Bungalow struct {
	ID int
  BungalowName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Reservation model
type Reservation struct {
	ID int
	FullName string
	Email string
	Phone string
	StartDate time.Time
	EndDate time.Time
	BungalowID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Processed int
	Bungalow Bungalow
}

// Restriction model
type Restriction struct {
	ID int
  RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BungalowRestriction model
type BungalowRestriction struct {
	ID int
	StartDate time.Time
	EndDate time.Time
	BungalowID int
	ReservationID int
	RestrictionID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Bungalow Bungalow
	Reservation Reservation
	Restriction Restriction
}