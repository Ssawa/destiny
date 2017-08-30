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
		choiceVal := int64(-1)
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
		utils.Verbose.Println("Picking a random category")
		var tagChoice *bolt.Bucket
		tagChoiceVal := int64(-1)

		tagsBucket := tx.Bucket(tagsKey)
		if tagsBucket == nil {
			utils.Verbose.Println("No tags bucket. Returning")
			return nil
		}
		for _, tag := range tags {
			t := tagsBucket.Bucket([]byte(tag))
			if t != nil {
				val := rand.Int63()
				if tagChoiceVal < val {
					tagChoice = t
					tagChoiceVal = val
				}
			}
		}

		if tagChoice == nil {
			utils.Verbose.Println("No valid tag. Returning")
			return nil
		}

		utils.Verbose.Println("Gathering excludes")
		excludeBuckets := []*bolt.Bucket{}
		for _, ex := range excludes {
			x := tagsBucket.Bucket([]byte(ex))
			if x != nil {
				excludeBuckets = append(excludeBuckets, x)
			}
		}

		utils.Verbose.Println("Picking random element from tag")
		c := tagChoice.Cursor()
		var choice []byte
		choiceVal := int64(-1)
	FindChoice:
		for k, _ := c.First(); k != nil; k, _ = c.Next() {

			for _, x := range excludeBuckets {
				if x.Get(k) != nil {
					continue FindChoice
				}
			}

			val := rand.Int63()
			if choiceVal < val {
				choice = k
				choiceVal = val
			}
		}

		adagesBucket := tx.Bucket(adagesKey)
		if adagesBucket == nil {
			utils.Verbose.Println("Adages bucket does not exist in the database")
			return nil
		}
		utils.Verbose.Println("Chose: ", choice)
		if choice == nil {
			return nil
		}

		var err error
		adage, err = DeserializeAdage(adagesBucket.Get(choice))

		return err
	})

	return adage, err
}
