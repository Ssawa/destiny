package storage

import (
	"log"
	"sync"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type serializeFunc func() ([]byte, error)
type deSerializeFunc func(data []byte) (*Adage, error)

func TestAdageSerializeDeSerialize(t *testing.T) {
	adage := Adage{
		ID:        uuid.NewV4(),
		Body:      "Test Body",
		Tags:      []string{"Tag1", "Tag2"},
		Author:    "Test Author",
		Source:    "Test Source",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testSerializeDeserialize := func(serialize serializeFunc, deserialize deSerializeFunc) {
		data, err := serialize()
		assert.Nil(t, err)

		adage2, err := deserialize(data)
		assert.Nil(t, err)

		assert.Equal(t, adage.ID, adage2.ID)
		assert.Equal(t, adage.Body, adage2.Body)
		assert.Equal(t, adage.Tags, adage2.Tags)
		assert.Equal(t, adage.Author, adage2.Author)
		assert.Equal(t, adage.Source, adage2.Source)
		assert.True(t, adage.CreatedAt.Equal(adage2.CreatedAt))
		assert.True(t, adage.UpdatedAt.Equal(adage2.UpdatedAt))
	}

	testSerial := func(description string, iterations int, serialize serializeFunc, deserialize deSerializeFunc) {
		start := time.Now()
		for i := 0; i < iterations; i++ {
			testSerializeDeserialize(serialize, deserialize)
		}
		elapsed := time.Since(start)
		log.Printf("%s: %s", description, elapsed)
	}

	testParallel := func(description string, workers int, iterations int, serialize serializeFunc, deserialize deSerializeFunc) {
		jobs := make(chan int, 10)
		wg := sync.WaitGroup{}
		start := time.Now()
		for i := 0; i < workers; i++ {
			go func() {
				for _ = range jobs {
					testSerializeDeserialize(serialize, deserialize)
					wg.Done()
				}
			}()
		}

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			jobs <- i
		}
		wg.Wait()
		elapsed := time.Since(start)
		log.Printf("%s: %s", description, elapsed)
	}

	testSerial("Cached Serial One Iteration", 1, adage.Serialize, DeserializeAdage)
	testSerial("Direct Serial One Iteration", 1, adage.SerializeDirect, DeserializeAdageDirect)

	testSerial("Cached Serial Many Iterations", 10000, adage.Serialize, DeserializeAdage)
	testSerial("Direct Serial Many Iterations", 10000, adage.SerializeDirect, DeserializeAdageDirect)
	testParallel("Cached Parallel Many Iterations", 10, 10000, adage.Serialize, DeserializeAdage)
	testParallel("Direct Parallel Many Iterations", 10, 10000, adage.SerializeDirect, DeserializeAdageDirect)
}
