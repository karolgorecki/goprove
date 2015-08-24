package util

import (
	"errors"
	"os"
	"strings"
	"time"

	"fmt"
)

var executionTime time.Time

// BenchmarkStart is a function for starting time counter
func BenchmarkStart() {
	executionTime = time.Now()
}

// ExecutionTime prints the time since the BenchmarkStart
func ExecutionTime() {
	fmt.Printf("\n\x1b[1;30mExecution time %s\x1b[0m\n", time.Since(executionTime))
}

// GetSuccessMessage is a function that returns formatted success message
func GetSuccessMessage(msg string) (successMessage string) {
	return fmt.Sprintf("\x1b[1;32m[✔]\x1b[0m %s", msg)
}

// GetFailMessage is a function that returns formatted fail message
func GetFailMessage(msg string) (failMessage string) {
	return fmt.Sprintf("\x1b[1;31m[✗]\x1b[0m %s", msg)
}

// PrintBanner prints the big Goprove
func PrintBanner() {
	fmt.Println("\x1b[0;32m", `

   ██████╗  ██████╗ ██████╗ ██████╗  ██████╗ ██╗   ██╗███████╗
  ██╔════╝ ██╔═══██╗██╔══██╗██╔══██╗██╔═══██╗██║   ██║██╔════╝
  ██║  ███╗██║   ██║██████╔╝██████╔╝██║   ██║██║   ██║█████╗  
  ██║   ██║██║   ██║██╔═══╝ ██╔══██╗██║   ██║╚██╗ ██╔╝██╔══╝  
  ╚██████╔╝╚██████╔╝██║     ██║  ██║╚██████╔╝ ╚████╔╝ ███████╗
   ╚═════╝  ╚═════╝ ╚═╝     ╚═╝  ╚═╝ ╚═════╝   ╚═══╝  ╚══════╝ 

`, "\x1b[0m")
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
