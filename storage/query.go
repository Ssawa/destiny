package storage

import (
	"math/rand"

	"github.com/cj-dimaggio/bolt"
	"github.com/cj-dimaggio/destiny/utils"
)

// GetAdageFromAll gets an Adage from anywhere in the database
func GetAdageFromAll(db *bolt.DB, excludes []string) (*Adage, error) {
	var adage *Adage

	err := db.View(func(tx *bolt.Tx) error {
		adagesBucket := tx.Bucket(adagesKey)
		if adagesBucket == nil {
			utils.Verbose.Println("Adages bucket does not exist in the database")
			return nil
		}

		tagsBucket := tx.Bucket(tagsKey)
		if tagsBucket == nil {
			utils.Verbose.Println("No tags bucket. Returning")
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

		c := adagesBucket.Cursor()

		utils.Verbose.Println("Iterating over keys")
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

func GetAdageFromCategories(db *bolt.DB, tags []string, excludes []string, exclusive bool) (*Adage, error) {
	var adage *Adage

	err := db.View(func(tx *bolt.Tx) error {
		tagsBucket := tx.Bucket(tagsKey)
		if tagsBucket == nil {
			utils.Verbose.Println("No tags bucket. Returning")
			return nil
		}
		utils.Verbose.Println("Gathering tags that exist")
		tagBuckets := []*bolt.Bucket{}
		for _, tag := range tags {
			t := tagsBucket.Bucket([]byte(tag))
			if t != nil {
				tagBuckets = append(tagBuckets, t)
			}
		}

		if len(tagBuckets) == 0 {
			utils.Verbose.Println("No valid tags. Returning")
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

		utils.Verbose.Println("Picking random element")

		adagesBucket := tx.Bucket(adagesKey)
		if adagesBucket == nil {
			utils.Verbose.Println("Adages bucket does not exist in the database")
			return nil
		}
		c := adagesBucket.Cursor()

		var choice []byte
		choiceVal := int64(-1)
	FindChoice:
		for k, _ := c.First(); k != nil; k, _ = c.Next() {

			for _, x := range excludeBuckets {
				if x.Get(k) != nil {
					continue FindChoice
				}
			}

			for _, t := range tagBuckets {
				if t.Get(k) == nil {
					continue FindChoice
				} else if !exclusive {
					break
				}
			}

			val := rand.Int63()
			if choiceVal < val {
				choice = k
				choiceVal = val
			}
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
