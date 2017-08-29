package utils

import "github.com/boltdb/bolt"

// OpenReadWrite opens a Bolt DB database with exclusive Read/Write access.
func OpenReadWrite(databasePath string) (*bolt.DB, error) {
	return openDatabase(databasePath, &bolt.Options{
		ReadOnly: false,
	})
}

// OpenReadOnly opens a Bolt DB database with shared Read only access.
func OpenReadOnly(databasePath string) (*bolt.DB, error) {
	return openDatabase(databasePath, &bolt.Options{
		ReadOnly: true,
	})
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
