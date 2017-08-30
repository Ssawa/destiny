package storage

import (
	"math/rand"

	"github.com/Ssawa/bolt"
	"github.com/Ssawa/destiny/utils"
)

// GetAdageFromAll gets an Adage from anywhere in the database
func GetAdageFromAll(db *bolt.DB) (*Adage, error) {
	var adage *Adage

	err := db.View(func(tx *bolt.Tx) error {
		adagesBucket := tx.Bucket(adagesKey)
		if adagesBucket == nil {
			utils.Verbose.Println("Adages bucket does not exist in the database")
			return nil
		}
		c := adagesBucket.Cursor()

		utils.Verbose.Println("Iterating over keys")
		var choice []byte
		choiceVal := int64(0)
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			val := rand.Int63()
			if choiceVal < val {
				choice = k
				choiceVal = val
			}
		}

		utils.Verbose.Println("Chose: ", choice)

		var err error
		adage, err = DeserializeAdage(adagesBucket.Get(choice))

		return err
	})

	return adage, err
}

func GetAdageFromCategories(db *bolt.DB, tags []string, excludes []string) (*Adage, error) {
	var adage *Adage

	err := db.View(func(tx *bolt.Tx) error {
		return nil
	})

	return adage, err
}
