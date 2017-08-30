package utils

import (
	"os"

	"github.com/Ssawa/bolt"
)

// OpenReadWrite opens a Bolt DB database with exclusive Read/Write access.
func OpenReadWrite(databasePath string) (*bolt.DB, error) {
	return openDatabase(databasePath, &bolt.Options{ReadOnly: false})
}

// OpenReadOnly opens a Bolt DB database with shared Read only access.
func OpenReadOnly(databasePath string) (*bolt.DB, error) {
	db, err := openDatabase(databasePath, &bolt.Options{ReadOnly: true})
	if err != nil {
		if _, isWriteError := err.(*os.PathError); isWriteError {
			// BoltDB can't initialize a new database in ReadOnly mode (in fact
			// it gets caught in a weird state, which is why we're using my personal
			// fork until the PR gets accepted). We could either return an error
			// or we can create it using ReadWrite mode. Both have their pros and
			// cons. We'll be creating the database, as this simplifies some error
			// handling (the main point of this function) and also forces us to make
			// sure that we're handling incomplete database initialization throughout
			// our entire codebase
			Verbose.Println("Tried to open a non-existant database.")
			Verbose.Println("Creating database in ReadWrite mode")
			db, err = OpenReadWrite(databasePath)
			if err != nil {
				Verbose.Println("Could not create database", err)
			} else {
				Verbose.Println("Closing ReadWrite database and opening again as ReadOnly")
				db.Close()
				db, err = openDatabase(databasePath, &bolt.Options{ReadOnly: true})
			}
		}
	}
	return db, err
}

// openDatabase attempts to open the bolt database with the given options
// while ensuring that the containing folder exists
func openDatabase(databasePath string, options *bolt.Options) (*bolt.DB, error) {
	if err := EnsureDirectory(databasePath); err != nil {
		return nil, err
	}

	db, err := bolt.Open(databasePath, 0600, options)
	return db, err
}
