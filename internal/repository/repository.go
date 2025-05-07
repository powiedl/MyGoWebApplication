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
}
