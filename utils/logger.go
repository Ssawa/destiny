package utils

import (
	"io/ioutil"
	"log"
)

// Verbose is used to log verbose commands
var Verbose = log.New(ioutil.Discard, "", 0)
