package gopgx

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

//D - Database instance to be used after running Setup

//Db struct is used to configure and init pg db instance
type Db struct {
	Host, DbName, Port, User, password string
	D                                  *pgx.Conn
}

// Setup connects to database server and checks connection
func (db *Db) Setup() (*pgx.Conn, error) {

	//cfg.SetConfig(host, port, user, password, dbname)

	fmt.Println(" > Connecting to database...")

	var err error
	db.D, err = pgx.Connect(context.Background(), fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		db.Host, db.Port, db.User, db.password, db.DbName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		return nil, err
		//os.Exit(1) //TODO: handle error gracefully or pass to the calling program
	}
	err = db.Ping("Database connected...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping: Failed to run QueryRow for Ping(): %v\n", err)
		return nil, err
		//os.Exit(1) //TODO: handle error gracefully or pass to the calling program
	}
	defer db.D.Close(context.Background())

	return db.D, nil
}

//Ping prints the message passed, used for checking database connection
//TODO add errors as return value
func (db *Db) Ping(msg string) error {
	var showMessage string
	if db.D != nil {
		err := db.D.QueryRow(context.Background(), "select '"+msg+"'").Scan(&showMessage)
		if err != nil {
			return err
		}
		fmt.Println(showMessage)
	} else {
		return errors.New("Db instance not found, Setup database first")
	}
	return nil
}

type passwordSetter func() (string, error)

func (db *Db) SetPassword(f passwordSetter) error {

	pwd, err := f()

	if err != nil {
		return err
	}

	db.password = pwd
	return nil

}

/*
**** Make this package configurable
// Config struct to create database url
// not exporting password and username
type Config struct {
	Host, Port, user, password, DBName string
}

// Config instance
var cfg *Config


// SetConfig to prepare cfg instance for use
// TODO: Update to check and use either env vars or defaults, or env files also
func (c *Config) SetConfig(host, port, user, password, dbname string) {
	cfg.Host = host
	cfg.DBName = dbname
	cfg.user = user
	cfg.password = password
	cfg.Port = port
}


//DbURL creates database url to connect to the db server
func (c *Config) DbURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		c.Host, c.Port, c.user, c.password, c.DBName)
}
*/
