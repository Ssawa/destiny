package storage

import (
	"github.com/Ssawa/bolt"
	"github.com/Ssawa/destiny/utils"
	uuid "github.com/satori/go.uuid"
)

// AddAdage inserts a new adage to the database
func AddAdage(db *bolt.DB, adage string, tags []string) error {
	utils.Verbose.Println("Inserting adage")

	id := uuid.NewV1()
	utils.Verbose.Println("UUID generated: ", id)

	utils.Verbose.Println("Starting transaction")
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(adagesBucket)
		if err != nil {
			return err
		}

		utils.Verbose.Println("Saving to database")
		err = bucket.Put(id.Bytes(), []byte(adage))
		return nil
	})
}
