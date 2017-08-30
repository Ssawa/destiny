package utils

import (
	"io/ioutil"
	"os"
	"os/exec"
)

// GetInputFromEditor spawns an editor and gets user input from there.
// Pass in a non empty string to override automatically deducing the editor
func GetInputFromEditor(editor string) (string, error) {
	if editor == "" {
		editor = GuessEditor()
	}

	Verbose.Println("Creating temporary file")
	tmpFile, err := ioutil.TempFile("", "")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		return "", err
	}

	Verbose.Println("Spawning editor:", editor)
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	tmpFile.Seek(0, 0)
	output, err := ioutil.ReadAll(tmpFile)

	return string(output), err
}

// GuessEditor tries to determine which editor to spawn
func GuessEditor() string {
	Verbose.Println("Guessing which editor to spawn")
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	return editor
}
