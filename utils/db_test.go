package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestOpenDatabase(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	db, err := openDatabase(filepath.Join(tmpDir, "nested-folder", "test.db"), nil)
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestOpenReadWrite(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	db, err := OpenReadWrite(filepath.Join(tmpDir, "nested-folder", "test.db"))
	defer db.Close()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("TestBucket"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("Test"), []byte("Hello world"))
		return err
	})
	assert.Nil(t, err)
}

func TestOpenReadOnly(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	db, err := OpenReadOnly(filepath.Join(tmpDir, "nested-folder", "test.db"))
	db.Close()

	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("TestBucket"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("Test"), []byte("Hello world"))
		return err
	})
	assert.Nil(t, err)
}
