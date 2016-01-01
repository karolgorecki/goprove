package checklist

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/karolgorecki/goprove/util"
)

func projectBuilds() bool {
	_, err := exec.Command("go", "build", sourcePath).Output()
	return err == nil
}

func isFormatted() bool {
	output, _ := exec.Command("gofmt", "-l -s", sourcePath).Output()
	return len(output) == 0
}

func testPassing() bool {
	output, _ := exec.Command("go", "test", sourcePath+"/...").Output()
	return strings.Index(string(output), `--- FAIL`) == -1
}

func hasLicense() bool {
	return util.FilesExistAny(sourcePath, "license", "licensing", "copying")
}

func hasReadme() bool {
	return util.FilesExistAny(sourcePath, "readme")
}

func hasContributing() bool {
	return util.FilesExistAny(sourcePath, "contribution", "contribute", "contributing")
}

func isLinted() bool {
	output, _ := exec.Command("golint", sourcePath+"/...").Output()
	return len(output) == 0
}

func isVetted() bool {
	_, err := exec.Command("go", "vet", sourcePath).Output()
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
	return util.FindPatternInTree(sourcePath, `func\sBenchmark\w+\(`, "*_test.go")
}

func hasBlackboxTests() bool {
	return util.FindPatternInTree(sourcePath, `"testing\/quick"`, "*_test.go")
}