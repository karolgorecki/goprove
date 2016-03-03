// Package util provides some helper methods.
package util

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"fmt"
)

const (
	nokColor     = "\x1b[1;31m[✗]" + defaultColor + " %s"
	okColor      = "\x1b[1;32m[✔]" + defaultColor + " %s"
	defaultColor = "\x1b[0m"
)

// FormatSuccess decorates a string for the text output.
func FormatSuccess(msg string) (successMessage string) {
	return fmt.Sprintf(okColor, msg)
}

// FormatFail decorates a string for the text output.
func FormatFail(msg string) (failMessage string) {
	return fmt.Sprintf(nokColor, msg)
}

// FilesExistAny checks if the given file(s) exists in the root folder.
func FilesExistAny(path string, files ...string) bool {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Print(err)
		return false
	}

	for _, f := range dirFiles {
		if f.IsDir() {
			continue
		}

		for _, file := range files {
			if strings.Index(strings.ToLower(f.Name()), file) != -1 {
				return true
			}
		}
	}

	return false
}

// FindOccurrencesInTree tries to match the regular expression in files matching the file pattern.
// It returns the number of matchings.
func FindOccurrencesInTree(path, regex, filePattern string) int {
	matches := 0

	err := filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only interested in files
		if f.IsDir() {
			return nil
		}

		if match, err := filepath.Match(filePattern, f.Name()); !match || err != nil {
			return nil
		}

		file, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		r, _ := regexp.Compile(regex)
		match := r.FindStringSubmatch(string(file))

		if len(match) > 0 {
			matches++
		}
		return nil
	})

	if err != nil {
		return 0
	}
	return matches
}
