package utils

import (
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
