package utils

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// OpenReadWrite opens a Bolt DB database with exclusive Read/Write access.
func OpenReadWrite(databasePath string) (*bolt.DB, error) {
	return openDatabase(databasePath, &bolt.Options{ReadOnly: false})
}

// OpenReadOnly opens a Bolt DB database with shared Read only access.
func OpenReadOnly(databasePath string) (*bolt.DB, error) {
	// So when you try to open a Bolt database that doesn't exist it will get a
	// filedescriptor error and return a nil database. *However*
	db, _ := openDatabase(databasePath, &bolt.Options{ReadOnly: true})
	if db != nil {
		db.Close()
	}
	fmt.Println("Opening again")
	return OpenReadWrite(databasePath)
	// If the database doesn't exist yet we'll get an error
	// fmt.Println(os.IsNotExist(err))
	// pathErr := err.(*os.PathError)
	// fmt.Println(pathErr.Op)
	// fmt.Println(reflect.TypeOf(pathErr.Err))
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
