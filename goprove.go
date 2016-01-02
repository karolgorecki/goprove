// Inspect your project for the best practices listed in the Go CheckList.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/karolgorecki/goprove/checklist"
	"github.com/karolgorecki/goprove/util"
)

var (
	flagOutput = flag.String("output", "text", "output formatting")
)

const usageDoc = `goprove: inspect your project for the best practices listed in the Go CheckList

Usage:

  goprove <directory>
`

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
	sourcePath := args[0]

	// Create 2 slices with passed/failed checklist items
	okTasks, nokTasks := checklist.RunTasks(sourcePath)

	printOutput(okTasks, nokTasks)
}

func printOutput(passed []map[string]interface{}, failed []map[string]interface{}) {
	switch *flagOutput {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(struct {
			Passed, Failed []map[string]interface{}
		}{
			passed, failed,
		})
	case "text":
		fmt.Println("Passed tests:", len(passed), "of", len(passed)+len(failed))
		fmt.Println("---------------------------------------------------------------")
		for _, task := range passed {
			fmt.Println(util.FormatSuccess(task["Desc"].(string)))
		}

		fmt.Println("---------------------------------------------------------------")
		for _, task := range failed {
			fmt.Println(util.FormatFail(task["Desc"].(string)))
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(1)
}
