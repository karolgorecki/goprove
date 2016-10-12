package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/karolgorecki/goprove"
	"github.com/karolgorecki/goprove/util"
)

var (
	flagOutput  = flag.String("output", "text", "output formatting")
	flagExclude = flag.String("exclude", "text", "tasks to exclude")
)

const usageDoc = `goprove: inspect your project for the best practices listed in the Go CheckList

Usage:

  SIMPLE:
  	goprove <directory>

  WITH OUTPUT:
  	goprove -output=<output: json or text> <directory>

  WITH EXCLUDE:
  	goprove -exclude=<tasks: separated by comma> <directory>

Available tasks for exclude:
  projectBuilds, isFormatted, hasLicense, isLinted, isVetted, hasReadme,
  testPassing, isDirMatch, hasContributing, hasBenches, hasBlackboxTests

`

func usage() {
	fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(0)
}

func main() {
	// Use sbrk for allocations rather than GC.
	// Improves performance by 20%.
	_ = os.Setenv("GODEBUG", "sbrk=1")
	flag.Parse()
	log.SetPrefix("goprove: ")

	args := flag.Args()
	if len(args) != 1 {
		usage()
	}

	var excludeTasks []string

	if len(*flagExclude) > 0 {
		excludeTasks = strings.Split(*flagExclude, ",")
	}

	sourcePath := args[0]
	okTasks, nokTasks := goprove.RunTasks(sourcePath, excludeTasks)
	printOutput(okTasks, nokTasks)
}

func printOutput(passed []map[string]interface{}, failed []map[string]interface{}) {
	switch *flagOutput {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(struct {
			Passed []map[string]interface{} `json:"passed"`
			Failed []map[string]interface{} `json:"failed"`
		}{
			passed, failed,
		})
	case "text":
		fmt.Println("Passed tests:", len(passed), "of", len(passed)+len(failed))
		fmt.Println("---------------------------------------------------------------")
		for _, task := range passed {
			fmt.Println(util.FormatSuccess(task["desc"].(string)))
		}

		fmt.Println("---------------------------------------------------------------")
		for _, task := range failed {
			fmt.Println(util.FormatFail(task["desc"].(string)))
		}
	}
}
