package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/powiedl/myGoWebApplication/internal/models"
)

func (m *testDBRepo)AllUsers() bool {
	return true
}

// InsertReservation inserts a record into the reservations table
func (m *testDBRepo)InsertReservation(res models.Reservation) (int,error) {
	if res.BungalowID == 99 {
		return 0,errors.New("some error")
	}
	return 1,nil
}

// InsertBungalowRestriction inserts a record into bungalow_restrictions table
func (m *testDBRepo)InsertBungalowRestriction(r models.BungalowRestriction) error {
	if r.BungalowID == 999 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByBungalowID checks if there is availabiltity for a date range for a given BungalowID, returns false if not
func (m *testDBRepo)SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowId int) (bool,error) {
	layout := "2006-01-02"
	testDateToFail,err := time.Parse(layout,"2038-01-01")
	if err != nil {
		log.Println(err)
	}
	if start == testDateToFail { // if start is 2038-01-01 return an error as availability
		return false,errors.New("some error")
	}

	strFalse := "2036-12-31" // after 2036-12-31 always return false as availability
	tFalse,err := time.Parse(layout,strFalse)
	if err != nil {
		log.Println(err)
	}
	if start.After(tFalse) {
		return false, nil
	}

	return true,nil
}

// SearchAvailabilityByDatesForAllBungalows returns a slice of available bungalows, if any for a queried range of dates
func (m *testDBRepo)SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow,error) {
	var bungalows []models.Bungalow
	layout := "2006-01-02"
	testDateToFail,err := time.Parse(layout,"2038-01-01")
	if err != nil {
		log.Println(err)
	}
	if start == testDateToFail { // if start is 2038-01-01 return an error as availability
		return bungalows,errors.New("some error")
	}

	strFalse := "2036-12-31" // after 2036-12-31 always return false as availability
	tFalse,err := time.Parse(layout,strFalse)
	if err != nil {
		log.Println(err)
	}
	if start.After(tFalse) {
		return bungalows, nil
	}

	bungalow := models.Bungalow{
		BungalowName: "The Solitude Stack",
	}
	return []models.Bungalow{bungalow},nil
}

// Get bungalow by id
func (m *testDBRepo)GetBungalowById(id int)(models.Bungalow,error) {
	var bungalow models.Bungalow
	if id > 3 {
		return bungalow,errors.New("an error occured")
	}
	return bungalow,nil
}

// #region User functions
func (m *testDBRepo) GetUserByID(id int)(models.User,error) {
	var u models.User
	return u,nil
}

// UpdateUser updates basic user data in the database
func (m *testDBRepo) UpdateUser(u models.User)error {
	return nil
}
// Authenticate authenticates a user by data
func (m *testDBRepo) Authenticate(email, testPassword string)(int,string,error) {
	return 1,"",nil
}

func (m *testDBRepo) AllReservations()([]models.Reservation,error) {
	var reservations []models.Reservation
	return reservations,nil
}

func (m *testDBRepo) AllNewReservations()([]models.Reservation,error) {
	var reservations []models.Reservation
	return reservations,nil
}

func (m *testDBRepo) GetReservationByID(id int)(models.Reservation,error) {
	var res models.Reservation
	return res,nil
}

func (m *testDBRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateStatusReservation(id,status int) error {
	return nil
}
// #endregion