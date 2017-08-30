package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessEditor(t *testing.T) {
	os.Setenv("EDITOR", "")
	editor := GuessEditor()
	assert.Equal(t, "vim", editor)

	os.Setenv("EDITOR", "TEST")
	editor = GuessEditor()
	assert.Equal(t, "TEST", editor)
}

func TestGetInputFromEditorUNIX(t *testing.T) {
	// Not really sure the best way of testing this so let's just
	// try it with the "touch" binary. I imagine that means this will
	// only pass on Unix systems. Also this will mostly just be a smoke test
	os.Setenv("EDITOR", "touch")
	_, err := GetInputFromEditor("")
	assert.Nil(t, err)
}

func TestParseFortuneFile(t *testing.T) {
	tmpFile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpFile.Name())

	data := `This is

a
test
%
Another%test
%
How did we do?`

	tmpFile.Write([]byte(data))
	adages := []string{}
	err := ParseFortuneFile(tmpFile.Name(), func(adage string) error {
		adages = append(adages, adage)
		return nil
	})
	assert.Nil(t, err)
	assert.Len(t, adages, 3)
	assert.Equal(t, "This is\n\na\ntest\n", adages[0])
	assert.Equal(t, "Another%test\n", adages[1])
	assert.Equal(t, "How did we do?", adages[2])

	count := 0
	err = ParseFortuneFile(tmpFile.Name(), func(adage string) error {
		count += 1
		return errors.New("From Test")
	})
	assert.NotNil(t, err)
	assert.Equal(t, "From Test", err.Error())
	assert.Equal(t, 1, count)
}
