package repository

import (
	"time"

	"github.com/powiedl/myGoWebApplication/internal/models"
)

type DatabaseRepo interface{
	AllUsers() bool
	InsertReservation(res models.Reservation) (int,error )
	InsertBungalowRestriction(r models.BungalowRestriction) error
	SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowId int) (bool,error)
	SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow,error)
	GetBungalowById(id int)(models.Bungalow,error)
	// User functions
	GetUserByID(id int)(models.User,error) 
  UpdateUser(u models.User)error
  Authenticate(email, testPassword string)(int,string,error)
	AllReservations()([]models.Reservation,error)
	AllNewReservations()([]models.Reservation,error)
	GetReservationByID(id int)(models.Reservation,error)
	UpdateReservation(r models.Reservation) error
	DeleteReservation(id int) error
	UpdateStatusReservation(id,status int) error
}
