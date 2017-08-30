package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathExists(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	tmpFile, _ := ioutil.TempFile(tmpDir, "")

	exists, err := PathExists(tmpFile.Name())
	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = PathExists(filepath.Join(tmpDir, "DoesNotExist"))
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestEnsureDirectory(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	err := EnsureDirectory(filepath.Join(tmpDir, "thefile.tmp"))
	assert.Nil(t, err)

	noDir := filepath.Join(tmpDir, "nested-dir", "nested-dir-2")
	err = EnsureDirectory(filepath.Join(noDir, "thefile.tmp"))
	assert.Nil(t, err)
	_, err = os.Stat(noDir)
	assert.Nil(t, err)
}
