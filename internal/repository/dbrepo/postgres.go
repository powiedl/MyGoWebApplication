package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/powiedl/myGoWebApplication/internal/models"
	"golang.org/x/crypto/bcrypt"
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
	defer rows.Close()

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

// #region user

// GetUserByID returns user data by id
func (m *postgresDBRepo) GetUserByID(id int)(models.User,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	query := `
	SELECT id, full_name, email, password, role, created_at, updated_at FROM users WHERE id = $1
	`
	row := m.DB.QueryRowContext(ctx,query,id)

	var u models.User

	err := row.Scan(&u.ID,&u.FullName,&u.Email,&u.Password,&u.Role,&u.CreatedAt,&u.UpdatedAt)
	if err != nil {
		return u,err
	}
	return u,nil
}

// UpdateUser updates basic user data in the database
func (m *postgresDBRepo) UpdateUser(u models.User)error {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

  statement := `
	UPDATE users SET full_name=$2,email=$3,role=$4,updated_at=$5 WHERE id=$1
	`
	_, err := m.DB.ExecContext(ctx,statement,u.ID,u.FullName,u.Email,u.Role,time.Now())
	if err != nil {
		return err
	}
	return nil
}
// Authenticate authenticates a user by data
func (m *postgresDBRepo) Authenticate(email, testPassword string)(int,string,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	var id int
	var passwordHash string

	row := m.DB.QueryRowContext(ctx,"SELECT id,password FROM users WHERE email=$1",email)
	err := row.Scan(&id,&passwordHash)
	if err != nil {
		return id,"",err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash),[]byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return id,"",errors.New("wrong password")
	} else if err != nil {
		return id,"",err
	}
	return id,passwordHash,nil
}
// #endregion

// AllReservations builds and returns a slice of all reservations from the database
func (m *postgresDBRepo) AllReservations()([]models.Reservation,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	var reservations[]models.Reservation

	query := `
	  select r.id, r.full_name,r.email,r.phone,r.start_date,r.end_date,r.bungalow_id,r.created_at,r.updated_at,r.status, b.id,b.bungalow_name
	  from reservations r
	  left join bungalows b on (r.bungalow_id = b.id)
	  order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx,query)
	if err != nil {
		return reservations,err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.BungalowID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Status,
			&i.Bungalow.ID,
			&i.Bungalow.BungalowName,
		)
		if err != nil {
		  return reservations,err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations,err
	}
	return reservations,nil
}

// AllNewReservations builds and returns a slice of all reservations from the database
func (m *postgresDBRepo) AllNewReservations()([]models.Reservation,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	var reservations[]models.Reservation

	query := `
	  select r.id, r.full_name,r.email,r.phone,r.start_date,r.end_date,r.bungalow_id,r.created_at,r.updated_at,r.status, b.id,b.bungalow_name
	  from reservations r
	  left join bungalows b on (r.bungalow_id = b.id)
		where r.status = 0
	  order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx,query)
	if err != nil {
		return reservations,err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.BungalowID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Status,
			&i.Bungalow.ID,
			&i.Bungalow.BungalowName,
		)
		if err != nil {
		  return reservations,err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations,err
	}
	return reservations,nil
}

// GetReservationByID builds and returns a slice of all reservations from the database
func (m *postgresDBRepo) GetReservationByID(id int)(models.Reservation,error) {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
	  select r.id, r.full_name,r.email,r.phone,r.start_date,r.end_date,r.bungalow_id,r.created_at,r.updated_at,r.status, b.id,b.bungalow_name
	  from reservations r
	  left join bungalows b on (r.bungalow_id = b.id)
		where r.id=$1`

	row := m.DB.QueryRowContext(ctx,query,id)

		err := row.Scan(
			&res.ID,
			&res.FullName,
			&res.Email,
			&res.Phone,
			&res.StartDate,
			&res.EndDate,
			&res.BungalowID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.Status,
			&res.Bungalow.ID,
			&res.Bungalow.BungalowName,
		)
		if err != nil {
		  return res,err
		}
	return res,nil
}

// UpdateReservation updates reservation data in the database
func (m *postgresDBRepo) UpdateReservation(r models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

  statement := `
	UPDATE reservations SET full_name=$2,email=$3,phone=$4,updated_at=$5 WHERE id=$1
	`
	_, err := m.DB.ExecContext(ctx,statement,r.ID,r.FullName,r.Email,r.Phone,time.Now())
	if err != nil {
		return err
	}
	return nil
}

// DeleteReservation deletes a reservation in the database
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

  statement := `
	DELETE FROM reservations WHERE id=$1
	`
	_, err := m.DB.ExecContext(ctx,statement,id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStatusReservation updates the status of a reservation in the database
func (m *postgresDBRepo) UpdateStatusReservation(id,status int) error {
	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

  statement := `
	UPDATE reservations SET status=$2 WHERE id=$1
	`
	_, err := m.DB.ExecContext(ctx,statement,id,1)
	if err != nil {
		return err
	}
	return nil
}