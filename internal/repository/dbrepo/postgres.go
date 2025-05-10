package dbrepo

import (
	"context"
	"time"

	"github.com/powiedl/myGoWebApplication/internal/models"
)

func (m *postgresDBRepo)AllUsers() bool {
	return true
}

// InsertReservation inserts a record into the reservations table
func (m *postgresDBRepo)InsertReservation(res models.Reservation) (int,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	statement := `INSERT INTO public.reservations (full_name,email,phone,start_date,end_date,bungalow_id,created_at,updated_at)
	  VALUES ($1,$2,$3,$4,$5,$6,$7,$8) returning id`
	
	var newId int
	
	err := m.DB.QueryRowContext(ctx,statement,
		res.FullName,res.Email,res.Phone,res.StartDate,res.EndDate,res.BungalowID,time.Now(),time.Now(),
	).Scan(&newId)
	if err != nil {
		return 0,err
	}

	return newId,nil
}

// InsertBungalowRestriction inserts a record into bungalow_restrictions table
func (m *postgresDBRepo)InsertBungalowRestriction(r models.BungalowRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	statement:=`INSERT INTO public.bungalow_restrictions (start_date,end_date,bungalow_id,reservation_id,created_at,updated_at,restriction_id)
	  VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_,err := m.DB.ExecContext(ctx,statement,
			r.StartDate,r.EndDate,r.BungalowID,r.ReservationID,time.Now(),time.Now(),r.RestrictionID)

	if err != nil {
		return err
	}
	//log.Println("InsertBungalowRestriction, executed statement '",statement,"'")
	return nil
}

// SearchAvailabilityByDatesByBungalowID checks if there is availabiltity for a date range for a given BungalowID, returns false if not
func (m *postgresDBRepo)SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowId int) (bool,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	var numRows int

	query:=`SELECT count(id) FROM public.bungalow_restrictions
    WHERE $1 <= end_date AND $2 >= start_date AND bungalow_id=$3`

  err := m.DB.QueryRowContext(ctx,query,
			start,end,bungalowId).Scan(&numRows)
	if err != nil {
		return false,err
	}
	return numRows == 0,nil
}

// SearchAvailabilityByDatesForAllBungalows returns a slice of available bungalows, if any for a queried range of dates
func (m *postgresDBRepo)SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	var bungalows []models.Bungalow

  query := `SELECT b.id, b.bungalow_name FROM public.bungalows b
    where b.id not in (
	    SELECT bungalow_id from public.bungalow_restrictions br
	    where $1 <= br.end_date and $2 >= br.start_date
    )`

	rows,err := m.DB.QueryContext(ctx,query,start,end)
	if err != nil {
		return bungalows,err
	}
	for rows.Next() {
		var bungalow models.Bungalow
		err := rows.Scan(
			&bungalow.ID,
			&bungalow.BungalowName,
		)
		if err != nil {
			return bungalows,err
		}
		bungalows = append(bungalows,bungalow)
	}
	if err=rows.Err(); err != nil {
		return bungalows,err
	}
	return bungalows,nil
}

// Get bungalow by id
func (m *postgresDBRepo)GetBungalowById(id int)(models.Bungalow,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	var bungalow models.Bungalow

  query := `SELECT b.id, b.bungalow_name, b.created_at, b.updated_at FROM public.bungalows b where b.id=$1`

	row := m.DB.QueryRowContext(ctx,query,id)
	err := row.Scan(&bungalow.ID,&bungalow.BungalowName,&bungalow.CreatedAt,&bungalow.UpdatedAt)
	if err != nil {
		return bungalow,err
	}
	return bungalow,nil
}