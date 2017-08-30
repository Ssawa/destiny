package storage

import (
	"math/rand"

	"github.com/Ssawa/bolt"
	"github.com/Ssawa/destiny/utils"
)

// GetAdage gets a random adage from the database
func GetAdage(db *bolt.DB) (string, error) {
	var adage string
	var keys [][]byte

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(adagesBucket)
		if bucket == nil {
			utils.Verbose.Println("Adages bucket does not exist in the database")
			return nil
		}
		c := bucket.Cursor()

		utils.Verbose.Println("Iterating over keys")
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, k)
		}

		choice := keys[rand.Intn(len(keys))]

		utils.Verbose.Println("Chose: ", choice)

		adage = string(bucket.Get(choice))

		return nil
	})

	return adage, err
}
