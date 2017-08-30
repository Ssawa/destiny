package storage

import (
	"log"
	"sync"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

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

	testSerializeDeserialize := func(serialize func() ([]byte, error), deserialize func(data []byte) (*Adage, error)) {
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

	iterations := 10000

	start := time.Now()
	for i := 0; i < iterations; i++ {
		testSerializeDeserialize(adage.Serialize, DeserializeAdage)
	}
	elapsed := time.Since(start)
	log.Printf("Optimized Serial Took %s", elapsed)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		testSerializeDeserialize(adage.SerializeSlow, DeserializeAdageSlow)
	}
	elapsed = time.Since(start)
	log.Printf("Slow Serial Took %s", elapsed)

	workers := 10
	jobs := make(chan *Adage, 10)
	wg := sync.WaitGroup{}
	start = time.Now()
	for i := 0; i < workers; i++ {
		go func() {
			for a := range jobs {
				testSerializeDeserialize(a.SerializeSlow, DeserializeAdageSlow)
				wg.Done()
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		jobs <- &adage
	}
	wg.Wait()
	elapsed = time.Since(start)
	log.Printf("Slow Parallel Took %s", elapsed)

	workers = 10
	jobs = make(chan *Adage, 10)
	wg = sync.WaitGroup{}
	start = time.Now()
	for i := 0; i < workers; i++ {
		go func() {
			for a := range jobs {
				testSerializeDeserialize(a.Serialize, DeserializeAdage)
				wg.Done()
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		jobs <- &adage
	}
	wg.Wait()
	elapsed = time.Since(start)
	log.Printf("Optimized Parallel Took %s", elapsed)

}
