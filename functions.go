// Package goprove contains lib for checking the Golang best practi
package goprove

import (
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/lint"
	"github.com/ryanuber/go-license"

	"github.com/karolgorecki/goprove/util"
)

func projectBuilds() bool {
	_, err := exec.Command("go", "build", sourcePath).Output()
	return err == nil
}

func isFormatted() bool {
	errors := 0
	filepath.Walk(sourcePath, func(path string, f os.FileInfo, err error) error {
		if !strings.HasSuffix(filepath.Ext(path), ".go") {
			return nil
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		fmtFile, _ := format.Source(file)

		if string(file) != string(fmtFile) {
			errors++
		}
		return nil
	})
	return errors == 0
}

func testPassing() bool {
	output, _ := exec.Command("go", "test", sourcePath+"/...").Output()
	return strings.Index(string(output), `--- FAIL`) == -1
}

func hasLicense() bool {
	if _, err := license.NewFromDir(sourcePath); err != nil {
		return false
	}
	return true
}

func hasReadme() bool {
	return util.FilesExistAny(sourcePath, "readme")
}

func hasContributing() bool {
	return util.FilesExistAny(sourcePath, "contribution", "contribute", "contributing")
}

func isLinted() bool {
	errors := 0
	l := new(lint.Linter)

	filepath.Walk(sourcePath+"/...", func(path string, f os.FileInfo, err error) error {

		if !strings.HasSuffix(filepath.Ext(path), ".go") {
			return nil
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		if lnt, _ := l.Lint(f.Name(), file); len(lnt) > 0 {
			if lnt[0].Confidence > 0.2 {
				errors++
				return nil
			}
		}
		return nil
	})

	return errors == 0
}

func isVetted() bool {
	_, err := exec.Command("go", "vet", sourceGoPath).Output()
	return err == nil
}

func isDirMatch() bool {
	ok := true

	filepath.Walk(sourcePath, func(p string, dir os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !dir.IsDir() || dir.Name() == "." {
			return nil
		}

		// If the dir is "cmd" or it starts with "_" we should skip it
		if dir.IsDir() && (dir.Name() == "cmd" || dir.Name()[0] == '_') {
			return filepath.SkipDir
		}

		files, _ := filepath.Glob(p + string(os.PathSeparator) + "*.go")
		if len(files) == 0 {
			return nil
		}

		file, err := ioutil.ReadFile(files[0])
		if err != nil {
			return err
		}

		r, _ := regexp.Compile(`package ([\w]+)`)
		match := r.FindStringSubmatch(string(file))
		if len(match) > 1 {
			if dir.Name() != match[1] {
				ok = false
			}
		}

		return nil
	})

	return ok
}

func hasBenches() bool {
	return util.FindOccurrencesInTree(sourcePath, `func\sBenchmark\w+\(`, "*_test.go") > 0
}

func hasBlackboxTests() bool {
	return util.FindOccurrencesInTree(sourcePath, `"testing\/quick"`, "*_test.go") > 0
}

func hasBuildPackage() bool {
	return util.FindOccurrencesInTree(sourcePath, `package\smain`, "*.go") > 0
}
