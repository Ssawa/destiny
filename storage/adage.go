package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	"github.com/Ssawa/bolt"
	"github.com/Ssawa/destiny/utils"
	"github.com/satori/go.uuid"
)

// Reuse components to speed up serialization and deserialization
var adageBuffer = new(bytes.Buffer)
var adageEncoder = gob.NewEncoder(adageBuffer)
var adageDecoder = gob.NewDecoder(adageBuffer)
var adageMutex sync.Mutex

// Adage is an entry in the database
type Adage struct {
	ID        uuid.UUID
	Body      string
	Tags      []string
	Author    string
	Source    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetAdage gets a random adage from the database
func GetAdage(db *bolt.DB) (*Adage, error) {
	var adage *Adage
	var keys [][]byte

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(adagesKey)
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

		var err error
		adage, err = DeserializeAdage(bucket.Get(choice))

		return err
	})

	return adage, err
}

// Insert the adage to the database
func (adage *Adage) Insert(db *bolt.DB) error {
	utils.Verbose.Println("Inserting adage")

	id := uuid.NewV1()
	utils.Verbose.Println("UUID generated: ", id)

	utils.Verbose.Println("Starting transaction")
	return db.Update(func(tx *bolt.Tx) error {
		dataBucket, err := tx.CreateBucketIfNotExists(adagesKey)
		if err != nil {
			return err
		}

		utils.Verbose.Println("Saving to database")
		data, err := adage.Serialize()
		if err != nil {
			return err
		}

		err = dataBucket.Put(id.Bytes(), data)
		if err != nil {
			return err
		}

		// Update our tags indexes
		tagsBucket, err := tx.CreateBucketIfNotExists(tagsKey)
		if err != nil {
			return err
		}

		for _, tag := range adage.Tags {
			utils.Verbose.Println("Saving index for tag", tag)
			index, err := tagsBucket.CreateBucketIfNotExists([]byte(tag))
			if err != nil {
				return err
			}
			err = index.Put(id.Bytes(), []byte{})
		}

		return nil
	})
}

// Serialize converts the structure to a byte array for saving into the database.
// We're just going to serialize it to JSON for now. I messed around with the
// gob package and it's actually slower to initialize all the object and buffers
// for each serialization and trying to reuse components can lead to corrupt
// binaries. The other alternative is something like Protobufs but I'm not sure
// if it's worth the added complexity. Right now we should be trying to make sure
// we almost never serialize/deserialize underlying data and instead do everything
// through indexes.
func (adage *Adage) Serialize() ([]byte, error) {
	return json.Marshal(*adage)
}

// DeserializeAdage converts a byte array into an Adage struct
func DeserializeAdage(data []byte) (*Adage, error) {
	adage := new(Adage)
	err := json.Unmarshal(data, adage)
	return adage, err
}
