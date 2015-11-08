package util

import (
	"errors"
	"os"
	"strings"
	"time"

	"fmt"
)

const (
	nokColor     = "\x1b[1;31m[✗]" + defaultColor + " %s"
	okColor      = "\x1b[1;32m[✔]" + defaultColor + " %s"
	bannerColor  = "\x1b[0;32m"
	defaultColor = "\x1b[0m"
	fadedColor   = "\x1b[1;30m"
)

var executionTime time.Time

// BenchmarkStart is a function for starting time counter
func BenchmarkStart() {
	executionTime = time.Now()
}

// ExecutionTime prints the time since the BenchmarkStart
func ExecutionTime() {
	fmt.Printf("\n"+fadedColor+"Execution time %s"+defaultColor+"\n", time.Since(executionTime))
}

// GetSuccessMessage is a function that returns formatted success message
func GetSuccessMessage(msg string) (successMessage string) {
	return fmt.Sprintf(okColor, msg)
}

// GetFailMessage is a function that returns formatted fail message
func GetFailMessage(msg string) (failMessage string) {
	return fmt.Sprintf(nokColor, msg)
}

// PrintBanner prints the big Goprove
func PrintBanner() {
	fmt.Println(bannerColor, `

=================================
	G O P R O ✔ E
=================================

`, defaultColor)
}

// FileExists check if the given file(s) exists. Returns true if file exists.
func FileExists(files ...string) (fileExists bool, err error) {

	fs, err := os.Open(".")
	if err != nil {
		return false, errors.New("There was a problem with opening the directory")
	}
	defer fs.Close()

	dir, err := fs.Readdir(-1)
	if err != nil {
		return false, errors.New("There was a problem with reading the directory")
	}

FILE_SEARCH:
	for _, file := range dir {
		fileName := strings.ToLower(file.Name())

		for _, file := range files {
			if strings.HasPrefix(fileName, file) {
				fileExists = true
				break FILE_SEARCH
			}
		}
	}

	return fileExists, nil
}
