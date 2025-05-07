package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbconn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

func ConnectSQL (dsn string) (*DB, error) {
  d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)
	dbconn.SQL = d

	if testDB(d)	!= nil {
		return nil,err
	}

	return dbconn,nil
}

func NewDatabase(dsn string) (*sql.DB,error) {
	db, err := sql.Open("pgx",dsn)
	if err != nil {
		return nil,err
	}

	if testDB(db)	!= nil {
		return nil,err
	}
	return db,nil
}

func testDB(d *sql.DB) error {
	err := d.Ping(); 
  return err
}