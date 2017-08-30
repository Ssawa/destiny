package storage

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

// Reuse components to speed up serialization and deserialization
var buffer = new(bytes.Buffer)
var encoder = gob.NewEncoder(buffer)
var decoder = gob.NewDecoder(buffer)
var mutex sync.Mutex

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

// Serialize converts the structure to a byte array for saving into the database
func (adage *Adage) Serialize() ([]byte, error) {
	mutex.Lock()
	buffer.Reset()
	err := encoder.Encode(*adage)
	if err != nil {
		return nil, err
	}
	data := buffer.Bytes()
	mutex.Unlock()
	return data, nil
}

// DeserializeAdage converts a byte array into an Adage struct
func DeserializeAdage(data []byte) (*Adage, error) {
	mutex.Lock()
	buffer.Reset()
	buffer.Write(data)
	adage := new(Adage)
	err := decoder.Decode(adage)
	mutex.Unlock()
	return adage, err
}

// SerializeSlow is another implementation of Serialize that doesn't reuse cached
// components. See adage_test's TestAdageSerializeDeSerialize for an example of
// the time difference.
func (adage *Adage) SerializeSlow() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(*adage)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// DeserializeAdageSlow is another implementation of DeserializeAdage that doesn't
// reuse cached components. See adage_test's TestAdageSerializeDeSerialize for an
// example of the time difference.
func DeserializeAdageSlow(data []byte) (*Adage, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	adage := new(Adage)
	err := decoder.Decode(adage)
	return adage, err
}
