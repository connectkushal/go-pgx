package gopgx

import (
	"context"
	"errors"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v5"
)

//Db struct is used to configuure database and store a Connection instance
type Db struct {
	Host, DbName, Port, User, password string
	Instance                           *pgx.Conn
}

// Setup connects to the postgres database server and verifies the connection using Ping
func (db *Db) Setup(ctx context.Context) error {

	fmt.Println(" > Connecting to database...")

	var err error
	db.Instance, err = pgx.Connect(ctx, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		db.Host, db.Port, db.User, db.password, db.DbName))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		return err
	}

	err = db.Ping("Connected to db " + db.DbName + "@" + db.Host + ":" + db.Port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping: Failed to run QueryRow for Ping(): %v\n", err)
		return err
	}

	return nil
}

//Ping prints the message passed, can be used for checking the database connection
func (db *Db) Ping(msg string) error {
	var showMessage string
	if db.Instance != nil {
		err := db.Instance.QueryRow(context.Background(), "select '"+msg+"'").Scan(&showMessage)
		if err != nil {
			return err
		}
		fmt.Println(showMessage)
	} else {
		return errors.New("Db instance not found, Setup database first")
	}
	return nil
}

// To enable setting password via env var, api or any other way
type passwordSetter func() (string, error)

// SetPassword sets the db password through passing a function which returns a string and an error
func (db *Db) SetPassword(f passwordSetter) error {
	pwd, err := f()
	if err != nil {
		return err
	}
	db.password = pwd
	return nil

}
